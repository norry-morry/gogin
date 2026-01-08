package controller

import (
	"net/http"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"resume/internal/adapter/http/middleware"
	repo "resume/internal/domain/repository"
	vo "resume/internal/domain/valueobject/i18n"
	uci18n "resume/internal/usecase/i18n"
)

// I18nHandler は i18n 関連エンドポイントのハンドラです。
type I18nHandler struct {
	uc   uci18n.Usecase
	repo repo.DictionaryRepository
}

// NewI18nHandler は I18nHandler を生成します。
func NewI18nHandler(
	uc uci18n.Usecase,
	repo repo.DictionaryRepository,
) *I18nHandler {
	return &I18nHandler{
		uc:   uc,
		repo: repo,
	}
}

// GetBundle は prefix で辞書をまとめて返します。
// GET /i18n/bundle?prefix=ui.page.profile.
// lang はミドルウェアで決まったものを使用。
// レスポンスは { "key": "value", ... } のフラットなマップとします。
func (h *I18nHandler) GetBundle(c *gin.Context) {
	loc := vo.NewLocale(middleware.From(c))
	prefix := c.Query("prefix")

	out, err := h.uc.ListBundle(c.Request.Context(), uci18n.ListBundleInput{
		Locale: loc,
		Prefix: prefix,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "i18n_bundle_failed", "message": err.Error(),
		}})
		return
	}
	// locale はレスポンスには含めず、辞書だけを返す
	c.JSON(http.StatusOK, out.Data) // or out.Bundle など、実際のフィールド名に合わせて
}

// Reload は辞書を再読込し、キャッシュを差し替えます。
// POST /i18n/reload（通常は local か管理者のみ）
func (h *I18nHandler) Reload(c *gin.Context) {
	if os.Getenv("APP_ENV") != "local" {
		c.Status(http.StatusForbidden)
		return
	}
	if err := h.uc.Reload(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "i18n_reload_failed", "message": err.Error(),
		}})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListLocales は利用可能なロケールの一覧を返します。
// GET /i18n/locales
func (h *I18nHandler) ListLocales(c *gin.Context) {
	all, err := h.repo.LoadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "i18n_list_failed", "message": err.Error(),
		}})
		return
	}
	type item struct {
		Code  string `json:"code"`
		Label string `json:"label"`
	}
	list := make([]item, 0, len(all))
	for loc := range all {
		code := loc.Code()
		lbl := labelFor(code)
		list = append(list, item{Code: code, Label: lbl})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Code < list[j].Code })
	c.JSON(http.StatusOK, gin.H{"locales": list})
}

func labelFor(code string) string {
	if tag, err := language.Parse(code); err == nil {
		return display.Self.Name(tag) // 例: ja → 日本語, en → English
	}
	return code
}

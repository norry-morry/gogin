// Package controller は メタ情報（国・住所用途・性別など）に関する
// 参照系エンドポイントを提供します。
package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"resume/internal/adapter/http/handlerutil"
	"resume/internal/adapter/http/middleware"
	ucmeta "resume/internal/usecase/meta"
)

// MetaHandler は、メタ情報（国、住所用途、性別など）を提供する
// HTTP エンドポイントのハンドラです。
type MetaHandler struct {
	uc ucmeta.Usecase
}

// NewMetaHandler は、メタ情報ユースケースを注入した新しい MetaHandler を生成します。
func NewMetaHandler(
	uc ucmeta.Usecase,
) *MetaHandler {
	return &MetaHandler{
		uc: uc,
	}
}

// ListCountry は、国一覧を返すエンドポイントです。
// クエリパラメータ ?only_supported=false を指定すると、
// 非サポート国も含めた全件を返します。
// 例: GET /meta/country?only_supported=false
func (h *MetaHandler) ListCountry(c *gin.Context) {
	lang := middleware.From(c)
	// ?only_supported=false の場合のみ false、それ以外は true
	only := c.DefaultQuery("only_supported", "true") != "false"
	in := ucmeta.ListCountryInput{
		Locale:        lang,
		OnlySupported: only,
	}
	out, err := h.uc.ListCountry(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// ListAddressPurpose は、住所用途（billing, home, office など）の一覧を返すエンドポイントです。
// ロケールに応じて翻訳済みラベルを返します。
// 例: GET /meta/address-purpose?lang=ja
func (h *MetaHandler) ListAddressPurpose(c *gin.Context) {
	lang := middleware.From(c)
	in := ucmeta.ListAddressPurposeInput{
		Locale: lang,
	}
	out, err := h.uc.ListAddressPurpose(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// ListGender は、性別（male, female, other, unspecified）の一覧を返すエンドポイントです。
// ロケールに応じて翻訳済みラベルを返します。
// 例: GET /meta/gender?lang=ja
func (h *MetaHandler) ListGender(c *gin.Context) {
	lang := middleware.From(c)
	in := ucmeta.ListGenderInput{
		Locale: lang,
	}
	out, err := h.uc.ListGender(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// ListEducationStatus は 学歴状態の一覧を返すエンドポイントです。
// ロケールに応じて翻訳済みラベルを返します。
// 例: GET /meta/education/status?lang=ja
func (h *MetaHandler) ListEducationStatus(c *gin.Context) {
	lang := middleware.From(c)
	slog.Debug("ListEducationStatus called. lang:", lang, lang)
	in := ucmeta.ListEducationStatusInput{
		Locale: lang,
	}
	out, err := h.uc.ListEducationStatus(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// ListDegreeType は 学位種別の一覧を返すエンドポイントです。
// ロケールに応じて翻訳済みラベルを返します。
// 例: GET /meta/education/degree?lang=ja
func (h *MetaHandler) ListDegreeType(c *gin.Context) {
	lang := middleware.From(c)
	slog.Debug("ListEducationStatus called. lang:", lang, lang)
	in := ucmeta.ListDegreeTypeInput{
		Locale: lang,
	}
	out, err := h.uc.ListDegreeType(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

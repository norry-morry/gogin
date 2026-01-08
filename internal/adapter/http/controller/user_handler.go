// Package controller はエンドポイント処理を担当します。
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	resp "resume/internal/adapter/http/dto/response"
	"resume/internal/adapter/http/handlerutil"
	"resume/internal/shared/apperr"
	"resume/internal/usecase/user"
)

// UserHandler はユーザー情報用のハンドラです。
type UserHandler struct {
	uc user.Usecase
}

// NewUserHandler は UserHandler の新しいインスタンスを生成して返します。
// ユーザー情報の HTTP リクエストを処理するためのハンドラを初期化します。
func NewUserHandler(uc user.Usecase) *UserHandler {
	return &UserHandler{uc: uc}
}

// GetByID は GET /api/users/:id を処理します。
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		handlerutil.WriteError(c,
			apperr.New(apperr.CodeBadRequest, "invalid user id", map[string]any{
				"param": idStr,
			}))
		return
	}

	out, err := h.uc.GetByID(c.Request.Context(), user.GetUserInput{ID: id64})
	if err != nil {
		handlerutil.WriteError(c, err) // usecase側がapperr返却する想定
		return
	}

	c.JSON(
		http.StatusOK,
		resp.FromUsecase(
			out.ID,
			out.DisplayName,
			out.Email,
		),
	)
}

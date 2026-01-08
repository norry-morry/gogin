// Package controller は 学歴に関するエンドポイントを提供します。
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"resume/internal/adapter/http/dto/request"
	"resume/internal/adapter/http/handlerutil"
	"resume/internal/shared/apperr"
	"resume/internal/shared/ctx/auth"
	"resume/internal/shared/util"
	uceducation "resume/internal/usecase/education"
)

// EducationHandler は 学歴 を提供するHTTPエンドポイントのハンドラです
type EducationHandler struct {
	uc uceducation.Usecase
}

// NewEducationHandler は 学歴ユースケースを注入した新しい EducationHandler を生成します
func NewEducationHandler(
	uc uceducation.Usecase,
) *EducationHandler {
	return &EducationHandler{
		uc: uc,
	}
}

// ListEducation は 現在ログイン中のユーザーに紐づく学歴一覧を取得します。
func (h *EducationHandler) ListEducation(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(
			c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil),
		)
		return
	}
	in := uceducation.ListInput{
		UserID: userID,
	}
	out, err := h.uc.ListUserEducation(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// AddEducation は 認証済みユーザーが 学歴情報 を登録する
func (h *EducationHandler) AddEducation(c *gin.Context) {
	_, ok := auth.From(c.Request.Context())
	if !ok {
		handlerutil.WriteError(
			c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil),
		)
		return
	}
	var req request.AddEducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			env := handlerutil.BuildValidationEnvelope(c, verrs, "domain.userEducation")
			handlerutil.WriteValidationError(c, env)
			return
		}
		handlerutil.WriteError(c, err)
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(
			c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil),
		)
		return
	}
	eventDate, err := util.ParseYearMonthVal(req.EventDate)
	if err != nil {
		handlerutil.WriteError(c, apperr.New(apperr.CodeBadRequest, "invalid enrollment date", nil))
	}
	in := uceducation.CreateInput{
		UserID:            userID,
		InstitutionName:   req.InstitutionName,
		FacultyName:       req.FacultyName,
		DepartmentName:    req.DepartmentName,
		DegreeTypeID:      req.DegreeTypeID,
		EducationStatusID: req.EducationStatusID,
		EventDate:         eventDate,
		Description:       req.Description,
		IsPublic:          req.IsPublic,
	}
	out, err := h.uc.CreateUserEducation(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// UpdateEducation は 認証済みユーザーが 学歴情報 を更新する
func (h *EducationHandler) UpdateEducation(c *gin.Context) {
	_, ok := auth.From(c.Request.Context())
	if !ok {
		handlerutil.WriteError(
			c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil),
		)
		return
	}
	var req request.UpdateEducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			env := handlerutil.BuildValidationEnvelope(c, verrs, "domain.userEducation")
			handlerutil.WriteValidationError(c, env)
			return
		}
		handlerutil.WriteError(c, err)
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(
			c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil),
		)
		return
	}
	educationIDStr := c.Param("education_id")
	educationID, err := strconv.ParseUint(educationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid education_id",
		})
		return
	}
	eventDate, err := util.ParseYearMonthVal(req.EventDate)
	if err != nil {
		handlerutil.WriteError(c, apperr.New(apperr.CodeBadRequest, "invalid enrollment date", nil))
	}
	in := uceducation.UpdateInput{
		EducationID:       educationID,
		UserID:            userID,
		InstitutionName:   req.InstitutionName,
		FacultyName:       req.FacultyName,
		DepartmentName:    req.DepartmentName,
		DegreeTypeID:      req.DegreeTypeID,
		EducationStatusID: req.EducationStatusID,
		EventDate:         eventDate,
		Description:       req.Description,
		IsPublic:          req.IsPublic,
	}
	out, err := h.uc.UpdateUserEducation(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// DeleteEducation は 認証済みユーザーが 学歴情報 を削除する
// 削除後 他の学歴で並び順を変更する
func (h *EducationHandler) DeleteEducation(c *gin.Context) {
	educationIDStr := c.Param("education_id")
	educationID, err := strconv.ParseUint(educationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid education_id",
		})
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := uceducation.DeleteInput{EducationID: educationID, UserID: userID}
	_, err = h.uc.DeleteUserEducation(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// OrderEducation は 認証済みユーザーが 学歴情報 の 並び順を変更 する
func (h *EducationHandler) OrderEducation(c *gin.Context) {
	_, ok := auth.From(c.Request.Context())
	if !ok {
		handlerutil.WriteError(
			c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil),
		)
		return
	}
	var req request.OrderEducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			env := handlerutil.BuildValidationEnvelope(c, verrs, "domain.userEducation")
			handlerutil.WriteValidationError(c, env)
			return
		}
		handlerutil.WriteError(c, err)
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(
			c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil),
		)
		return
	}
	in := uceducation.ReorderInput{
		UserID:       userID,
		EducationIDs: req.EducationIDs,
	}
	_, err := h.uc.ReorderEducation(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ExistsEducation は 認証済みユーザーが 学歴情報 の 有無を返す
func (h *EducationHandler) ExistsEducation(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := uceducation.HasEducationInput{
		UserID: userID,
	}
	out, err := h.uc.HasUserEducation(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

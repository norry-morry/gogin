// Package controller は フリーランサー情報のエンドポイント処理を担当します。
package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	fba "resume/internal/adapter/gateway/firebase"
	"resume/internal/adapter/http/dto/request"
	"resume/internal/adapter/http/handlerutil"
	"resume/internal/adapter/http/presenter"
	"resume/internal/shared/apperr"
	"resume/internal/shared/ctx/auth"
	"resume/internal/usecase/profile"
)

// ProfileHandler は フリーランサー情報のハンドラです
type ProfileHandler struct {
	fb fba.Auth
	uc profile.Usecase
}

// NewProfileHandler は ProfileHandler の新しいインスタンスを生成して返す
// フリーランサー情報の HTTP リクエストを処理するためのハンドラを初期化します
func NewProfileHandler(
	fb fba.Auth,
	uc profile.Usecase,
) *ProfileHandler {
	return &ProfileHandler{
		fb: fb,
		uc: uc,
	}
}

// Fba は認証済みユーザーのFirebaseから取得できる情報を返します。
// コンテキストから認証クレームを取得し、メールアドレスとUID、UserIdをjson形式で返します
// 多分表立っては使いません
func (h *ProfileHandler) Fba(c *gin.Context) {
	cl, ok := auth.From(c.Request.Context())
	if !ok {
		handlerutil.WriteError(c,
			apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
	}
	c.JSON(http.StatusOK, gin.H{"uid": cl.UID, "email": cl.Email, "id": userID})
}

// Me は 認証済みユーザーの情報を返します
// コンテキストから認証クレームを取得し、UIDとメールアドレスをJSON形式で応答します。
func (h *ProfileHandler) Me(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
	}
	in := profile.UserInput{
		UserID: userID,
	}
	out, err := h.uc.GetUser(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
	}
	c.JSON(http.StatusOK, out)
}

// Identities は認証ユーザに紐づく外部ID一覧をページングして返します。
// クエリ: page(>=1), per_page(-1|1..100), sort(created_at|provider|uid|email_at_signup), order(asc|desc), q(任意)
func (h *ProfileHandler) Identities(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	var q request.ListIdentitiesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		// ▼ ここで必ず validation 用のユーティリティに流す
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			slog.Info("validator.ValidationErrors hit", "errs", verrs)
			// scope は他と揃えて domain.user_identity にしておく
			env := handlerutil.BuildValidationEnvelope(c, verrs, "meta")
			handlerutil.WriteValidationError(c, env) // 422 + 新フォーマット
			return
		}

		// バリデーション以外の Bind エラーだけ通常の WriteError
		handlerutil.WriteError(c, err)
		return
	}
	q.Normalize()

	in := profile.UserIdentityInput{
		UserID:    userID,
		Page:      q.Page,
		PerPage:   q.PerPage,
		Sort:      q.Sort,
		SortOrder: q.Order,
		Query:     q.Q,
	}
	//slog.Debug("uid=%d page=%d per=%d sort=%s order=%s q=%v", userID, page, perPage, col, order, q)
	out, err := h.uc.ListUserIdentity(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch identities"})
	}
	c.JSON(http.StatusOK, out)
}

// GetPersonalInfo は 現在ログイン中のユーザーに紐づくプロフィール情報を取得します
// 成功時には HTTP 200 でプロフィール情報を JSON 形式で返します
// 認証エラー時は HTTP 401、内部エラー時は HTTP 500 を返します
func (h *ProfileHandler) GetPersonalInfo(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := profile.UserInput{
		UserID: userID,
	}
	out, err := h.uc.GetUserProfile(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user_profile"})
	}
	res := presenter.PresentUserProfile(out)
	c.JSON(http.StatusOK, res)
}

// PatchPersonalInfo は 現在ログイン中のユーザーに紐づくプロフィール情報を登録・更新します
// 成功時には HTTP 200 で返します
func (h *ProfileHandler) PatchPersonalInfo(c *gin.Context) {
	var req request.PatchProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			env := handlerutil.BuildValidationEnvelope(c, verrs, "domain.user_profile")
			handlerutil.WriteValidationError(c, env) // 422 + {"error": {...}}
			return
		}
		handlerutil.WriteError(c, err)
		return
	}
	req.Normalize()
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := profile.PatchUserProfileInput{
		UserID:         userID,
		FamilyName:     req.FamilyName,
		GivenName:      req.GivenName,
		FamilyNameKana: req.FamilyNameKana,
		GivenNameKana:  req.GivenNameKana,
		BirthDate:      req.BirthDate,
		GenderID:       req.GenderID,
		Initial:        req.Initial,
	}
	if err := h.uc.PatchUserProfile(c.Request.Context(), in); err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HasPersonalInfo は 現在ログイン中のユーザーに紐づくプロフィール情報の有無を取得します
func (h *ProfileHandler) HasPersonalInfo(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := profile.UserInput{
		UserID: userID,
	}
	out, err := h.uc.HasUserProfile(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// HasUserAddress は 現在ログイン中のユーザーに紐づく住所情報の有無を取得します
func (h *ProfileHandler) HasUserAddress(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := profile.UserInput{
		UserID: userID,
	}
	out, err := h.uc.HasUserAddress(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// ListUserAddress は 現在ログイン中のユーザーに紐づく住所一覧を取得します。
// クエリでページング・ソート・目的IDなどを指定できます。
func (h *ProfileHandler) ListUserAddress(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	var q request.ListUserAddressQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			env := handlerutil.BuildValidationEnvelope(c, verrs, "meta")
			handlerutil.WriteValidationError(c, env) // 422 + 新フォーマット
			return
		}
		handlerutil.WriteError(c, err)
		return
	}
	q.Normalize()

	in := profile.ListUserAddressInput{
		UserID:    userID,
		Page:      q.Page,
		PerPage:   q.PerPage,
		Sort:      q.Sort,
		SortOrder: q.Order,

		PurposeID: q.PurposeID,
	}
	slog.Debug("handler", "in", in)
	out, err := h.uc.ListUserAddress(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch user address",
		})
	}

	c.JSON(http.StatusOK, out)
}

// CreateUserAddress は 認証済みユーザーが住所を登録する
func (h *ProfileHandler) CreateUserAddress(c *gin.Context) {
	var req request.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			env := handlerutil.BuildValidationEnvelope(c, verrs, "domain.address")
			handlerutil.WriteValidationError(c, env) // 422 + {"error": {...}}
			return
		}
		handlerutil.WriteError(c, err)
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := profile.CreateUserAddressInput{
		UserID:             userID,
		PurposeID:          *req.PurposeID,
		IsPrimary:          false,
		CountryCode:        req.CountryCode,
		AdministrativeArea: req.AdministrativeArea,
		Locality:           req.Locality,
		DependentLocality:  req.DependentLocality,
		PostalCode:         req.PostalCode,
		AddressLine1:       req.AddressLine1,
		AddressLine2:       req.AddressLine2,
		AddressLine3:       req.AddressLine3,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
	}
	out, err := h.uc.CreateUserAddress(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	res := presenter.ToAddressResponse(out.UserAddress)

	c.Header("location", fmt.Sprintf("/fl/address/%d", res.ID))
	c.JSON(http.StatusCreated, res)
}

// UpdateUserAddress は 現在ログイン中のユーザーに紐づく住所情報を更新します。
// パスパラメータの address_id と JSON ボディの内容をもとに更新を行います。
func (h *ProfileHandler) UpdateUserAddress(c *gin.Context) {
	// 1. パスパラメータから address_id を取得
	addressIDStr := c.Param("address_id")

	// 2. uint64 に変換
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		// 不正な ID → 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid address_id",
		})
		return
	}
	var req request.UpdateUserAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			env := handlerutil.BuildValidationEnvelope(c, verrs, "domain.address")
			handlerutil.WriteValidationError(c, env) // 422 + {"error": {...}}
			return
		}
		handlerutil.WriteError(c, err)
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnauthorized, "unauthorized", nil))
		return
	}
	in := profile.UpdateUserAddressInput{
		AddressID: addressID,
		CreateUserAddressInput: profile.CreateUserAddressInput{
			UserID:             userID,
			PurposeID:          *req.PurposeID,
			IsPrimary:          false,
			CountryCode:        req.CountryCode,
			AdministrativeArea: req.AdministrativeArea,
			Locality:           req.Locality,
			DependentLocality:  req.DependentLocality,
			PostalCode:         req.PostalCode,
			AddressLine1:       req.AddressLine1,
			AddressLine2:       req.AddressLine2,
			AddressLine3:       req.AddressLine3,
			Latitude:           req.Latitude,
			Longitude:          req.Longitude,
		},
	}
	out, err := h.uc.UpdateUserAddress(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	res := presenter.ToAddressResponse(out.UserAddress)
	c.Header("location", fmt.Sprintf("/fl/address/%d", res.ID))
	c.JSON(http.StatusOK, res)
}

// DeleteUserAddress は 現在ログイン中のユーザーに紐付く住所を削除します
// パスパラメータとfirebaseから求められたUserIDを元に削除を行います
func (h *ProfileHandler) DeleteUserAddress(c *gin.Context) {
	addressIGStr := c.Param("address_id")

	addressID, err := strconv.ParseUint(addressIGStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid address_id",
		})
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnprocessable, "unauthorized", nil))
		return
	}

	in := profile.DeleteUserAddressInput{AddressID: addressID, UserID: userID}
	_, err = h.uc.DeleteUserAddress(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}

	// ★ REST 的にもっとも自然なレスポンス
	c.Status(http.StatusNoContent)
}

// DetailUserAddress は 認証済みユーザーの住所を取得します
func (h *ProfileHandler) DetailUserAddress(c *gin.Context) {
	addressIDStg := c.Param("address_id")

	addressID, err := strconv.ParseUint(addressIDStg, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid address id",
		})
		return
	}
	userID, ok := auth.UserID(c)
	if !ok || userID == 0 {
		handlerutil.WriteError(c, apperr.New(apperr.CodeUnprocessable, "unauthorized", nil))
		return
	}
	in := profile.DetailUserAddressInput{
		AddressID: addressID,
		UserID:    userID,
	}
	out, err := h.uc.DetailUserAddress(c.Request.Context(), in)
	if err != nil {
		handlerutil.WriteError(c, err)
		return
	}
	res := presenter.ToAddressResponse(out.UserAddress)
	c.Header("location", fmt.Sprintf("/fl/address/%d", res.ID))
	c.JSON(http.StatusOK, res)
}

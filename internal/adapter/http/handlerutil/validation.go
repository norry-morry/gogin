package handlerutil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	vtrans "resume/internal/adapter/validation/rules"
	"resume/internal/shared/requestid"
	"resume/internal/shared/util"
)

// BuildValidationEnvelope は vtrans.MapValidationErrors の結果を “既存エンベロープに詰める”
func BuildValidationEnvelope(
	c *gin.Context, // gin.Context 互換（GetHeader, MustGetなど必要なものだけ）
	verrs validator.ValidationErrors,
	scope string, // 例: "domain.address"
) ErrorEnvelope {

	//locale := middleware.From(c) // あなたの i18n middleware から取得
	// ① i18n 版 details を作る
	i18nDetails := vtrans.MapValidationErrors(verrs, scope)
	// ② 旧形（param/tag/value）も必要なら作る
	legacy := map[string]map[string]string{}
	for _, fe := range verrs {
		// フィールドキーは既存に合わせる（例: StructField ではなく scope + snake/camel 統一）
		fieldKey := scope + "." + util.ToCamel(fe.Field())
		legacy[fieldKey] = map[string]string{
			"param": fe.Param(),
			"tag":   fe.Tag(),
			"value": fmt.Sprint(fe.Value()), // 型に関係なく安全に文字列化
		}
	}

	return ErrorEnvelope{
		Code:          "UNPROCESSABLE",
		Message:       "The request contains semantically invalid data.",
		RequestID:     requestid.Get(c), // 既存のリクエストID取り出し
		Details:       i18nDetails,      // ← 新フォーマット（フロントはここを見る）
		LegacyDetails: legacy,           // ← 当面併載。要らなくなったら削除
	}
}

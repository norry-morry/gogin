// Package request はtest_handler.createのリクエストバリデータ
package request

import "time"

// CreateTestRequest はテスト用の入力DTO
type CreateTestRequest struct {
	// 必須・文字数制約
	Name string `json:"name" binding:"required,min=2,max=50"`

	// Email型チェック
	Email string `json:"email" binding:"required,email"`

	// パスワード強度（カスタムルール）
	Password string `json:"password" binding:"required,passwd_strong"`

	// 日付相関（from < to のときのみ有効）
	From *time.Time `json:"from" time_format:"2006-01-02" binding:"omitempty,ltfield=To,required_with=To"`
	To   *time.Time `json:"to"   time_format:"2006-01-02" binding:"omitempty,gtfield=From,required_with=From"`

	// 数値の最小・最大
	Quantity int `json:"quantity" binding:"gte=1,lte=100"`

	// 数値同士の相関（minAmount < maxAmount）
	MinAmount *int `json:"minAmount" binding:"omitempty,gte=0,ltfield=MaxAmount,required_with=MaxAmount"`
	MaxAmount *int `json:"maxAmount" binding:"omitempty,gtfield=MinAmount,required_with=MinAmount"`
}

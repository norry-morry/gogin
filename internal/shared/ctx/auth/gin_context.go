// Package auth は Gin の Context で認証情報を扱うヘルパを提供します。
package auth

import "github.com/gin-gonic/gin"

type key string

const (
	kUserID      key = "userID"
	kUID         key = "UID"
	kFirebaseUID key = "firebaseUID"
	kEmail       key = "email"
)

// SetUserID はユーザーIDを Gin の Context に設定します。
func SetUserID(c *gin.Context, id uint64) {
	c.Set(string(kUserID), id)
}

// SetFirebaseUID は Firebase UID を設定します。
func SetFirebaseUID(c *gin.Context, uid string) {
	c.Set(string(kUID), uid)
	c.Set(string(kFirebaseUID), uid)
}

// SetEmail はメールアドレスを設定します。
func SetEmail(c *gin.Context, email *string) {
	c.Set(string(kEmail), email)
}

// UserID は設定されたユーザーIDを取得します。
func UserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get(string(kUserID))
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}

// FirebaseUID は設定された Firebase UID を取得します。
func FirebaseUID(c *gin.Context) (string, bool) {
	v, ok := c.Get(string(kFirebaseUID))
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// UID は設定されたUID（互換キー）を取得します。
func UID(c *gin.Context) (string, bool) {
	v, ok := c.Get(string(kUID))
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// Email は設定されたメールアドレスを取得します。
func Email(c *gin.Context) (*string, bool) {
	v, ok := c.Get(string(kEmail))
	if !ok {
		return nil, false
	}
	if v == nil {
		return nil, true // セット済みだが nil
	}
	email, ok := v.(*string)
	return email, ok
}

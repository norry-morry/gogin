package log

import (
	"log/slog"
	"os"
)

// NewBaseLogger は JSON 構造化・ローテーション設定済みのベースロガーを返します。
func NewBaseLogger(path string) *slog.Logger {
	var w *os.File
	if path == "" || path == "-" {
		w = os.Stdout
	} else {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			w = os.Stdout
		} else {
			w = f
		}
	}
	// JSON ハンドラ(タイムスタンプ・レベルなどはhandler側で付与)
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	})
	return slog.New(h)
}

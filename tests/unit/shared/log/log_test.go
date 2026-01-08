// Package logutil は internal/shared/logutil の multiHandler に対する
// ユニットテストを提供します。
package log

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	logutil "resume/internal/shared/log"
)

// fakeHandler は slog.Handler のテスト用スタブです。
type fakeHandler struct {
	enabled        bool
	enabledCalls   int
	handleCalls    int
	handleErr      error
	withAttrsCalls int
	withGroupCalls int
	lastRecord     slog.Record
}

func (f *fakeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	f.enabledCalls++
	return f.enabled
}

func (f *fakeHandler) Handle(ctx context.Context, r slog.Record) error {
	f.handleCalls++
	f.lastRecord = r
	return f.handleErr
}

func (f *fakeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	f.withAttrsCalls++
	return f
}

func (f *fakeHandler) WithGroup(name string) slog.Handler {
	f.withGroupCalls++
	return f
}

//
// ---------- NewMultiHandler / Enabled ----------
//

// TestMultiHandler_EnabledAggregates は、いずれかの Handler が有効なら true を返すことを確認する。
func TestMultiHandler_EnabledAggregates(t *testing.T) {
	t.Parallel()

	h1 := &fakeHandler{enabled: false}
	h2 := &fakeHandler{enabled: true}

	mh := logutil.NewMultiHandler(h1, h2)

	ctx := context.Background()
	if !mh.Enabled(ctx, slog.LevelInfo) {
		t.Fatalf("expected Enabled to return true when one handler is enabled")
	}

	if h1.enabledCalls == 0 || h2.enabledCalls == 0 {
		t.Errorf("expected Enabled to be called on all handlers, got h1=%d, h2=%d",
			h1.enabledCalls, h2.enabledCalls)
	}
}

//
// ---------- Handle ----------
//

// TestMultiHandler_HandleCallsOnlyEnabled は Enabled が true の Handler のみ Handle が呼ばれることを確認する。
func TestMultiHandler_HandleCallsOnlyEnabled(t *testing.T) {
	t.Parallel()

	h1 := &fakeHandler{enabled: false}
	h2 := &fakeHandler{enabled: true}
	h3 := &fakeHandler{enabled: true}

	mh := logutil.NewMultiHandler(h1, h2, h3)

	ctx := context.Background()
	rec := slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", 0)

	if err := mh.Handle(ctx, rec); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if h1.handleCalls != 0 {
		t.Errorf("expected h1.Handle not to be called, got %d", h1.handleCalls)
	}
	if h2.handleCalls != 1 {
		t.Errorf("expected h2.Handle to be called once, got %d", h2.handleCalls)
	}
	if h3.handleCalls != 1 {
		t.Errorf("expected h3.Handle to be called once, got %d", h3.handleCalls)
	}
}

// TestMultiHandler_HandleErrorAggregation は最初のエラーを返しつつ、
// 他の Handler にも処理を続行することを確認する。
func TestMultiHandler_HandleErrorAggregation(t *testing.T) {
	t.Parallel()

	errFirst := errors.New("first")
	h1 := &fakeHandler{enabled: true, handleErr: errFirst}
	h2 := &fakeHandler{enabled: true, handleErr: nil}

	mh := logutil.NewMultiHandler(h1, h2)

	ctx := context.Background()
	rec := slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", 0)

	err := mh.Handle(ctx, rec)
	if !errors.Is(err, errFirst) {
		t.Fatalf("expected first error to be returned, got %v", err)
	}

	if h1.handleCalls != 1 || h2.handleCalls != 1 {
		t.Errorf("expected both handlers to be called, got h1=%d, h2=%d",
			h1.handleCalls, h2.handleCalls)
	}
}

//
// ---------- WithAttrs / WithGroup ----------
//

// TestMultiHandler_WithAttrs は各 Handler の WithAttrs が呼ばれることを確認する。
func TestMultiHandler_WithAttrs(t *testing.T) {
	t.Parallel()

	h1 := &fakeHandler{enabled: true}
	h2 := &fakeHandler{enabled: true}

	mh := logutil.NewMultiHandler(h1, h2)

	attr := slog.String("key", "value")
	mh2 := mh.WithAttrs([]slog.Attr{attr})

	if mh2 == nil {
		t.Fatalf("expected non-nil handler from WithAttrs")
	}

	if h1.withAttrsCalls != 1 || h2.withAttrsCalls != 1 {
		t.Errorf("expected WithAttrs to be called once on each handler, got h1=%d, h2=%d",
			h1.withAttrsCalls, h2.withAttrsCalls)
	}
}

// TestMultiHandler_WithGroup は各 Handler の WithGroup が呼ばれることを確認する。
func TestMultiHandler_WithGroup(t *testing.T) {
	t.Parallel()

	h1 := &fakeHandler{enabled: true}
	h2 := &fakeHandler{enabled: true}

	mh := logutil.NewMultiHandler(h1, h2)

	mh2 := mh.WithGroup("grp")
	if mh2 == nil {
		t.Fatalf("expected non-nil handler from WithGroup")
	}

	if h1.withGroupCalls != 1 || h2.withGroupCalls != 1 {
		t.Errorf("expected WithGroup to be called once on each handler, got h1=%d, h2=%d",
			h1.withGroupCalls, h2.withGroupCalls)
	}
}

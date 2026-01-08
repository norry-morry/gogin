package di_test

import (
	"context"
	"testing"

	"resume/internal/di"
)

func TestBackground_ReturnsContextBackground(t *testing.T) {
	t.Parallel()

	ctx := di.ExportBackground()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
	if ctx != context.Background() {
		t.Errorf("expected context.Background(), got %#v", ctx)
	}
}

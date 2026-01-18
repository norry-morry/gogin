// Package meta_test は internal/usecase/meta のユニットテストを行います。
// 本ファイルは ListGender インタラクタのユニットテストです。
package meta_test

import (
	"context"
	"errors"
	"testing"

	"resume/internal/domain/entity"
	"resume/internal/usecase/meta"
)

// ============================================================
// 正常系
// ============================================================
func TestListGender_Normal(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	gn := &fakeGnRepo{
		ReturnList: []entity.Gender{
			{ID: 1, Code: "male"}, // 他フィールドは全部ゼロ値でOK
			{ID: 2, Code: "female"},
		},
	}
	ap := &fakeApRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{
			"master.gender.male":   "Male (EN)",
			"master.gender.female": "Female (EN)",
		},
	}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	out, err := uc.ListGender(ctx, meta.ListGenderInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Tx / Repo 呼び出し
	if tx.CalledCount != 1 {
		t.Errorf("expected TxRunner.Do to be called once")
	}
	if !gn.Called {
		t.Errorf("expected GenderRepo.ListActive to be called")
	}

	// 件数
	if len(out.Items) != 2 {
		t.Fatalf("expected 2 items, got=%d", len(out.Items))
	}

	// 結果検証
	m := map[string]meta.GenderDTO{}
	for _, it := range out.Items {
		m[it.Code] = it
	}

	// male
	male := m["male"]
	if male.Label != "Male (EN)" {
		t.Errorf("expected Male (EN), got=%s", male.Label)
	}
	if male.Value != male.ID {
		t.Errorf("expected Value == ID, got=%d vs %d", male.Value, male.ID)
	}

	// female
	female := m["female"]
	if female.Label != "Female (EN)" {
		t.Errorf("expected Female (EN), got=%s", female.Label)
	}
	if female.Value != female.ID {
		t.Errorf("expected Value == ID, got=%d vs %d", female.Value, female.ID)
	}
}

// ============================================================
// フォールバック（Translate が失敗）
// ============================================================
func TestListGender_Fallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	gn := &fakeGnRepo{
		ReturnList: []entity.Gender{
			{ID: 9, Code: "other"},
		},
	}
	ap := &fakeApRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{}, // 翻訳失敗させる
	}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	out, err := uc.ListGender(ctx, meta.ListGenderInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Items) != 1 {
		t.Fatalf("expected 1 item, got=%d", len(out.Items))
	}

	item := out.Items[0]

	// fallback = Code
	if item.Label != "other" {
		t.Errorf("expected fallback label \"other\", got=%s", item.Label)
	}
	if item.Code != "other" {
		t.Errorf("unexpected code: %s", item.Code)
	}
}

// ============================================================
// Tx エラー
// ============================================================
func TestListGender_TxError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
			return errors.New("tx failed")
		},
	}
	gn := &fakeGnRepo{}
	ap := &fakeApRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	_, err := uc.ListGender(ctx, meta.ListGenderInput{
		Locale: "en",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "tx failed" {
		t.Errorf("unexpected error: %v", err)
	}
}

// ============================================================
// Repo エラー
// ============================================================
func TestListGender_RepoError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	gn := &fakeGnRepo{
		ReturnErr: errors.New("repo error"),
	}
	ap := &fakeApRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	_, err := uc.ListGender(ctx, meta.ListGenderInput{
		Locale: "en",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "repo error" {
		t.Errorf("unexpected error: %v", err)
	}
}

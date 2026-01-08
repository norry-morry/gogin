// Package meta_test は internal/usecase/meta パッケージのユニットテストを行います。
//
// 本ファイルでは ListEducationStatus の正常系・フォールバック系・異常系・ETag の性質を検証します。
//
// 依存の TxRunner / Repository / Translator は mocks_test.go の Fake を利用し、
// usecase 層のロジックのみを純粋にテストします。
package meta_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"resume/internal/domain/entity"
	"resume/internal/usecase/meta"
)

// ============================================================
// 正常系
// ============================================================

func TestListEducationStatus_Normal(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{
		ReturnList: []entity.EducationStatus{
			{ID: 1, Code: "graduated", SortOrder: 2},
			{ID: 2, Code: "enrolled", SortOrder: 1},
		},
	}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{
			"master.education_status.graduated": "Graduated",
			"master.education_status.enrolled":  "Enrolled",
		},
	}

	uc := newMetaUC(tx, ap, gn, cn, es, dt, tr)

	in := meta.ListEducationStatusInput{
		Locale: "en",
	}

	out, err := uc.ListEducationStatus(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !es.Called {
		t.Errorf("esRepo.ListIsActive should be called")
	}
	if tx.CalledCount != 1 {
		t.Errorf("TxRunner.Do should be called once")
	}

	// SortOrder: enrolled(1) → graduated(2)
	if len(out.Items) != 2 {
		t.Fatalf("expected 2 items, got=%d", len(out.Items))
	}

	gotCodes := []string{out.Items[0].Code, out.Items[1].Code}
	expCodes := []string{"enrolled", "graduated"}
	if !reflect.DeepEqual(gotCodes, expCodes) {
		t.Errorf("expected code order=%v, got=%v", expCodes, gotCodes)
	}

	// Label は Translator の結果
	if out.Items[0].Label != "Enrolled" {
		t.Errorf("label mismatch: %s", out.Items[0].Label)
	}
	if out.Items[1].Label != "Graduated" {
		t.Errorf("label mismatch: %s", out.Items[1].Label)
	}

	if out.ETag == "" {
		t.Errorf("ETag should be generated")
	}
}

// ============================================================
// フォールバック系：Translator が失敗時の Label の挙動
// ============================================================

func TestListEducationStatus_Fallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{
		ReturnList: []entity.EducationStatus{
			{ID: 1, Code: "graduated", SortOrder: 1},
		},
	}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		// 翻訳失敗させる（LabelMap 空）
		LabelMap: map[string]string{},
	}

	uc := newMetaUC(tx, ap, gn, cn, es, dt, tr)

	in := meta.ListEducationStatusInput{
		Locale: "en",
	}

	out, err := uc.ListEducationStatus(ctx, in)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if len(out.Items) != 1 {
		t.Fatalf("expected 1 item, got=%d", len(out.Items))
	}

	// 翻訳に失敗した場合、Label は Code にフォールバックする
	if out.Items[0].Label != "graduated" {
		t.Errorf("expected Label=Code fallback, got=%s", out.Items[0].Label)
	}
}

// ============================================================
// 異常系：TxRunner がエラー
// ============================================================

func TestListEducationStatus_TxRunnerError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
			return errors.New("tx failed")
		},
	}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{}

	uc := newMetaUC(tx, ap, gn, cn, es, dt, tr)

	_, err := uc.ListEducationStatus(ctx, meta.ListEducationStatusInput{
		Locale: "en",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "tx failed" {
		t.Errorf("unexpected error: %v", err)
	}
}

// ============================================================
// 異常系：Repository がエラー
// ============================================================

func TestListEducationStatus_RepoError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{
		ReturnErr: errors.New("repo error"),
	}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{}

	uc := newMetaUC(tx, ap, gn, cn, es, dt, tr)

	_, err := uc.ListEducationStatus(ctx, meta.ListEducationStatusInput{
		Locale: "en",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "repo error" {
		t.Errorf("unexpected error: %v", err)
	}
}

// ============================================================
// ETag: 内容が同じなら同じ / 順序が違っても同じ / 内容が違うと変わる
// ============================================================

func TestListEducationStatus_ETag(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{
			"master.education_status.a": "A",
			"master.education_status.b": "B",
		},
	}

	// FakeRepo を操作して順番だけ変える（コード集合は同じ）
	es1 := &fakeEsRepo{
		ReturnList: []entity.EducationStatus{
			{ID: 1, Code: "a", SortOrder: 2},
			{ID: 2, Code: "b", SortOrder: 1},
		},
	}
	es2 := &fakeEsRepo{
		ReturnList: []entity.EducationStatus{
			{ID: 2, Code: "b", SortOrder: 1},
			{ID: 1, Code: "a", SortOrder: 2},
		},
	}
	esChanged := &fakeEsRepo{
		ReturnList: []entity.EducationStatus{
			{ID: 1, Code: "a", SortOrder: 2},
			{ID: 2, Code: "b", SortOrder: 1},
			{ID: 3, Code: "c", SortOrder: 3},
		},
	}

	uc1 := newMetaUC(tx, ap, gn, cn, es1, dt, tr)
	out1, err := uc1.ListEducationStatus(ctx, meta.ListEducationStatusInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("uc1.ListEducationStatus error: %v", err)
	}

	uc2 := newMetaUC(tx, ap, gn, cn, es2, dt, tr)
	out2, err := uc2.ListEducationStatus(ctx, meta.ListEducationStatusInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("uc2.ListEducationStatus error: %v", err)
	}

	uc3 := newMetaUC(tx, ap, gn, cn, esChanged, dt, tr)
	out3, err := uc3.ListEducationStatus(ctx, meta.ListEducationStatusInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("uc3.ListEducationStatus error: %v", err)
	}

	// 1) 同じコード集合 → 同じ ETag
	if out1.ETag != out2.ETag {
		t.Errorf("expected same ETag for same code set, got %s and %s", out1.ETag, out2.ETag)
	}

	// 2) コード集合が異なる → ETag も異なる
	if out1.ETag == out3.ETag {
		t.Errorf("expected different ETag when code set changed")
	}
}

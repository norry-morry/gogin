// Package meta_test は internal/usecase/meta パッケージのユニットテストを行います。
//
// 本ファイルでは ListAddressPurpose の正常系・フォールバック系・異常系・ETag の性質を検証します。
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
// セットアップヘルパ
// ============================================================

func newMetaUC(
	tx *fakeTxRunner,
	ap *fakeApRepo,
	gn *fakeGnRepo,
	cn *fakeCnRepo,
	es *fakeEsRepo,
	dt *fakeDtRepo,
	tr *fakeTranslator,
) meta.Usecase {
	return meta.New(tx, ap, gn, cn, es, dt, tr)
}

// ============================================================
// 正常系
// ============================================================

func TestListAddressPurpose_Normal(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	ap := &fakeApRepo{
		ReturnList: []entity.AddressPurpose{
			{ID: 1, Code: "billing", DisplayName: "Billing JP", SortOrder: 2},
			{ID: 2, Code: "shipping", DisplayName: "Shipping JP", SortOrder: 1},
		},
	}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{
			"master.address_purpose.billing":  "Billing (EN)",
			"master.address_purpose.shipping": "Shipping (EN)",
		},
	}

	uc := newMetaUC(tx, ap, gn, cn, es, dt, tr)

	in := meta.ListAddressPurposeInput{
		Locale: "en",
	}

	out, err := uc.ListAddressPurpose(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ap.Called {
		t.Errorf("apRepo.GetAllActive should be called")
	}
	if tx.CalledCount != 1 {
		t.Errorf("TxRunner.Do should be called once")
	}

	// SortOrder: shipping(1) → billing(2)
	if len(out.Items) != 2 {
		t.Fatalf("expected 2 items, got=%d", len(out.Items))
	}

	gotCodes := []string{out.Items[0].Code, out.Items[1].Code}
	expCodes := []string{"shipping", "billing"}
	if !reflect.DeepEqual(gotCodes, expCodes) {
		t.Errorf("expected code order=%v, got=%v", expCodes, gotCodes)
	}

	// Label は Translator の結果
	if out.Items[0].Label != "Shipping (EN)" {
		t.Errorf("label mismatch: %s", out.Items[0].Label)
	}
	if out.Items[1].Label != "Billing (EN)" {
		t.Errorf("label mismatch: %s", out.Items[1].Label)
	}

	if out.ETag == "" {
		t.Errorf("ETag should be generated")
	}
}

// ============================================================
// フォールバック系：Translator が失敗時の Label の挙動
// ============================================================

func TestListAddressPurpose_Fallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	ap := &fakeApRepo{
		ReturnList: []entity.AddressPurpose{
			{ID: 1, Code: "billing", DisplayName: "Billing JP", SortOrder: 1},
		},
	}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		// 翻訳失敗させる（LabelMap 空）
		LabelMap: map[string]string{},
	}

	uc := newMetaUC(tx, ap, gn, cn, es, dt, tr)

	in := meta.ListAddressPurposeInput{
		Locale: "en",
	}

	out, err := uc.ListAddressPurpose(ctx, in)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if len(out.Items) != 1 {
		t.Fatalf("expected 1 item, got=%d", len(out.Items))
	}

	// フォールバックの優先順位：
	// resolve(...) 失敗 → fallback (DisplayName) → それも空 → Code
	if out.Items[0].Label != "Billing JP" {
		t.Errorf("expected Label=DisplayName fallback, got=%s", out.Items[0].Label)
	}
}

// ============================================================
// 異常系：TxRunner がエラー
// ============================================================

func TestListAddressPurpose_TxRunnerError(t *testing.T) {
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

	_, err := uc.ListAddressPurpose(ctx, meta.ListAddressPurposeInput{
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

func TestListAddressPurpose_RepoError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	ap := &fakeApRepo{
		ReturnErr: errors.New("repo error"),
	}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{}

	uc := newMetaUC(tx, ap, gn, cn, es, dt, tr)

	_, err := uc.ListAddressPurpose(ctx, meta.ListAddressPurposeInput{
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

func TestListAddressPurpose_ETag(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	gn := &fakeGnRepo{}
	cn := &fakeCnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{
			"master.address_purpose.a": "A",
			"master.address_purpose.b": "B",
		},
	}

	// FakeRepo を操作して順番だけ変える（コード集合は同じ）
	ap1 := &fakeApRepo{
		ReturnList: []entity.AddressPurpose{
			{ID: 1, Code: "a", DisplayName: "AA", SortOrder: 2},
			{ID: 2, Code: "b", DisplayName: "BB", SortOrder: 1},
		},
	}
	ap2 := &fakeApRepo{
		ReturnList: []entity.AddressPurpose{
			{ID: 2, Code: "b", DisplayName: "BB", SortOrder: 1},
			{ID: 1, Code: "a", DisplayName: "AA", SortOrder: 2},
		},
	}
	apChanged := &fakeApRepo{
		ReturnList: []entity.AddressPurpose{
			{ID: 1, Code: "a", DisplayName: "AA", SortOrder: 2},
			{ID: 2, Code: "b", DisplayName: "BB", SortOrder: 1},
			{ID: 3, Code: "c", DisplayName: "CC", SortOrder: 3},
		},
	}

	uc1 := newMetaUC(tx, ap1, gn, cn, es, dt, tr)
	out1, err := uc1.ListAddressPurpose(ctx, meta.ListAddressPurposeInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("uc1.ListAddressPurpose error: %v", err)
	}

	uc2 := newMetaUC(tx, ap2, gn, cn, es, dt, tr)
	out2, err := uc2.ListAddressPurpose(ctx, meta.ListAddressPurposeInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("uc2.ListAddressPurpose error: %v", err)
	}

	uc3 := newMetaUC(tx, apChanged, gn, cn, es, dt, tr)
	out3, err := uc3.ListAddressPurpose(ctx, meta.ListAddressPurposeInput{
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("uc3.ListAddressPurpose error: %v", err)
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

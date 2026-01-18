// Package meta_test は internal/usecase/meta のユニットテストを行います。
// 本ファイルは ListCountry インタラクタのユニットテストです。
//
// テスト範囲：
// - 正常系：Tx → CountryRepo → Translator → assemble までの流れ
// - フォールバック：翻訳失敗時は Code をそのままラベルとして使う
// - 異常系：TxRunner の失敗、Repository の失敗
// - ETag：同じ入力なら同じ ETag、内容が変われば ETag も変わる
//
// mocks_test.go の Fake 実装を用いて、usecase のビジネスロジックのみを検証する。
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

// TestListCountry_Normal は、Country 一覧が正常に取得され、
// Translator によりラベルが付与されることを確認します。
func TestListCountry_Normal(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	cn := &fakeCnRepo{
		// リポジトリが返す順序を、そのままレスポンスに反映する前提。
		ReturnList: []entity.Country{
			{Code: "jp", IsSupported: true, SortOrder: 2},
			{Code: "us", IsSupported: true, SortOrder: 1},
		},
	}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{
			"master.country.jp": "Japan (EN)",
			"master.country.us": "United States (EN)",
		},
	}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	in := meta.ListCountryInput{
		Locale:        "en",
		OnlySupported: false,
	}

	out, err := uc.ListCountry(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cn.Called {
		t.Errorf("expected CountryRepo.List to be called")
	}
	if tx.CalledCount != 1 {
		t.Errorf("expected TxRunner.Do to be called once, got %d", tx.CalledCount)
	}

	if len(out.Items) != 2 {
		t.Fatalf("expected 2 items, got=%d", len(out.Items))
	}

	// リポジトリの順序をそのまま返していること
	if out.Items[0].Value != "jp" || out.Items[1].Value != "us" {
		t.Errorf("unexpected value order: got=[%s %s], want=[jp us]",
			out.Items[0].Value, out.Items[1].Value)
	}

	// 翻訳が効いていること
	if out.Items[0].Label != "Japan (EN)" {
		t.Errorf("unexpected label[0]: got=%s, want=Japan (EN)", out.Items[0].Label)
	}
	if out.Items[1].Label != "United States (EN)" {
		t.Errorf("unexpected label[1]: got=%s, want=United States (EN)", out.Items[1].Label)
	}

	if out.ETag == "" {
		t.Errorf("ETag should not be empty")
	}
}

// ============================================================
// フォールバック確認（Translate が失敗した場合）
// ============================================================

// TestListCountry_Fallback は Translator でキーが見つからなかった場合に、
// Country.Code がそのままラベルとして使われることを確認します。
func TestListCountry_Fallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	cn := &fakeCnRepo{
		ReturnList: []entity.Country{
			{Code: "jp", IsSupported: true, SortOrder: 1},
		},
	}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		// LabelMap を空にして、常に Translate が失敗するようにする。
		LabelMap: map[string]string{},
	}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	out, err := uc.ListCountry(ctx, meta.ListCountryInput{
		Locale:        "en",
		OnlySupported: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Items) != 1 {
		t.Fatalf("expected 1 item, got=%d", len(out.Items))
	}

	// fallback: Code をそのままラベルに使う
	if out.Items[0].Label != "jp" {
		t.Errorf("expected fallback label `jp`, got=%s", out.Items[0].Label)
	}
}

// ============================================================
// 異常系：TxRunner 失敗
// ============================================================

// TestListCountry_TxError は TxRunner.Do がエラーを返した場合、
// そのエラーがそのまま呼び出し元へ伝播することを確認します。
func TestListCountry_TxError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{
		DoFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
			return errors.New("tx failed")
		},
	}
	cn := &fakeCnRepo{}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	_, err := uc.ListCountry(ctx, meta.ListCountryInput{
		Locale:        "en",
		OnlySupported: false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "tx failed" {
		t.Errorf("unexpected tx error: %v", err)
	}
}

// ============================================================
// 異常系：Repository 失敗
// ============================================================

// TestListCountry_RepoError は CountryRepo.List がエラーを返した場合、
// そのエラーがそのまま呼び出し元へ伝播することを確認します。
func TestListCountry_RepoError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tx := &fakeTxRunner{}
	cn := &fakeCnRepo{
		ReturnErr: errors.New("repo error"),
	}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{}

	uc := meta.New(tx, ap, gn, cn, es, dt, tr)

	_, err := uc.ListCountry(ctx, meta.ListCountryInput{
		Locale:        "en",
		OnlySupported: false,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "repo error" {
		t.Errorf("unexpected repo error: %v", err)
	}
}

// ============================================================
// ETag：同じ内容 → 同じ ETag / 異なる内容 → 異なる ETag
// ============================================================

// TestListCountry_ETag は、同じ Country 一覧を返す場合は同一の ETag が生成され、
// 一覧の内容が変わると ETag も変わることを確認します。
func TestListCountry_ETag(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tx := &fakeTxRunner{}
	ap := &fakeApRepo{}
	gn := &fakeGnRepo{}
	es := &fakeEsRepo{}
	dt := &fakeDtRepo{}
	tr := &fakeTranslator{
		LabelMap: map[string]string{
			"master.country.jp": "Japan",
			"master.country.us": "USA",
			"master.country.fr": "France",
		},
	}

	// 同じ内容（順序も同じ）
	cn1 := &fakeCnRepo{
		ReturnList: []entity.Country{
			{Code: "jp", IsSupported: true, SortOrder: 2},
			{Code: "us", IsSupported: true, SortOrder: 1},
		},
	}
	cn2 := &fakeCnRepo{
		ReturnList: []entity.Country{
			{Code: "jp", IsSupported: true, SortOrder: 2},
			{Code: "us", IsSupported: true, SortOrder: 1},
		},
	}
	// 内容が異なる（fr が追加されている）
	cnChanged := &fakeCnRepo{
		ReturnList: []entity.Country{
			{Code: "jp", IsSupported: true, SortOrder: 2},
			{Code: "us", IsSupported: true, SortOrder: 1},
			{Code: "fr", IsSupported: false, SortOrder: 3},
		},
	}

	uc1 := meta.New(tx, ap, gn, cn1, es, dt, tr)
	out1, err := uc1.ListCountry(ctx, meta.ListCountryInput{Locale: "en"})
	if err != nil {
		t.Fatalf("unexpected error on uc1: %v", err)
	}

	uc2 := meta.New(tx, ap, gn, cn2, es, dt, tr)
	out2, err := uc2.ListCountry(ctx, meta.ListCountryInput{Locale: "en"})
	if err != nil {
		t.Fatalf("unexpected error on uc2: %v", err)
	}

	uc3 := meta.New(tx, ap, gn, cnChanged, es, dt, tr)
	out3, err := uc3.ListCountry(ctx, meta.ListCountryInput{Locale: "en"})
	if err != nil {
		t.Fatalf("unexpected error on uc3: %v", err)
	}

	// 1) 同じ一覧 → 同じ ETag
	if out1.ETag != out2.ETag {
		t.Errorf("expected same ETag, got %s vs %s", out1.ETag, out2.ETag)
	}

	// 2) 一覧の内容が変わる → ETag も変わる
	if out1.ETag == out3.ETag {
		t.Errorf("expected different ETag, got %s vs %s", out1.ETag, out3.ETag)
	}
}

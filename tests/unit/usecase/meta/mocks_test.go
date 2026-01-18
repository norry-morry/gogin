// Package meta_test は usecase/meta のユニットテスト用偽実装をまとめたものです。
// Repository / TxRunner / Translator をモック化し、usecase のロジックのみを検証します。
package meta_test

import (
	"context"
	"strings"
	"sync"

	"resume/internal/domain/entity"
	vo "resume/internal/domain/valueobject/i18n"
	uci18n "resume/internal/usecase/i18n"
)

// ------------------------------------------------------------
// FakeTxRunner
// ------------------------------------------------------------

// fakeTxRunner は TxRunner をモック化した実装です。
// Do が呼ばれた回数を記録し、DoFunc で挙動を差し替えられます。
type fakeTxRunner struct {
	CalledCount int
	DoFunc      func(ctx context.Context, fn func(ctx context.Context) error) error
	mu          sync.Mutex
}

// Do は TxRunner.Do のモック実装です。
func (f *fakeTxRunner) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	f.mu.Lock()
	f.CalledCount++
	f.mu.Unlock()

	if f.DoFunc != nil {
		return f.DoFunc(ctx, fn)
	}
	return fn(ctx)
}

// ------------------------------------------------------------
// Fake AddressPurpose Repository
// ------------------------------------------------------------

// fakeApRepo は AddressPurpose 用 Repository のモックです。
type fakeApRepo struct {
	Called     bool
	ReturnList []entity.AddressPurpose
	ReturnErr  error
}

// GetAllActive は AddressPurposeRepo.GetAllActive のモック実装です。
func (f *fakeApRepo) GetAllActive(ctx context.Context) ([]entity.AddressPurpose, error) {
	f.Called = true
	if f.ReturnErr != nil {
		return nil, f.ReturnErr
	}
	return f.ReturnList, nil
}

// ------------------------------------------------------------
// Fake Gender Repository
// ------------------------------------------------------------

// fakeGnRepo は Gender 用 Repository のモックです。
type fakeGnRepo struct {
	Called     bool
	ReturnList []entity.Gender
	ReturnErr  error
}

// ListActive は GenderRepo.ListActive のモック実装です。
func (f *fakeGnRepo) ListActive(ctx context.Context) ([]entity.Gender, error) {
	f.Called = true
	if f.ReturnErr != nil {
		return nil, f.ReturnErr
	}
	return f.ReturnList, nil
}

// ------------------------------------------------------------
// Fake Country Repository
// ------------------------------------------------------------

// fakeCnRepo は Country 用 Repository のモックです。
// 本番側の CnRepo は List(ctx, onlySupported bool) を持っている想定。
type fakeCnRepo struct {
	Called            bool
	LastOnlySupported *bool
	ReturnList        []entity.Country
	ReturnErr         error
}

// List は CnRepo.List のモック実装です。
func (f *fakeCnRepo) List(ctx context.Context, onlySupported bool) ([]entity.Country, error) {
	f.Called = true
	v := onlySupported
	f.LastOnlySupported = &v

	if f.ReturnErr != nil {
		return nil, f.ReturnErr
	}
	return f.ReturnList, nil
}

// ------------------------------------------------------------
// Fake Translator
// ------------------------------------------------------------

// TranslateCall は Translate の呼び出しを記録するための構造体です。
type TranslateCall struct {
	Locale vo.Locale
	Key    vo.Key
}

// BundleCall は Bundle の呼び出しを記録するための構造体です。
type BundleCall struct {
	Locale vo.Locale
	Prefix string
}

// fakeTranslator は uci18n.Translator をモック化した実装です。
// LabelMap のキーを vo.Key の中身（string）として扱います。
type fakeTranslator struct {
	LabelMap      map[string]string
	TranslateLogs []TranslateCall
	BundleLogs    []BundleCall
}

// Translate は uci18n.Translator.Translate のモック実装です。
func (f *fakeTranslator) Translate(loc vo.Locale, key vo.Key) (string, bool) {
	f.TranslateLogs = append(f.TranslateLogs, TranslateCall{
		Locale: loc,
		Key:    key,
	})

	if f.LabelMap == nil {
		return "", false
	}

	// vo.Key → string は key.String()
	if v, ok := f.LabelMap[key.String()]; ok {
		return v, true
	}
	return "", false
}

// Bundle は uci18n.Translator.Bundle のモック実装です。
// 現状 meta.Usecase からは呼ばれていない想定なので、最低限の実装にしています。
func (f *fakeTranslator) Bundle(loc vo.Locale, prefix string) map[string]string {
	f.BundleLogs = append(f.BundleLogs, BundleCall{
		Locale: loc,
		Prefix: prefix,
	})

	if f.LabelMap == nil {
		return nil
	}

	out := map[string]string{}
	for k, v := range f.LabelMap {
		if strings.HasPrefix(k, prefix) {
			out[k] = v
		}
	}
	return out
}

// ... existing code ...
// ------------------------------------------------------------
// Fake DegreeType Repository
// ------------------------------------------------------------

// fakeDtRepo は DegreeType 用 Repository のモックです。
type fakeDtRepo struct {
	Called     bool
	ReturnList []entity.DegreeType
	ReturnErr  error
}

// ListIsActive は DegreeTypeRepository.ListIsActive のモック実装です。
func (f *fakeDtRepo) ListIsActive(ctx context.Context) ([]entity.DegreeType, error) {
	f.Called = true
	if f.ReturnErr != nil {
		return nil, f.ReturnErr
	}
	return f.ReturnList, nil
}

// ------------------------------------------------------------
// Fake EducationStatus Repository
// ------------------------------------------------------------

// fakeEsRepo は EducationStatus 用 Repository のモックです。
type fakeEsRepo struct {
	Called     bool
	ReturnList []entity.EducationStatus
	ReturnErr  error
}

// ListIsActive は EducationStatusRepository.ListIsActive のモック実装です。
func (f *fakeEsRepo) ListIsActive(ctx context.Context) ([]entity.EducationStatus, error) {
	f.Called = true
	if f.ReturnErr != nil {
		return nil, f.ReturnErr
	}
	return f.ReturnList, nil
}

// コンパイル時に fakeTranslator が uci18n.Translator を実装していることを確認する。
var _ uci18n.Translator = (*fakeTranslator)(nil)

// Package router は、HTTPルーティングおよび関連アダプタを提供します。
package router

import uci18n "resume/internal/usecase/i18n"

// usecaseLocaleSource は、i18n Usecase から利用可能ロケール一覧を取得するための
// アダプタ構造体です。middleware.Opts.Store に渡して利用します。
type usecaseLocaleSource struct {
	uc uci18n.Usecase
}

// NewUsecaseLocaleSource は、i18n Usecase をラップして LocaleSource として返します。
func NewUsecaseLocaleSource(uc uci18n.Usecase) usecaseLocaleSource {
	return usecaseLocaleSource{uc: uc}
}

// AvailableLocaleCodes は、Usecase 層の ListLocales() の結果から言語コード一覧を返します。
// middleware.Opts.Store インターフェースに適合させるために使用されます。
func (s usecaseLocaleSource) AvailableLocaleCodes() []string {
	// usecase 層の ListLocales() は、locale 情報（code, nameなど）を返す前提
	locales := s.uc.ListLocales()
	codes := make([]string, 0, len(locales))
	for _, l := range locales {
		codes = append(codes, l.Code())
	}
	return codes
}

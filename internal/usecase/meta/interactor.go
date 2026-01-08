// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

import "resume/internal/usecase/i18n"

// Interactor は、メタ情報（国・住所用途・性別など）に関するユースケースの実装です。
type Interactor struct {
	tx         TxRunner
	apRepo     ApRepo
	gnRepo     GnRepo
	cnRepo     CnRepo
	esRepo     EsRepo
	dtRepo     DtRepo
	translator i18n.Translator
}

// New は、トランザクションランナーと各種リポジトリ、翻訳器を受け取り、
// メタ情報ユースケースの実装である Interactor を生成します。
func New(
	tx TxRunner,
	apRepo ApRepo,
	gnRepo GnRepo,
	cnRepo CnRepo,
	esRepo EsRepo,
	dtRepo DtRepo,
	translator i18n.Translator,
) Usecase {
	return &Interactor{
		tx:         tx,
		apRepo:     apRepo,
		gnRepo:     gnRepo,
		cnRepo:     cnRepo,
		esRepo:     esRepo,
		dtRepo:     dtRepo,
		translator: translator,
	}
}

// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

// Interactor は Education に関するユースケース処置の実装を提供する
// リポジトリ層及びトランザクション実行を通じて、ドメインロジックを仲介します
type Interactor struct {
	tx     TxRunner
	ueRepo UeRepo
	esRepo EsRepo
	dtRepo DtRepo
}

// New は 学歴 用ユースケースの Interactor を生成して返す
// 各リポジトリとトランザクションランナーを依存として受け取る
func New(
	tx TxRunner,
	ueRepo UeRepo,
	esRepo EsRepo,
	dtRepo DtRepo,
) Usecase {
	return &Interactor{
		tx:     tx,
		ueRepo: ueRepo,
		esRepo: esRepo,
		dtRepo: dtRepo,
	}
}

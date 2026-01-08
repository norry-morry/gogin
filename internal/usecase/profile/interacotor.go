package profile

// Interactor は、Profile に関するユースケース処理の実装を提供します。
// リポジトリ層およびトランザクション実行を通じて、ドメインロジックを仲介します。
type Interactor struct {
	uRepo  URepo
	uIRepo UIRepo
	upRepo UPRepo
	uaRepo UARepo
	tx     TxRunner
}

// New は、Profile 用ユースケースの Interactor を生成して返します。
// 各リポジトリとトランザクションランナーを依存として受け取ります。
func New(
	uRepo URepo,
	uIRepo UIRepo,
	upRepo UPRepo,
	uaRepo UARepo,
	tx TxRunner,
) Usecase {
	return &Interactor{
		uRepo:  uRepo,
		uIRepo: uIRepo,
		upRepo: upRepo,
		uaRepo: uaRepo,
		tx:     tx,
	}
}

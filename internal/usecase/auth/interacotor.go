package auth

type interactor struct {
	users  UserRepository
	idents AuthIdentityRepository
	tx     TxRunner
	clock  Clock
}

// New は認証ユースケースの実装（interactor）を生成して返します。
// tx または clock が nil の場合、内部で noTxRunner / sysClock にフォールバックします
func New(users UserRepository, idents AuthIdentityRepository, tx TxRunner, clock Clock) Usecase {
	if clock == nil {
		clock = sysClock{}
	}
	return &interactor{users: users, idents: idents, tx: tx, clock: clock}
}

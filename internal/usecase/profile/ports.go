package profile

import (
	"resume/internal/domain/repository"
	stx "resume/internal/shared/tx"
)

// URepo は User ドメインに関する永続化処理を行うリポジトリのエイリアスです。
// repository.UserRepository を短縮形で参照するための型定義です。
type URepo = repository.UserRepository

// UIRepo は Firebase などの外部 ID 情報を扱う UserIdentity リポジトリのエイリアスです。
// repository.UserIdentityRepository を短縮形で参照するための型定義です。
type UIRepo = repository.UserIdentityRepository

// TxRunner はトランザクション制御を提供するインターフェースのエイリアスです。
// stx.Runner を短縮形で参照し、ユースケース層での依存注入を簡略化します。
type TxRunner = stx.Runner

// UPRepo はユーザープロフィール情報を扱うリポジトリインターフェースのエイリアスです。
// repository.UserProfileRepository を短縮形で参照するために定義されています。
// 実際のデータ取得（JOIN 含む）はインフラ層の実装に委譲されます。
type UPRepo = repository.UserProfileRepository

// UARepo は ユーザー住所の情報を扱うリポジトリインターフェースのエイリアス
// repository.UserAddressRepository を 短縮形で参照するために定義
// 実際のデータ取得(joinも含めて)は インフラ層の実装に移譲
type UARepo = repository.UserAddressRepository

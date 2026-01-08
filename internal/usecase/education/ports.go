// Package education は、学歴に関するユースケース（アプリケーションロジック）を提供します。
// ドメイン層のエンティティやリポジトリを操作し、アプリケーション全体で再利用可能なビジネスフローを実装します。
package education

import (
	"resume/internal/domain/repository"
	stx "resume/internal/shared/tx"
)

// TxRunner はトランザクション制御を提供するインターフェースのエイリアスです。
// stx.Runner を短縮形で参照し、ユースケース層での依存注入を簡略化します。
type TxRunner = stx.Runner

// UeRepo は 学歴 リポジトリのエイリアスです。
// repository.UserEducationRepository を短縮形で参照するための型定義です。
type UeRepo = repository.UserEducationRepository

// EsRepo は 学歴状態 リポジトリのエイリアスです。
// repository.EducationStatusRepository を短縮形で参照するための型定義です。
type EsRepo = repository.EducationStatusRepository

// DtRepo は 学位種別 リポジトリのエイリアスです。
// repository.DegreeTypeRepository を短縮形で参照するための型定義です。
type DtRepo = repository.DegreeTypeRepository

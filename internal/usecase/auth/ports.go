// Package auth は、ユースケースが依存するポート（インターフェース）と
// ドメインエンティティの型エイリアスを定義します。
package auth

import (
	"time"

	"resume/internal/domain/entity"
	"resume/internal/domain/repository"
	stx "resume/internal/shared/tx"
)

// UserRepository is the repository interface for users table used by the auth use case.
// It aliases the domain repository implementation so the use case depends on abstractions.
type UserRepository = repository.UserRepository

// AuthIdentityRepository is the repository interface for auth_identities table used by the auth use case.
// It aliases the domain repository implementation so the use case depends on abstractions.
// revive:disable:exported  // パッケージ名と併せた命名の都合上、stutter の警告を抑制
type AuthIdentityRepository = repository.UserIdentityRepository

// revive:enable:exported

// User is the domain user entity the auth use case works with.
// It aliases entity.User to avoid leaking infra details into the use case.
type User = entity.User

// Identity is the domain auth identity entity the auth use case works with.
// It aliases entity.UserIdentity to avoid leaking infra details into the use case.
type Identity = entity.UserIdentity

// TxRunner defines a transaction boundary abstraction used by the use case layer.
// The actual implementation (e.g., GormTxRunner) is provided by the infra layer via DI.
type TxRunner = stx.Runner

// Clock provides current time (preferably UTC) for deterministic testing and clearer ownership of "now".
type Clock interface{ Now() time.Time }

// sysClock is the default Clock implementation returning time.Now().UTC().
type sysClock struct{}

// Now returns the current time in UTC.
func (sysClock) Now() time.Time { return time.Now().UTC() }

// Package repository は、ドメイン層における永続化のための
// リポジトリインターフェースを定義します。
// このファイルでは、国情報を扱うリポジトリの契約を定義します。
package repository

import (
	"context"

	"resume/internal/domain/entity"
)

// CountryRepository は、国情報の取得に関するリポジトリの契約を定義します。
// データソース（データベースや静的ファイルなど）から国一覧を取得するために使用します。
type CountryRepository interface {
	// List は、利用可能なすべての国情報を返します。
	List(ctx context.Context, onlySupported bool) ([]entity.Country, error)
}

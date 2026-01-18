// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

import (
	"resume/internal/domain/repository"
	stx "resume/internal/shared/tx"
)

// TxRunner はユースケース内でトランザクションを張るためのポートです。
type TxRunner = stx.Runner

// ApRepo は住所目的の読み書きを行うリポジトリポートです。
type ApRepo = repository.AddressPurposeRepository

// GnRepo は 性別の読み書きを行うリポジトリポートです
type GnRepo = repository.GenderRepository

// CnRepo は 国の読み書きを行うリポジトリポートです
type CnRepo = repository.CountryRepository

// EsRepo は 学歴状態の読み書きを行うリポジトリポートです
type EsRepo = repository.EducationStatusRepository

// DtRepo は 学位種別の読み書きを行うリポジトリポートです
type DtRepo = repository.DegreeTypeRepository

// Package meta は、国・住所用途・性別などのメタ情報を扱うユースケースを提供します。
package meta

import "context"

// Usecase は、メタ情報（国・住所用途・性別など）を取得するユースケースのインターフェースです。
type Usecase interface {
	// ListAddressPurpose は、住所用途の選択肢一覧を取得します。
	ListAddressPurpose(ctx context.Context, in ListAddressPurposeInput) (ListAddressPurposeOutput, error)

	// ListGender は、性別の選択肢一覧を取得します。
	ListGender(ctx context.Context, in ListGenderInput) (ListGenderOutput, error)

	// ListCountry は、国の選択肢一覧を取得します。
	ListCountry(ctx context.Context, in ListCountryInput) (ListCountryOutput, error)

	// ListEducationStatus は 学歴状態一覧を取得します
	ListEducationStatus(ctx context.Context, in ListEducationStatusInput) (ListEducationStatusOutput, error)

	// ListDegreeType は 学位種別一覧を取得します
	ListDegreeType(ctx context.Context, in ListDegreeTypeInput) (ListDegreeTypeOutput, error)
}

// Package auth は、entity↔DTO の詰め替え（アセンブラ）を提供します。
package auth

import "resume/internal/shared/util"

// NOTE: entity は domain/entity の型（User, Identity）を使用します。

func toUserOutput(u *User, ps []*Identity) UserOutput {
	out := UserOutput{
		ID:            uint64(u.ID),
		UID:           u.UID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		DisplayName:   u.DisplayName,
		PhotoURL:      u.PhotoURL,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Providers:     make([]ProviderOutput, 0, len(ps)),
	}
	for _, p := range ps {
		out.Providers = append(out.Providers, ProviderOutput{
			ID:                  uint64(p.ID),
			Provider:            p.Provider,
			ProviderUserID:      p.ProviderUserID,
			ProviderDisplayName: p.ProviderDisplayName,
			EmailAtSignup:       util.Clone(p.EmailAtSignup),
			LinkedAt:            p.CreatedAt,
		})
	}
	return out
}

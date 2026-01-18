// Package inmemory はアプリケーション内キャッシュ層の Gateway 実装を提供します。
// i18n辞書キャッシュなどの Translator/CacheStore 実装を含みます。
package inmemory

import (
	"log/slog"
	"sync/atomic"

	vo "resume/internal/domain/valueobject/i18n"
	usecase "resume/internal/usecase/i18n"
)

// snapshot は全ロケール分の辞書バンドルを保持します。
type snapshot struct {
	ByLocale map[vo.Locale]*vo.Bundle
}

// CacheStoreImpl は usecase 層の Translator および CacheStore の実装です。
// atomic.Value に *snapshot を格納し、ロックレスな読み取りを実現します。
type CacheStoreImpl struct {
	val atomic.Value // stores *snapshot
}

// NewCacheStoreImpl は CacheStoreImpl を初期化して返します。
func NewCacheStoreImpl() *CacheStoreImpl {
	cs := &CacheStoreImpl{}
	cs.val.Store(&snapshot{ByLocale: map[vo.Locale]*vo.Bundle{}})
	return cs
}

// load は現在のスナップショットを取得します。
func (c *CacheStoreImpl) load() *snapshot {
	v, ok := c.val.Load().(*snapshot)
	if !ok {
		// 型が違う場合や初期化前の状態を安全に処理
		return nil
	}
	return v
}

// SwapAll は全ロケール分の辞書データをキャッシュに一括適用します。
func (c *CacheStoreImpl) SwapAll(all map[vo.Locale]*vo.Bundle) error {
	next := &snapshot{ByLocale: make(map[vo.Locale]*vo.Bundle, len(all))}
	for loc, b := range all {
		next.ByLocale[loc] = b
	}
	c.val.Store(next)
	return nil
}

// Get は指定ロケールとキーに対応する翻訳文字列を現在のスナップショットから取得します。
// フォールバックは行わず、生のキャッシュ内容のみを参照します。
func (c *CacheStoreImpl) Get(loc vo.Locale, key vo.Key) (string, bool) {
	snap := c.load()
	if b, ok := snap.ByLocale[loc]; ok {
		return b.Get(key)
	}
	return "", false
}

// Set は指定ロケール・キーの値をキャッシュ上で上書きします。
// 高頻度更新を想定していないため、スナップショットをコピーして丸ごと差し替えます。
func (c *CacheStoreImpl) Set(loc vo.Locale, key vo.Key, value string) {
	old := c.load()

	// まず snapshot 自体をコピー
	next := &snapshot{ByLocale: make(map[vo.Locale]*vo.Bundle, len(old.ByLocale))}
	for l, b := range old.ByLocale {
		next.ByLocale[l] = b
	}

	// 対象ロケールの bundle をコピーして値を差し替え
	if b, ok := next.ByLocale[loc]; ok {
		msgs := b.Messages() // map[string]string をコピーで取得
		msgs[key.Value()] = value
		next.ByLocale[loc] = vo.NewBundle(loc, msgs)
	} else {
		next.ByLocale[loc] = vo.NewBundle(loc, map[string]string{
			key.Value(): value,
		})
	}

	c.val.Store(next)
}

// Clear は全ロケール分の辞書キャッシュをクリアします。
func (c *CacheStoreImpl) Clear() {
	c.val.Store(&snapshot{ByLocale: map[vo.Locale]*vo.Bundle{}})
}

// Translate は指定ロケールとキーに対応する翻訳文字列を返します。
// 存在しない場合は英語ロケールをフォールバックとして参照します。
func (c *CacheStoreImpl) Translate(loc vo.Locale, key vo.Key) (string, bool) {
	slog.Debug("i18n.Translate",
		"loc", loc.Code(),
		"key", key.Value(),
	)
	snap := c.load()
	if b, ok := snap.ByLocale[loc]; ok {
		if v, ok := b.Get(key); ok {
			return v, true
		}
	}
	if b, ok := snap.ByLocale[vo.NewLocale("en")]; ok {
		if v, ok := b.Get(key); ok {
			return v, true
		}
	}
	slog.Debug("i18n.Translate not found",
		"loc", loc.Code(),
		"key", key.Value(),
	)
	return "", false
}

// Bundle は prefix 前方一致でキーを抽出して返します。
func (c *CacheStoreImpl) Bundle(loc vo.Locale, prefix string) map[string]string {
	snap := c.load()
	out := map[string]string{}
	if b, ok := snap.ByLocale[loc]; ok {
		for k, v := range b.Messages() {
			if len(prefix) == 0 || (len(k) >= len(prefix) && k[:len(prefix)] == prefix) {
				out[k] = v
			}
		}
	}
	return out
}

// Locales は現在キャッシュされているロケール一覧を返します。
func (c *CacheStoreImpl) Locales() []vo.Locale {
	snap := c.load()
	out := make([]vo.Locale, 0, len(snap.ByLocale))
	for loc := range snap.ByLocale {
		out = append(out, loc)
	}
	return out
}

// インターフェース適合保証
var (
	_ usecase.Translator = (*CacheStoreImpl)(nil)
	_ usecase.CacheStore = (*CacheStoreImpl)(nil)
)

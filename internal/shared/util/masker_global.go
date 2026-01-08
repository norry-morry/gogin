package util

import (
	"sync/atomic"
)

// maskerHolder は常に同じ具体型で atomic.Value に載せるためのラッパ。
type maskerHolder struct {
	M Masker
}

var globalMasker atomic.Value // 中身は常に maskerHolder

func init() {
	// 初期値: Noop
	globalMasker.Store(maskerHolder{M: NoopMasker{}})
}

// SetGlobalMasker はグローバル Masker を差し替える。
// ラッパ型に包むことで atomic.Value の具体型が不変になり、panic を防ぐ。
func SetGlobalMasker(m Masker) {
	if m == nil {
		return
	}
	globalMasker.Store(maskerHolder{M: m})
}

// GlobalMasker は現在のグローバル Masker を返す。
func GlobalMasker() Masker {
	v := globalMasker.Load()
	if v == nil {
		return NoopMasker{}
	}
	h, ok := v.(maskerHolder)
	if !ok || h.M == nil {
		return NoopMasker{}
	}
	return h.M
}

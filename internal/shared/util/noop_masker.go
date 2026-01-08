package util

// NoopMasker は何もしないマスカー。デフォルトで安全に使える。
type NoopMasker struct{}

// MaskByKey は、指定されたキーに対して値をマスクします。
// NoopMasker では実際のマスク処理は行わず、入力値 v をそのまま返します。
func (NoopMasker) MaskByKey(_ string, v any) any { return v }

// MaskMapShallow は、map[string]any の1階層目のキーを対象にマスクを行います。
// NoopMasker ではマスク処理を行わず、入力マップ src をそのまま返します。
func (NoopMasker) MaskMapShallow(src map[string]any) map[string]any { return src }

// MaskAnyRecursive は、任意の値を再帰的にマスクします。
// NoopMasker ではマスク処理を行わず、入力値 v をそのまま返します。
func (NoopMasker) MaskAnyRecursive(_ string, v any) any { return v }

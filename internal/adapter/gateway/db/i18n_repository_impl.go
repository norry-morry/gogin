// Package db は 永続化層（データベースやファイル等）にアクセスするGateway 実装を提供します。
package db

import (
	"os"
	"path/filepath"
	"strings"

	repo "resume/internal/domain/repository"
	vo "resume/internal/domain/valueobject/i18n"
	"resume/internal/shared/util"
)

// I18nRepositoryImpl は DictionaryRepository のファイルシステム実装です。
// YAML ファイル群から全ロケールの辞書データをロードします。
type I18nRepositoryImpl struct {
	Root string // 辞書ルートディレクトリ（例: "旧/i18n_localesパス"）
}

// NewI18nRepositoryImpl は I18nRepositoryImpl を生成します。
func NewI18nRepositoryImpl(root string) *I18nRepositoryImpl {
	return &I18nRepositoryImpl{Root: root}
}

// LoadAll は全ロケール分の YAML ファイルを走査し、
// flatten 済みの map を持つ vo.Bundle にまとめて返します。
func (r *I18nRepositoryImpl) LoadAll() (map[vo.Locale]*vo.Bundle, error) {
	tree := make(map[string]map[string]any) // locale -> nested map

	err := filepath.Walk(r.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		locale := util.DetectLocale(path)
		y, err := util.ReadYAMLFile(path)
		if err != nil {
			return err
		}

		// locales/ja/ui/page/home.yaml → ["ui","page","home"]
		rel, err := filepath.Rel(r.Root, path)
		if err != nil {
			// 想定外のパスなら絶対パスで代替
			rel = path
		}
		segs := strings.Split(rel, string(os.PathSeparator))
		if len(segs) < 2 {
			return nil
		}
		prefix := segs[1:]
		last := prefix[len(prefix)-1]
		prefix[len(prefix)-1] = strings.TrimSuffix(last, ext)

		nested := util.WrapWithPath(prefix, y)

		if _, ok := tree[locale]; !ok {
			tree[locale] = map[string]any{}
		}
		if err := util.DeepMerge(tree[locale], nested); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := make(map[vo.Locale]*vo.Bundle, len(tree))
	for loc, nested := range tree {
		flat := map[string]any{}
		util.Flatten("", nested, flat)

		msgs := map[string]string{}
		for k, v := range flat {
			if s, ok := v.(string); ok {
				msgs[k] = s
			}
		}

		l := vo.NewLocale(loc)
		out[l] = vo.NewBundle(l, msgs)
	}
	return out, nil
}

// インターフェース適合保証
var _ repo.DictionaryRepository = (*I18nRepositoryImpl)(nil)

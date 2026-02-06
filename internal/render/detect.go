package render

import (
	"path/filepath"
	"strings"
)

var markdownExt = map[string]struct{}{
	".md":       {},
	".markdown": {},
	".mdown":    {},
	".mkd":      {},
}

func DetectKind(name string) Kind {
	ext := strings.ToLower(filepath.Ext(name))
	if _, ok := markdownExt[ext]; ok {
		return KindMarkdown
	}
	if ext == "" {
		return KindPlain
	}
	return KindCode
}

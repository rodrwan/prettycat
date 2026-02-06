package render

import (
	"path/filepath"
	"regexp"
	"strings"
)

var languageKeywords = map[string][]string{
	".go":   {"package", "import", "func", "return", "if", "else", "for", "range", "struct", "interface", "var", "const", "type"},
	".js":   {"function", "return", "const", "let", "var", "if", "else", "for", "class", "new", "import", "export"},
	".ts":   {"function", "return", "const", "let", "if", "else", "for", "class", "interface", "type", "import", "export"},
	".py":   {"def", "return", "if", "elif", "else", "for", "while", "class", "import", "from", "try", "except"},
	".rb":   {"def", "end", "class", "module", "if", "elsif", "else", "do", "begin", "rescue", "require"},
	".java": {"class", "interface", "public", "private", "protected", "static", "void", "return", "if", "else", "for", "new"},
}

func renderCode(name string, in []byte, opts Options) (string, error) {
	plain := renderPlain(in)
	if !opts.Color {
		return plain, nil
	}

	ext := strings.ToLower(filepath.Ext(name))
	keywords := languageKeywords[ext]
	if len(keywords) == 0 {
		return colorizeGenericCode(plain), nil
	}

	out := plain
	for _, kw := range keywords {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(kw) + `\b`)
		out = re.ReplaceAllString(out, "\x1b[38;5;81m"+kw+"\x1b[0m")
	}
	out = colorizeStrings(out)
	out = colorizeComments(ext, out)
	return out, nil
}

func colorizeGenericCode(s string) string {
	s = colorizeStrings(s)
	return "\x1b[38;5;250m" + s + "\x1b[0m"
}

func colorizeStrings(s string) string {
	reDouble := regexp.MustCompile(`"([^"\\]|\\.)*"`)
	s = reDouble.ReplaceAllString(s, "\x1b[38;5;216m$0\x1b[0m")
	reSingle := regexp.MustCompile(`'([^'\\]|\\.)*'`)
	s = reSingle.ReplaceAllString(s, "\x1b[38;5;216m$0\x1b[0m")
	return s
}

func colorizeComments(ext, s string) string {
	if ext == ".py" || ext == ".rb" || ext == ".sh" {
		re := regexp.MustCompile(`(?m)#.*$`)
		return re.ReplaceAllString(s, "\x1b[38;5;244m$0\x1b[0m")
	}
	re := regexp.MustCompile(`(?m)//.*$`)
	return re.ReplaceAllString(s, "\x1b[38;5;244m$0\x1b[0m")
}

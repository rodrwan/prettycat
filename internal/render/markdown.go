package render

import (
	"fmt"
	"strings"
)

func renderMarkdown(in []byte, opts Options) (string, error) {
	s := string(in)
	if !opts.Color {
		return renderPlain([]byte(s)), nil
	}

	lines := strings.Split(s, "\n")
	var out strings.Builder
	inCode := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCode = !inCode
			if inCode {
				out.WriteString("\x1b[38;5;244m┌ code\x1b[0m\n")
			} else {
				out.WriteString("\x1b[38;5;244m└\x1b[0m\n")
			}
			continue
		}
		if inCode {
			out.WriteString("\x1b[38;5;179m" + line + "\x1b[0m\n")
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "### "):
			out.WriteString(fmt.Sprintf("\x1b[1;38;5;111m%s\x1b[0m\n", strings.TrimPrefix(trimmed, "### ")))
		case strings.HasPrefix(trimmed, "## "):
			out.WriteString(fmt.Sprintf("\x1b[1;38;5;117m%s\x1b[0m\n", strings.TrimPrefix(trimmed, "## ")))
		case strings.HasPrefix(trimmed, "# "):
			out.WriteString(fmt.Sprintf("\x1b[1;38;5;159m%s\x1b[0m\n", strings.TrimPrefix(trimmed, "# ")))
		case strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* "):
			out.WriteString("\x1b[38;5;212m•\x1b[0m " + strings.TrimSpace(trimmed[2:]) + "\n")
		case strings.HasPrefix(trimmed, "> "):
			out.WriteString("\x1b[38;5;244m│ \x1b[0m" + strings.TrimSpace(trimmed[2:]) + "\n")
		default:
			out.WriteString(styleInlineMarkdown(line) + "\n")
		}
	}

	return strings.TrimRight(out.String(), "\n") + "\n", nil
}

func styleInlineMarkdown(line string) string {
	line = stylePair(line, "**", "\x1b[1m", "\x1b[0m")
	line = stylePair(line, "`", "\x1b[38;5;179m", "\x1b[0m")
	return line
}

func stylePair(line, delim, start, end string) string {
	parts := strings.Split(line, delim)
	if len(parts) < 3 {
		return line
	}
	var out strings.Builder
	for i, p := range parts {
		if i%2 == 1 {
			out.WriteString(start + p + end)
		} else {
			out.WriteString(p)
		}
	}
	return out.String()
}

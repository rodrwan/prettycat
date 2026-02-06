package render

import "strings"

func renderPlain(in []byte) string {
	return strings.TrimRight(string(in), "\n") + "\n"
}

package style

import (
	"fmt"
	"strings"
)

const (
	reset  = "\x1b[0m"
	bold   = "\x1b[1m"
	gray   = "\x1b[38;5;244m"
	pink   = "\x1b[38;5;212m"
	border = "\x1b[38;5;240m"
)

func Header(title string, color bool) string {
	if !color {
		return fmt.Sprintf("==> %s <==\n", title)
	}
	return bold + pink + "==> " + title + " <==" + reset + "\n"
}

func Separator(color bool) string {
	if !color {
		return strings.Repeat("-", 40) + "\n"
	}
	return gray + border + strings.Repeat("─", 40) + reset + "\n"
}

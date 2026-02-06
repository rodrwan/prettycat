package pager

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/rodrwan/prettycat/internal/render"
)

func Run(docs []render.Doc, color bool, stdout io.Writer) error {
	content := joinDocs(docs)
	lines := strings.Split(content, "\n")
	height := terminalHeight()

	offset := 0
	query := ""
	matches := []int{}
	matchIdx := 0
	statusExtra := ""
	searching := false
	searchInput := ""

	restore, err := makeRaw(int(os.Stdin.Fd()))
	if err == nil {
		defer restore()
	}

	r := bufio.NewReader(os.Stdin)
	for {
		pageSize := computePageSize(height, searching)
		renderPage(stdout, lines, offset, pageSize, color, statusExtra, searching, searchInput)
		statusExtra = ""

		key, err := readKey(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if searching {
			switch key {
			case "enter":
				searching = false
				query = strings.TrimSpace(searchInput)
				matches = findMatches(lines, query)
				matchIdx = 0
				if len(matches) > 0 {
					offset = matches[0]
					statusExtra = fmt.Sprintf("match 1/%d for %q", len(matches), query)
				} else if query != "" {
					statusExtra = fmt.Sprintf("no matches for %q", query)
				}
			case "esc":
				searching = false
				searchInput = ""
				statusExtra = "search canceled"
			case "backspace":
				if len(searchInput) > 0 {
					searchInput = searchInput[:len(searchInput)-1]
				}
			default:
				if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
					searchInput += key
				}
			}
			continue
		}

		switch key {
		case "q":
			return nil
		case "j", "down":
			if offset+1 < len(lines) {
				offset++
			}
		case "k", "up":
			if offset > 0 {
				offset--
			}
		case "f", "pgdn", "space":
			offset += pageSize
			if offset >= len(lines) {
				offset = max(0, len(lines)-1)
			}
		case "b", "pgup":
			offset -= pageSize
			if offset < 0 {
				offset = 0
			}
		case "g":
			offset = 0
		case "G":
			offset = max(0, len(lines)-pageSize)
		case "/":
			searching = true
			searchInput = ""
			statusExtra = "type search and press Enter"
		case "n":
			if len(matches) > 0 {
				matchIdx = (matchIdx + 1) % len(matches)
				offset = matches[matchIdx]
				statusExtra = fmt.Sprintf("match %d/%d for %q", matchIdx+1, len(matches), query)
			}
		case "N":
			if len(matches) > 0 {
				matchIdx--
				if matchIdx < 0 {
					matchIdx = len(matches) - 1
				}
				offset = matches[matchIdx]
				statusExtra = fmt.Sprintf("match %d/%d for %q", matchIdx+1, len(matches), query)
			}
		}

		if offset > max(0, len(lines)-pageSize) {
			offset = max(0, len(lines)-pageSize)
		}
	}
}

func joinDocs(docs []render.Doc) string {
	var b strings.Builder
	for _, doc := range docs {
		b.WriteString(doc.Body)
	}
	return b.String()
}

func renderPage(out io.Writer, lines []string, offset, pageSize int, color bool, extra string, searching bool, searchInput string) {
	// Always repaint the screen in pager mode to avoid drifting output.
	fmt.Fprint(out, "\x1b[2J\x1b[H")

	end := min(len(lines), offset+pageSize)
	for i := offset; i < end; i++ {
		fmt.Fprintln(out, lines[i])
	}

	status := fmt.Sprintf("[%d-%d/%d] q quit | j/k or arrows | f/b/space page | / find | n/N next/prev", offset+1, end, len(lines))
	if extra != "" {
		status += " | " + extra
	}
	if color {
		status = "\x1b[38;5;244m" + status + "\x1b[0m"
	}
	fmt.Fprintln(out, status)

	if searching {
		prompt := "/" + searchInput
		if color {
			prompt = "\x1b[38;5;212m" + prompt + "\x1b[0m"
		}
		fmt.Fprintln(out, prompt)
	}
}

func readKey(r *bufio.Reader) (string, error) {
	c, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	switch c {
	case '\r', '\n':
		return "enter", nil
	case 127, 8:
		return "backspace", nil
	case 32:
		return "space", nil
	case 27:
		if r.Buffered() == 0 {
			return "esc", nil
		}
		c2, err := r.ReadByte()
		if err != nil || c2 != '[' {
			return "esc", nil
		}
		c3, err := r.ReadByte()
		if err != nil {
			return "esc", nil
		}
		switch c3 {
		case 'A':
			return "up", nil
		case 'B':
			return "down", nil
		case '5':
			if r.Buffered() > 0 {
				_, _ = r.ReadByte()
			}
			return "pgup", nil
		case '6':
			if r.Buffered() > 0 {
				_, _ = r.ReadByte()
			}
			return "pgdn", nil
		default:
			return "esc", nil
		}
	default:
		return string(c), nil
	}
}

func terminalHeight() int {
	if h, err := ttyHeight(int(os.Stdout.Fd())); err == nil && h > 0 {
		return h
	}
	height := 24
	if h := strings.TrimSpace(os.Getenv("LINES")); h != "" {
		if n, err := strconv.Atoi(h); err == nil && n > 5 {
			height = n
		}
	}
	return height
}

func ttyHeight(fd int) (int, error) {
	ws := &winsize{}
	_, _, errno := syscall.Syscall(
		uintptr(syscall.SYS_IOCTL),
		uintptr(fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return 0, errno
	}
	if ws.Row == 0 {
		return 0, fmt.Errorf("terminal row size is 0")
	}
	return int(ws.Row), nil
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func computePageSize(height int, searching bool) int {
	footerLines := 1
	if searching {
		footerLines++
	}
	pageSize := height - footerLines
	if pageSize < 3 {
		return 10
	}
	return pageSize
}

func makeRaw(fd int) (func(), error) {
	orig, err := getTermios(fd)
	if err != nil {
		return func() {}, err
	}
	raw := *orig
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	raw.Cflag |= syscall.CS8
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if err := setTermios(fd, &raw); err != nil {
		return func() {}, err
	}
	return func() {
		_ = setTermios(fd, orig)
	}, nil
}

func getTermios(fd int) (*syscall.Termios, error) {
	termios := &syscall.Termios{}
	_, _, errno := syscall.Syscall6(
		uintptr(syscall.SYS_IOCTL),
		uintptr(fd),
		uintptr(syscall.TIOCGETA),
		uintptr(unsafe.Pointer(termios)),
		0,
		0,
		0,
	)
	if errno != 0 {
		return nil, errno
	}
	return termios, nil
}

func setTermios(fd int, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall6(
		uintptr(syscall.SYS_IOCTL),
		uintptr(fd),
		uintptr(syscall.TIOCSETA),
		uintptr(unsafe.Pointer(termios)),
		0,
		0,
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}

func findMatches(lines []string, q string) []int {
	if q == "" {
		return nil
	}
	q = strings.ToLower(q)
	matches := make([]int, 0)
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), q) {
			matches = append(matches, i)
		}
	}
	return matches
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

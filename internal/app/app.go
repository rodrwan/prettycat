package app

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rodrwan/prettycat/internal/exitcode"
	"github.com/rodrwan/prettycat/internal/input"
	"github.com/rodrwan/prettycat/internal/pager"
	"github.com/rodrwan/prettycat/internal/render"
	"github.com/rodrwan/prettycat/internal/style"
)

type Config struct {
	Args      []string
	Version   string
	NoColor   bool
	Stdin     *os.File
	Stdout    *os.File
	Stderr    io.Writer
	IsTTYIn   func(*os.File) bool
	IsTTYOut  func(*os.File) bool
	OpenFile  input.FileOpener
	ReadAll   input.ReadAllFn
	PagerOpen func([]render.Doc, bool, io.Writer) error
}

func Run(cfg Config) int {
	if len(cfg.Args) == 0 && cfg.IsTTYIn(cfg.Stdin) {
		fmt.Fprintln(cfg.Stderr, "prettycat: no input (pass files or pipe stdin)")
		return exitcode.Usage
	}

	loaded, stdinHadData, err := input.Load(cfg.Args, cfg.Stdin, cfg.IsTTYIn, cfg.OpenFile, cfg.ReadAll)
	if err != nil {
		fmt.Fprintf(cfg.Stderr, "prettycat: %v\n", err)
		return exitcode.Error
	}

	if stdinHadData && len(cfg.Args) > 0 {
		fmt.Fprintln(cfg.Stderr, "prettycat: stdin data ignored because file arguments were provided")
	}

	color := useColor(cfg.NoColor, cfg.Stdout)
	docs := make([]render.Doc, 0, len(loaded.Sources))
	hadErr := len(loaded.Errors) > 0

	for _, e := range loaded.Errors {
		fmt.Fprintf(cfg.Stderr, "prettycat: %v\n", e)
	}

	for i, src := range loaded.Sources {
		doc, err := render.Render(src, render.Options{Color: color, Width: 100})
		if err != nil {
			hadErr = true
			fmt.Fprintf(cfg.Stderr, "prettycat: %v\n", err)
			continue
		}
		if len(loaded.Sources) > 1 {
			doc.Body = style.Header(src.Name, color) + doc.Body
			if i < len(loaded.Sources)-1 {
				doc.Body += style.Separator(color)
			}
		}
		docs = append(docs, doc)
	}

	if len(docs) == 0 {
		return exitcode.Error
	}

	if cfg.IsTTYOut(cfg.Stdout) {
		if err := cfg.PagerOpen(docs, color, cfg.Stdout); err != nil {
			fmt.Fprintf(cfg.Stderr, "prettycat: pager error: %v\n", err)
			hadErr = true
		}
	} else {
		for _, doc := range docs {
			if _, err := io.WriteString(cfg.Stdout, normalize(doc.Body)); err != nil {
				fmt.Fprintf(cfg.Stderr, "prettycat: write output: %v\n", err)
				return exitcode.Error
			}
		}
	}

	if hadErr {
		return exitcode.Error
	}
	return exitcode.OK
}

func useColor(noColor bool, stdout *os.File) bool {
	if noColor {
		return false
	}
	if strings.TrimSpace(os.Getenv("NO_COLOR")) != "" {
		return false
	}
	return IsTTYFile(stdout)
}

func normalize(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s
	}
	return s + "\n"
}

func IsTTYFile(f *os.File) bool {
	if f == nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func RunPager(docs []render.Doc, color bool, stdout io.Writer) error {
	return pager.Run(docs, color, stdout)
}

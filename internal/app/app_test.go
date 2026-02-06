package app

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rodrwan/prettycat/internal/exitcode"
	"github.com/rodrwan/prettycat/internal/render"
)

func TestRunContinuesOnMissingFileAndReturnsError(t *testing.T) {
	tmp := t.TempDir()
	okFile := filepath.Join(tmp, "ok.txt")
	if err := os.WriteFile(okFile, []byte("hello world\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	outFile := filepath.Join(tmp, "out.txt")
	out, err := os.Create(outFile)
	if err != nil {
		t.Fatalf("create out: %v", err)
	}
	defer out.Close()

	var stderr bytes.Buffer
	cfg := Config{
		Args:     []string{okFile, filepath.Join(tmp, "missing.txt")},
		NoColor:  true,
		Stdin:    os.Stdin,
		Stdout:   out,
		Stderr:   &stderr,
		IsTTYIn:  func(*os.File) bool { return true },
		IsTTYOut: func(*os.File) bool { return false },
		OpenFile: os.Open,
		ReadAll:  io.ReadAll,
		PagerOpen: func([]render.Doc, bool, io.Writer) error {
			t.Fatalf("pager should not be called")
			return nil
		},
	}

	code := Run(cfg)
	if code != exitcode.Error {
		t.Fatalf("Run code = %d, want %d", code, exitcode.Error)
	}

	if _, err := out.Seek(0, 0); err != nil {
		t.Fatalf("seek out: %v", err)
	}
	data, err := io.ReadAll(out)
	if err != nil {
		t.Fatalf("read out: %v", err)
	}
	if !strings.Contains(string(data), "hello world") {
		t.Fatalf("expected output to include valid file content, got %q", string(data))
	}
	if !strings.Contains(stderr.String(), "missing.txt") {
		t.Fatalf("expected stderr to include missing file error, got %q", stderr.String())
	}
}

func TestRunNoInputUsage(t *testing.T) {
	var stderr bytes.Buffer
	cfg := Config{
		Args:      nil,
		NoColor:   true,
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    &stderr,
		IsTTYIn:   func(*os.File) bool { return true },
		IsTTYOut:  func(*os.File) bool { return false },
		OpenFile:  os.Open,
		ReadAll:   io.ReadAll,
		PagerOpen: func([]render.Doc, bool, io.Writer) error { return nil },
	}
	if got := Run(cfg); got != exitcode.Usage {
		t.Fatalf("Run() = %d, want %d", got, exitcode.Usage)
	}
}

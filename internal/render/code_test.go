package render

import (
	"regexp"
	"strings"
	"testing"
)

func TestRenderCodeNoColorHasNoANSI(t *testing.T) {
	out, err := renderCode("main.go", []byte("package main\n\nfunc main() {}\n"), Options{Color: false})
	if err != nil {
		t.Fatalf("renderCode returned error: %v", err)
	}

	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	if ansi.MatchString(out) {
		t.Fatalf("expected no ANSI sequences, got %q", out)
	}
}

func TestRenderCodeColorUsesRealANSIBytes(t *testing.T) {
	out, err := renderCode("main.go", []byte("package main\nfunc main() {}\n"), Options{Color: true})
	if err != nil {
		t.Fatalf("renderCode returned error: %v", err)
	}
	if strings.Contains(out, "\\x1b[") {
		t.Fatalf("expected real ANSI escape bytes, got literal escape text: %q", out)
	}
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected ANSI escape bytes in color output, got %q", out)
	}
}

package render

import "testing"

func TestDetectKind(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Kind
	}{
		{name: "markdown", in: "README.md", want: KindMarkdown},
		{name: "markdown uppercase", in: "notes.MARKDOWN", want: KindMarkdown},
		{name: "code", in: "main.go", want: KindCode},
		{name: "plain no ext", in: "LICENSE", want: KindPlain},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := DetectKind(tc.in); got != tc.want {
				t.Fatalf("DetectKind(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

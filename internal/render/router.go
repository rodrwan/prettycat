package render

import (
	"fmt"

	"github.com/rodrwan/prettycat/internal/input"
)

func Render(src input.Source, opts Options) (Doc, error) {
	kind := DetectKind(src.Name)

	var (
		body string
		err  error
	)

	switch kind {
	case KindMarkdown:
		body, err = renderMarkdown(src.Data, opts)
	case KindCode:
		body, err = renderCode(src.Name, src.Data, opts)
	default:
		body = renderPlain(src.Data)
	}
	if err != nil {
		return Doc{}, fmt.Errorf("render %s: %w", src.Name, err)
	}

	return Doc{Title: src.Name, Body: body, Kind: kind}, nil
}

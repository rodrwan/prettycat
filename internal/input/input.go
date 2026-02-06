package input

import (
	"fmt"
	"io"
	"os"
)

type Source struct {
	Name    string
	Data    []byte
	IsStdin bool
}

type LoadResult struct {
	Sources []Source
	Errors  []error
}

type FileOpener func(name string) (*os.File, error)
type ReadAllFn func(r io.Reader) ([]byte, error)
type IsTTYFn func(*os.File) bool

func Load(args []string, stdin *os.File, isTTY IsTTYFn, openFile FileOpener, readAll ReadAllFn) (LoadResult, bool, error) {
	result := LoadResult{}
	stdinHasData := !isTTY(stdin)

	if len(args) == 0 {
		if !stdinHasData {
			return result, false, fmt.Errorf("no input: provide a file or pipe data through stdin")
		}
		data, err := readAll(stdin)
		if err != nil {
			return result, false, fmt.Errorf("read stdin: %w", err)
		}
		result.Sources = append(result.Sources, Source{Name: "stdin", Data: data, IsStdin: true})
		return result, false, nil
	}

	for _, path := range args {
		f, err := openFile(path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", path, err))
			continue
		}
		data, err := readAll(f)
		f.Close()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", path, err))
			continue
		}
		result.Sources = append(result.Sources, Source{Name: path, Data: data, IsStdin: false})
	}

	if len(result.Sources) == 0 {
		return result, stdinHasData, fmt.Errorf("no readable files")
	}

	return result, stdinHasData, nil
}

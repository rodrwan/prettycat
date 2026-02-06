package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rodrwan/prettycat/internal/app"
	"github.com/rodrwan/prettycat/internal/exitcode"
)

const version = "0.1.0"

func main() {
	var (
		showVersion bool
		noColor     bool
	)

	flag.BoolVar(&showVersion, "version", false, "print version")
	flag.BoolVar(&noColor, "no-color", false, "disable ANSI colors")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] [file ...]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Render beautiful terminal output for text, markdown and code files.")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		os.Exit(exitcode.Usage)
	}

	if showVersion {
		fmt.Println(version)
		os.Exit(exitcode.OK)
	}

	code := app.Run(app.Config{
		Args:      flag.Args(),
		Version:   version,
		NoColor:   noColor,
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		IsTTYIn:   app.IsTTYFile,
		IsTTYOut:  app.IsTTYFile,
		OpenFile:  os.Open,
		ReadAll:   app.ReadAll,
		PagerOpen: app.RunPager,
	})
	os.Exit(code)
}

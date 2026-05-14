package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

func OK(msg string) {
	fmt.Printf("  %s %s\n", green("✓"), msg)
}

func Fail(msg string) {
	fmt.Printf("  %s %s\n", red("✗"), msg)
}

func Warn(msg string) {
	fmt.Printf("  %s %s\n", yellow("!"), msg)
}

func Header(msg string) {
	fmt.Printf("\n%s\n", bold(msg))
}

func Fatal(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s %s: %v\n", red("error"), msg, err)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, red("error")+" "+format+"\n", args...)
	os.Exit(1)
}

package logger

import (
	"os"

	"github.com/mattn/go-isatty"
)

var (
	nocolor = "\033[0m"
	red     = "\033[0;31m"
	green   = "\033[0;32m"
	orange  = "\033[0;33m"
	blue    = "\033[0;34m"
	purple  = "\033[0;35m"
	cyan    = "\033[0;36m"
	gray    = "\033[0;37m"
)

func colorize(text string, color string) string {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return color + text + nocolor
	}
	return text
}

// Gray ANSI color applied to a string
func Gray(text string) string {
	return colorize(text, gray)
}

// Green ANSI color applied to a string
func Green(text string) string {
	return colorize(text, green)
}

// Orange ANSI color applied to a string
func Orange(text string) string {
	return colorize(text, orange)
}

// Red ANSI color applied to a string
func Red(text string) string {
	return colorize(text, red)
}

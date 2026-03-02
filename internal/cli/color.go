package cli

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

// colorEnabled tracks whether color output is on.
var colorEnabled = true

// DisableColor turns off color output.
func DisableColor() { colorEnabled = false }

// EnableColor turns on color output (default).
func EnableColor() { colorEnabled = true }

func init() {
	// Respect NO_COLOR env var (https://no-color.org)
	if os.Getenv("NO_COLOR") != "" {
		colorEnabled = false
	}
}

// Colored formatting helpers

func green(s string) string {
	if !colorEnabled {
		return s
	}
	return colorGreen + s + colorReset
}

func red(s string) string {
	if !colorEnabled {
		return s
	}
	return colorRed + s + colorReset
}

func yellow(s string) string {
	if !colorEnabled {
		return s
	}
	return colorYellow + s + colorReset
}

func cyan(s string) string {
	if !colorEnabled {
		return s
	}
	return colorCyan + s + colorReset
}

func bold(s string) string {
	if !colorEnabled {
		return s
	}
	return colorBold + s + colorReset
}

func dim(s string) string {
	if !colorEnabled {
		return s
	}
	return colorDim + s + colorReset
}

func successf(format string, a ...interface{}) string {
	return green(fmt.Sprintf(format, a...))
}

func errorf(format string, a ...interface{}) string {
	return red(fmt.Sprintf(format, a...))
}

func warnf(format string, a ...interface{}) string {
	return yellow(fmt.Sprintf(format, a...))
}

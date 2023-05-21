package tracerr

import (
	"fmt"
)

// Colorize outputs using [ANSI Escape Codes](https://en.wikipedia.org/wiki/ANSI_escape_code)

func color(code int, in string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, in)
}

func bold(in string) string {
	return color(1, in)
}

func black(in string) string {
	return color(30, in)
}

func red(in string) string {
	return color(31, in)
}

func yellow(in string) string {
	return color(33, in)
}

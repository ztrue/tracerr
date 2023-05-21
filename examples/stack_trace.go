package main

import (
	"fmt"
	"os"

	"github.com/ztrue/tracerr"
)

func main() {
	if err := read(); err != nil {
		// Dump raw stack trace.
		frames := tracerr.StackTrace(err)
		fmt.Printf("%#v\n", frames)
	}
}

func read() error {
	return readNonExistent()
}

func readNonExistent() error {
	_, err := os.ReadFile("/tmp/non_existent_file")
	return tracerr.Wrap(err)
}

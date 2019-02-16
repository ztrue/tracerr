package main

import (
	"io/ioutil"

	"github.com/ztrue/tracerr"
)

func main() {
	if err := read(); err != nil {
		// Save output to variable.
		text := tracerr.SprintSource(err)
		ioutil.WriteFile("/tmp/tracerr.log", []byte(text), 0644)
	}
}

func read() error {
	return readNonExistent()
}

func readNonExistent() error {
	_, err := ioutil.ReadFile("/tmp/non_existent_file")
	return tracerr.Wrap(err)
}

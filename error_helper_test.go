// This file is added for purpose of having an example of different path in tests.
package tracerr_test

import (
	"github.com/ztrue/tracerr"
)

func addFrameA(message string) *tracerr.Error {
	return addFrameB(message)
}

func addFrameB(message string) *tracerr.Error {
	return addFrameC(message)
}

func addFrameC(message string) *tracerr.Error {
	return tracerr.New(message)
}

package tracerr_test

import (
	"fmt"
	"testing"

	"github.com/ztrue/tracerr"
)

func BenchmarkNew(b *testing.B) {
	for _, frames := range []int{5, 10, 20, 40} {
		suffix := fmt.Sprintf("%d", frames)
		b.Run(suffix, func(b *testing.B) {
			err := tracerr.New("")
			// Reduce by number of parent frames in order to have a correct depth.
			depth := frames - len(err.StackTrace())
			if depth < 1 {
				panic("number of frames is negative")
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				addFrames(depth, "test error")
			}
		})
	}
}

func addFrames(depth int, message string) error {
	if depth <= 1 {
		return tracerr.New(message)
	}
	return addFrames(depth-1, message)
}

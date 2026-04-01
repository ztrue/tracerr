// Package tracerr makes error output more informative.
// It adds stack trace to error and can display error with source fragments.
//
// Check example of output here https://github.com/ztrue/tracerr
package tracerr

import (
	"fmt"
	"runtime"
)

// DefaultCap is a default cap for frames array.
// It can be changed to number of expected frames
// for purpose of performance optimisation.
var DefaultCap = 20

// Error is an error with stack trace.
type Error interface {
	Error() string
	StackTrace() []Frame
	Unwrap() error
}

type errorData struct {
	// err contains original error.
	err error
	// pcs contains raw program counters, resolved lazily to frames.
	pcs []uintptr
	// frames contains pre-resolved stack trace.
	frames []Frame
}

// CustomError creates an error with provided frames.
func CustomError(err error, frames []Frame) Error {
	return &errorData{
		err:    err,
		frames: frames,
	}
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
func Errorf(message string, args ...interface{}) Error {
	return trace(fmt.Errorf(message, args...), 2)
}

// New creates new error with stacktrace.
func New(message string) Error {
	return trace(fmt.Errorf(message), 2)
}

// Wrap adds stacktrace to existing error.
func Wrap(err error) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	return trace(err, 2)
}

// Unwrap returns the original error.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if !ok {
		return err
	}
	return e.Unwrap()
}

// Error returns error message.
func (e *errorData) Error() string {
	return e.err.Error()
}

// StackTrace resolves and returns the stack trace, caching the result.
func (e *errorData) StackTrace() []Frame {
	if e.pcs == nil {
		return e.frames
	}
	cf := runtime.CallersFrames(e.pcs)
	frames := make([]Frame, 0, len(e.pcs))
	for {
		f, more := cf.Next()
		frames = append(frames, Frame{
			Func: f.Function,
			Line: f.Line,
			Path: f.File,
		})
		if !more {
			break
		}
	}
	e.frames = frames
	e.pcs = nil
	return e.frames
}

// Unwrap returns the original error.
func (e *errorData) Unwrap() error {
	return e.err
}

// Frame is a single step in stack trace.
type Frame struct {
	// Func contains a function name.
	Func string
	// Line contains a line number.
	Line int
	// Path contains a file path.
	Path string
}

// StackTrace returns stack trace of an error.
// It will be empty if err is not of type Error.
func StackTrace(err error) []Frame {
	e, ok := err.(Error)
	if !ok {
		return nil
	}
	return e.StackTrace()
}

// String formats Frame to string.
func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}

func trace(err error, skip int) Error {
	pcs := make([]uintptr, DefaultCap)
	for {
		n := runtime.Callers(skip+1, pcs)
		if n < len(pcs) {
			pcs = pcs[:n]
			break
		}
		pcs = make([]uintptr, len(pcs)*2)
	}
	return &errorData{
		err: err,
		pcs: pcs,
	}
}

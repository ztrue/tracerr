// Package tracerr makes error output more informative.
// It adds stack trace to error and can display error with source fragments.
//
// Check example of output here https://github.com/ztrue/tracerr
package tracerr

import (
	"fmt"
	"runtime"
)

// Error is an error with stack trace.
type Error struct {
	frames []Frame
	err    error
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
func Errorf(message string, args ...interface{}) *Error {
	return trace(fmt.Errorf(message, args...), 2)
}

// New creates new error with stacktrace.
func New(message string) *Error {
	return trace(fmt.Errorf(message), 2)
}

// Wrap adds stacktrace to existing error.
func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	e, ok := err.(*Error)
	if ok {
		return e
	}
	return trace(err, 2)
}

// Error returns error message.
func (e *Error) Error() string {
	return e.err.Error()
}

// StackTrace returns stack trace of an error.
func (e *Error) StackTrace() []Frame {
	if e == nil {
		return nil
	}
	return e.frames
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
// It will be empty if err is not of type *Error.
func StackTrace(err error) []Frame {
	e, ok := err.(*Error)
	if !ok {
		return nil
	}
	return e.StackTrace()
}

// String formats Frame to string.
func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}

func trace(err error, skip int) *Error {
	var frames []Frame
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		frame := Frame{
			Func: fn.Name(),
			Line: line,
			Path: path,
		}
		frames = append(frames, frame)
		skip++
	}
	return &Error{
		frames: frames,
		err:    err,
	}
}

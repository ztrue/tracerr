package tracerr_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ztrue/tracerr"
)

type ErrorTestCase struct {
	Error              tracerr.Error
	ExpectedMessage    string
	ExpectedStackTrace []tracerr.Frame
}

func TestError(t *testing.T) {
	cases := []ErrorTestCase{
		{
			Error:              nil,
			ExpectedMessage:    "",
			ExpectedStackTrace: nil,
		},
		{
			Error:              tracerr.Wrap(nil),
			ExpectedMessage:    "",
			ExpectedStackTrace: nil,
		},
		{
			Error:           tracerr.New("error message text"),
			ExpectedMessage: "error message text",
			ExpectedStackTrace: []tracerr.Frame{
				{
					Func: "github.com/ztrue/tracerr_test.TestError",
					Line: 31,
					Path: "/src/github.com/ztrue/tracerr/error_test.go",
				},
			},
		},
		{
			Error:           tracerr.Errorf("invalid argument %d: %#v", 5, "foo"),
			ExpectedMessage: "invalid argument 5: \"foo\"",
			ExpectedStackTrace: []tracerr.Frame{
				{
					Func: "github.com/ztrue/tracerr_test.TestError",
					Line: 42,
					Path: "/src/github.com/ztrue/tracerr/error_test.go",
				},
			},
		},
		{
			Error:           tracerr.Wrap(errors.New("wrapped error")),
			ExpectedMessage: "wrapped error",
			ExpectedStackTrace: []tracerr.Frame{
				{
					Func: "github.com/ztrue/tracerr_test.TestError",
					Line: 53,
					Path: "/src/github.com/ztrue/tracerr/error_test.go",
				},
			},
		},
		{
			Error:           addFrameA("error with stack trace").(tracerr.Error),
			ExpectedMessage: "error with stack trace",
			ExpectedStackTrace: []tracerr.Frame{
				{
					Func: "github.com/ztrue/tracerr_test.addFrameC",
					Line: 17,
					Path: "/src/github.com/ztrue/tracerr/error_helper_test.go",
				},
				{
					Func: "github.com/ztrue/tracerr_test.addFrameB",
					Line: 13,
					Path: "/src/github.com/ztrue/tracerr/error_helper_test.go",
				},
				{
					Func: "github.com/ztrue/tracerr_test.addFrameA",
					Line: 9,
					Path: "/src/github.com/ztrue/tracerr/error_helper_test.go",
				},
				{
					Func: "github.com/ztrue/tracerr_test.TestError",
					Line: 64,
					Path: "/src/github.com/ztrue/tracerr/error_test.go",
				},
			},
		},
		{
			Error:           tracerr.Wrap(addFrameA("error wrapped twice")),
			ExpectedMessage: "error wrapped twice",
			ExpectedStackTrace: []tracerr.Frame{
				{
					Func: "github.com/ztrue/tracerr_test.addFrameC",
					Line: 17,
					Path: "/src/github.com/ztrue/tracerr/error_helper_test.go",
				},
				{
					Func: "github.com/ztrue/tracerr_test.addFrameB",
					Line: 13,
					Path: "/src/github.com/ztrue/tracerr/error_helper_test.go",
				},
				{
					Func: "github.com/ztrue/tracerr_test.addFrameA",
					Line: 9,
					Path: "/src/github.com/ztrue/tracerr/error_helper_test.go",
				},
				{
					Func: "github.com/ztrue/tracerr_test.TestError",
					Line: 90,
					Path: "/src/github.com/ztrue/tracerr/error_test.go",
				},
			},
		},
	}

	for i, c := range cases {
		if c.Error == nil {
			if c.ExpectedMessage != "" {
				t.Errorf(
					"cases[%#v].Error = nil; want %#v",
					i, c.ExpectedMessage,
				)
			}
		} else if c.Error.Error() != c.ExpectedMessage {
			t.Errorf(
				"cases[%#v].Error.Error() = %#v; want %#v",
				i, c.Error.Error(), c.ExpectedMessage,
			)
		}

		if c.ExpectedStackTrace == nil {
			if c.Error != nil && c.Error.StackTrace() != nil {
				t.Errorf(
					"cases[%#v].Error.StackTrace() = %#v; want %#v",
					i, c.Error.StackTrace(), nil,
				)
			}
			if tracerr.StackTrace(c.Error) != nil {
				t.Errorf(
					"tracerr.StackTrace(cases[%#v].Error) = %#v; want %#v",
					i, tracerr.StackTrace(c.Error), nil,
				)
			}
			continue
		}

		frames1 := c.Error.StackTrace()
		frames2 := tracerr.StackTrace(c.Error)
		for k, frames := range [][]tracerr.Frame{frames1, frames2} {
			// Different failing message, depend on stack trace method.
			var pattern string
			if k == 0 {
				pattern = "cases[%#v].Error.StackTrace()"
			} else {
				pattern = "tracerr.StackTrace(cases[%#v].Error)"
			}
			prefix := fmt.Sprintf(pattern, i)
			// There must be at least two frames of test runner.
			expectedMinLen := len(c.ExpectedStackTrace) + 2
			if len(frames) < expectedMinLen {
				t.Errorf(
					"len(%s) = %#v; want >= %#v",
					prefix, len(frames), expectedMinLen,
				)
			}
			for j, expectedFrame := range c.ExpectedStackTrace {
				if frames[j].Func != expectedFrame.Func {
					t.Errorf(
						"%s[%#v].Func = %#v; want %#v",
						prefix, j, frames[j].Func, expectedFrame.Func,
					)
				}
				if frames[j].Line != expectedFrame.Line {
					t.Errorf(
						"%s[%#v].Line = %#v; want %#v",
						prefix, j, frames[j].Line, expectedFrame.Line,
					)
				}
				if !strings.HasSuffix(frames[j].Path, expectedFrame.Path) {
					t.Errorf(
						"%s[%#v].Path = %#v; want to has suffix %#v",
						prefix, j, frames[j].Path, expectedFrame.Path,
					)
				}
			}
		}

	}
}

func TestCustomError(t *testing.T) {
	err := errors.New("some error")
	frames := []tracerr.Frame{
		{
			Func: "main.foo",
			Line: 42,
			Path: "/src/github.com/john/doe/foobar.go",
		},
		{
			Func: "main.bar",
			Line: 43,
			Path: "/src/github.com/john/doe/bazqux.go",
		},
	}
	customErr := tracerr.CustomError(err, frames)
	message := customErr.Error()
	if message != err.Error() {
		t.Errorf(
			"customErr.Error() = %#v; want %#v",
			message, err.Error(),
		)
	}
	unwrapped := customErr.Unwrap()
	if unwrapped != err {
		t.Errorf(
			"customErr.Unwrap() = %#v; want %#v",
			unwrapped, err,
		)
	}
	stackTrace := customErr.StackTrace()
	if len(stackTrace) != len(frames) {
		t.Errorf(
			"len(customErr.StackTrace()) = %#v; want %#v",
			len(stackTrace), len(frames),
		)
	}
	for i, frame := range frames {
		if stackTrace[i] != frame {
			t.Errorf(
				"customErr.StackTrace()[%#v] = %#v; want %#v",
				i, stackTrace[i], frame,
			)
		}
	}
}

func TestErrorNil(t *testing.T) {
	wrapped := wrapError(nil)
	if wrapped != nil {
		t.Errorf(
			"wrapped = %#v; want nil",
			wrapped,
		)
	}
}

func TestFrameString(t *testing.T) {
	frame := tracerr.Frame{
		Func: "main.read",
		Line: 1337,
		Path: "/src/github.com/john/doe/foobar.go",
	}
	expected := "/src/github.com/john/doe/foobar.go:1337 main.read()"
	if frame.String() != expected {
		t.Errorf(
			"frame.String() = %#v; want %#v",
			frame.String(), expected,
		)
	}
}

func TestStackTraceNotInstance(t *testing.T) {
	err := errors.New("regular error")
	if tracerr.StackTrace(err) != nil {
		t.Errorf(
			"tracerr.StackTrace(%#v) = %#v; want %#v",
			err, tracerr.StackTrace(err), nil,
		)
	}
}

type UnwrapTestCase struct {
	Error error
	Wrap  bool
}

func TestUnwrap(t *testing.T) {
	cases := []UnwrapTestCase{
		{
			Error: nil,
		},
		{
			Error: fmt.Errorf("some error #%d", 9),
			Wrap:  false,
		},
		{
			Error: fmt.Errorf("some error #%d", 9),
			Wrap:  true,
		},
	}

	for i, c := range cases {
		err := c.Error
		if c.Wrap {
			err = tracerr.Wrap(err)
		}
		unwrappedError := tracerr.Unwrap(err)
		if unwrappedError != c.Error {
			t.Errorf(
				"tracerr.Unwrap(cases[%#v].Error) = %#v; want %#v",
				i, unwrappedError, c.Error,
			)
		}
	}
}

func wrapError(err error) error {
	return tracerr.Wrap(err)
}

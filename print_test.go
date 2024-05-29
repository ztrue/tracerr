package tracerr_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ztrue/tracerr"
)

// PrintTestCase
type PrintTestCase struct {
	Output               string
	Printer              func()
	ExpectedRows         []string
	ExpectedMinExtraRows int
}

func TestPrint(t *testing.T) {
	message := "runtime error: index out of range"
	err := addFrameA(message)

	cases := []PrintTestCase{
		{
			Output: tracerr.Sprint(nil),
			Printer: func() {
				tracerr.Print(nil)
			},
			ExpectedRows: []string{
				"",
			},
			ExpectedMinExtraRows: 0,
		},
		{
			Output: tracerr.Sprint(errors.New("regular error")),
			Printer: func() {
				tracerr.Print(errors.New("regular error"))
			},
			ExpectedRows: []string{
				"regular error",
			},
			ExpectedMinExtraRows: 0,
		},
		{
			Output: tracerr.Sprint(err),
			Printer: func() {
				tracerr.Print(err)
			},
			ExpectedRows: []string{
				message,
				"/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
			},
			ExpectedMinExtraRows: 2,
		},
		{
			Output: tracerr.SprintSource(err),
			Printer: func() {
				tracerr.PrintSource(err)
			},
			ExpectedRows: []string{
				message,
				"",
				"/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"14\t}",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"19\t",
				"",
				"/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"10\t}",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"15\t",
				"",
				"/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"6\t)",
				"7\t",
				"8\tfunc addFrameA(message string) error {",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"11\t",
				"",
				"/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
				"23\t",
				"24\tfunc TestPrint(t *testing.T) {",
				"25\t\tmessage := \"runtime error: index out of range\"",
				"26\t\terr := addFrameA(message)",
				"27\t",
				"28\t\tcases := []PrintTestCase{",
				"",
			},
			ExpectedMinExtraRows: 2,
		},
		{
			Output: tracerr.SprintSource(err, 2, 1),
			Printer: func() {
				tracerr.PrintSource(err, 2, 1)
			},
			ExpectedRows: []string{
				message,
				"",
				"/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"",
				"/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"",
				"/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"7\t",
				"8\tfunc addFrameA(message string) error {",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"",
				"/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
				"24\tfunc TestPrint(t *testing.T) {",
				"25\t\tmessage := \"runtime error: index out of range\"",
				"26\t\terr := addFrameA(message)",
				"27\t",
				"",
			},
			ExpectedMinExtraRows: 2,
		},
		{
			Output: tracerr.SprintSource(err, 4),
			Printer: func() {
				tracerr.PrintSource(err, 4)
			},
			ExpectedRows: []string{
				message,
				"",
				"/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"",
				"/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"",
				"/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"7\t",
				"8\tfunc addFrameA(message string) error {",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"",
				"/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
				"24\tfunc TestPrint(t *testing.T) {",
				"25\t\tmessage := \"runtime error: index out of range\"",
				"26\t\terr := addFrameA(message)",
				"27\t",
				"",
			},
			ExpectedMinExtraRows: 2,
		},
		{
			Output: tracerr.SprintSource(err, -1, -1),
			Printer: func() {
				tracerr.PrintSource(err, -1, -1)
			},
			ExpectedRows: []string{
				message,
				"",
				"/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"17\t\treturn tracerr.New(message)",
				"",
				"/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"13\t\treturn addFrameC(message)",
				"",
				"/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"9\t\treturn addFrameB(message)",
				"",
				"/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
				"26\t\terr := addFrameA(message)",
				"",
			},
			ExpectedMinExtraRows: 2,
		},
		{
			Output: tracerr.SprintSource(err, 0, 4),
			Printer: func() {
				tracerr.PrintSource(err, 0, 4)
			},
			ExpectedRows: []string{
				message,
				"",
				"/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"19\t",
				"",
				"/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"",
				"/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"",
				"/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
				"26\t\terr := addFrameA(message)",
				"27\t",
				"28\t\tcases := []PrintTestCase{",
				"29\t\t\t{",
				"30\t\t\t\tOutput: tracerr.Sprint(nil),",
			},
			ExpectedMinExtraRows: 2,
		},
		{
			Output: tracerr.SprintSourceColor(err, 1, 1),
			Printer: func() {
				tracerr.PrintSourceColor(err, 1, 1)
			},
			ExpectedRows: []string{
				message,
				"",
				bold("/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()"),
				black("16") + "\tfunc addFrameC(message string) error {",
				red("17\t\treturn tracerr.New(message)"),
				black("18") + "\t}",
				"",
				bold("/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()"),
				black("12") + "\tfunc addFrameB(message string) error {",
				red("13\t\treturn addFrameC(message)"),
				black("14") + "\t}",
				"",
				bold("/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()"),
				black("8") + "\tfunc addFrameA(message string) error {",
				red("9\t\treturn addFrameB(message)"),
				black("10") + "\t}",
				"",
				bold("/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()"),
				black("25") + "\t\tmessage := \"runtime error: index out of range\"",
				red("26\t\terr := addFrameA(message)"),
				black("27") + "\t",
				"",
			},
			ExpectedMinExtraRows: 2,
		},
	}

	for i, c := range cases {
		assertRows(t, i, c.Output, c.ExpectedRows, c.ExpectedMinExtraRows)
		output := captureOutput(c.Printer)
		assertRows(t, i, output, c.ExpectedRows, c.ExpectedMinExtraRows)
	}
}

func TestNoLine(t *testing.T) {
	err := tracerr.CustomError(
		errors.New("some error"),
		[]tracerr.Frame{
			{
				Func: "main.Foo",
				Line: 1337,
				Path: "error_helper_test.go",
			},
			{
				Func: "main.Bar",
				Line: 1338,
				Path: "error_helper_test.go",
			},
		},
	)
	output := tracerr.SprintSource(err)
	expectedRows := []string{
		"some error",
		"",
		"error_helper_test.go:1337 main.Foo()",
		"tracerr: too few lines, got 19, want 1337",
		"",
		"error_helper_test.go:1338 main.Bar()",
		"tracerr: too few lines, got 19, want 1338",
		"",
	}
	expected := strings.Join(expectedRows, "\n")
	if output != expected {
		t.Errorf(
			"tracerr.SprintSource(err) = %#v; want %#v",
			output, expected,
		)
	}
}

func TestNoLineColor(t *testing.T) {
	err := tracerr.CustomError(
		errors.New("some error"),
		[]tracerr.Frame{
			{
				Func: "main.Foo",
				Line: 1337,
				Path: "error_helper_test.go",
			},
			{
				Func: "main.Bar",
				Line: 1338,
				Path: "error_helper_test.go",
			},
		},
	)
	output := tracerr.SprintSourceColor(err)
	expectedRows := []string{
		"some error",
		"",
		bold("error_helper_test.go:1337 main.Foo()"),
		yellow("tracerr: too few lines, got 19, want 1337"),
		"",
		bold("error_helper_test.go:1338 main.Bar()"),
		yellow("tracerr: too few lines, got 19, want 1338"),
		"",
	}
	expected := strings.Join(expectedRows, "\n")
	if output != expected {
		t.Errorf(
			"tracerr.SprintSource(err) = %#v; want %#v",
			output, expected,
		)
	}
}

func TestNoSourceFile(t *testing.T) {
	err := tracerr.CustomError(
		errors.New("some error"),
		[]tracerr.Frame{
			{
				Func: "main.Foo",
				Line: 42,
				Path: "/tmp/not_exists.go",
			},
			{
				Func: "main.Bar",
				Line: 43,
				Path: "/tmp/not_exists_2.go",
			},
		},
	)
	output := tracerr.SprintSource(err)
	expectedRows := []string{
		"some error",
		"",
		"/tmp/not_exists.go:42 main.Foo()",
		"tracerr: file /tmp/not_exists.go not found",
		"",
		"/tmp/not_exists_2.go:43 main.Bar()",
		"tracerr: file /tmp/not_exists_2.go not found",
		"",
	}
	expected := strings.Join(expectedRows, "\n")
	if output != expected {
		t.Errorf(
			"tracerr.SprintSource(err) = %#v; want %#v",
			output, expected,
		)
	}
}

func TestNoSourceFileColor(t *testing.T) {
	err := tracerr.CustomError(
		errors.New("some error"),
		[]tracerr.Frame{
			{
				Func: "main.Foo",
				Line: 42,
				Path: "/tmp/not_exists.go",
			},
			{
				Func: "main.Bar",
				Line: 43,
				Path: "/tmp/not_exists_2.go",
			},
		},
	)
	output := tracerr.SprintSourceColor(err)
	expectedRows := []string{
		"some error",
		"",
		bold("/tmp/not_exists.go:42 main.Foo()"),
		yellow("tracerr: file /tmp/not_exists.go not found"),
		"",
		bold("/tmp/not_exists_2.go:43 main.Bar()"),
		yellow("tracerr: file /tmp/not_exists_2.go not found"),
		"",
	}
	expected := strings.Join(expectedRows, "\n")
	if output != expected {
		t.Errorf(
			"tracerr.SprintSource(err) = %#v; want %#v",
			output, expected,
		)
	}
}

func assertRows(t *testing.T, i int, output string, expectedRows []string, extra int) {
	rows := strings.Split(output, "\n")
	// There must be at least "extra" frames of test runner.
	expectedMinLen := len(expectedRows) + extra
	if len(rows) < expectedMinLen {
		t.Fatalf(
			"case #%d: len(rows) = %#v; want >= %#v",
			i, len(rows), expectedMinLen,
		)
	}
	// Remove root path, cause it could be different on different environments.
	re := regexp.MustCompile("([^/]*)/.*(/tracerr/.*)")
	for j, expectedRow := range expectedRows {
		row := re.ReplaceAllString(rows[j], "$1$2")
		if row != expectedRow {
			t.Errorf(
				"case #%d: rows[%#v] = %#v; want %#v",
				i, j, row, expectedRow,
			)
		}
	}
}

func captureOutput(fn func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err.Error())
	}
	stdout := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func bold(in string) string {
	return fmt.Sprintf("\x1b[1m%s\x1b[0m", in)
}

func black(in string) string {
	return fmt.Sprintf("\x1b[30m%s\x1b[0m", in)
}

func red(in string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", in)
}

func yellow(in string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", in)
}

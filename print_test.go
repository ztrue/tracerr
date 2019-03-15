package tracerr_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/logrusorgru/aurora"

	"github.com/ztrue/tracerr"
)

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
				"/src/github.com/ztrue/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"/src/github.com/ztrue/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
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
				"/src/github.com/ztrue/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"14\t}",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"19\t",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"10\t}",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"15\t",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"6\t)",
				"7\t",
				"8\tfunc addFrameA(message string) error {",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"11\t",
				"",
				"/src/github.com/ztrue/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
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
				"/src/github.com/ztrue/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"7\t",
				"8\tfunc addFrameA(message string) error {",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"",
				"/src/github.com/ztrue/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
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
				"/src/github.com/ztrue/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"7\t",
				"8\tfunc addFrameA(message string) error {",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"",
				"/src/github.com/ztrue/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
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
				"/src/github.com/ztrue/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"17\t\treturn tracerr.New(message)",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"13\t\treturn addFrameC(message)",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"9\t\treturn addFrameB(message)",
				"",
				"/src/github.com/ztrue/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
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
				"/src/github.com/ztrue/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()",
				"17\t\treturn tracerr.New(message)",
				"18\t}",
				"19\t",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()",
				"13\t\treturn addFrameC(message)",
				"14\t}",
				"15\t",
				"16\tfunc addFrameC(message string) error {",
				"17\t\treturn tracerr.New(message)",
				"",
				"/src/github.com/ztrue/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()",
				"9\t\treturn addFrameB(message)",
				"10\t}",
				"11\t",
				"12\tfunc addFrameB(message string) error {",
				"13\t\treturn addFrameC(message)",
				"",
				"/src/github.com/ztrue/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()",
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
				aurora.Bold("/src/github.com/ztrue/tracerr/error_helper_test.go:17 github.com/ztrue/tracerr_test.addFrameC()").String(),
				aurora.Black("16").String() + "\tfunc addFrameC(message string) error {",
				aurora.Red("17\t\treturn tracerr.New(message)").String(),
				aurora.Black("18").String() + "\t}",
				"",
				aurora.Bold("/src/github.com/ztrue/tracerr/error_helper_test.go:13 github.com/ztrue/tracerr_test.addFrameB()").String(),
				aurora.Black("12").String() + "\tfunc addFrameB(message string) error {",
				aurora.Red("13\t\treturn addFrameC(message)").String(),
				aurora.Black("14").String() + "\t}",
				"",
				aurora.Bold("/src/github.com/ztrue/tracerr/error_helper_test.go:9 github.com/ztrue/tracerr_test.addFrameA()").String(),
				aurora.Black("8").String() + "\tfunc addFrameA(message string) error {",
				aurora.Red("9\t\treturn addFrameB(message)").String(),
				aurora.Black("10").String() + "\t}",
				"",
				aurora.Bold("/src/github.com/ztrue/tracerr/print_test.go:26 github.com/ztrue/tracerr_test.TestPrint()").String(),
				aurora.Black("25").String() + "\t\tmessage := \"runtime error: index out of range\"",
				aurora.Red("26\t\terr := addFrameA(message)").String(),
				aurora.Black("27").String() + "\t",
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
		aurora.Bold("error_helper_test.go:1337 main.Foo()").String(),
		aurora.Brown("tracerr: too few lines, got 19, want 1337").String(),
		"",
		aurora.Bold("error_helper_test.go:1338 main.Bar()").String(),
		aurora.Brown("tracerr: too few lines, got 19, want 1338").String(),
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
		aurora.Bold("/tmp/not_exists.go:42 main.Foo()").String(),
		aurora.Brown("tracerr: file /tmp/not_exists.go not found").String(),
		"",
		aurora.Bold("/tmp/not_exists_2.go:43 main.Bar()").String(),
		aurora.Brown("tracerr: file /tmp/not_exists_2.go not found").String(),
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
	re := regexp.MustCompile("([^/]*)/.*(/src/github\\.com/ztrue/tracerr/.*)")
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

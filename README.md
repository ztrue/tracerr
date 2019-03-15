# Golang Errors with Stack Trace and Source Fragments

[![GoDoc](https://godoc.org/github.com/ztrue/tracerr?status.svg)](https://godoc.org/github.com/ztrue/tracerr)
[![Report](https://goreportcard.com/badge/github.com/ztrue/tracerr)](https://goreportcard.com/report/github.com/ztrue/tracerr)
[![Coverage Status](https://coveralls.io/repos/github/ztrue/tracerr/badge.svg?branch=master)](https://coveralls.io/github/ztrue/tracerr?branch=master)
[![Build Status](https://travis-ci.com/ztrue/tracerr.svg?branch=master)](https://travis-ci.com/ztrue/tracerr)

Tired of uninformative error output? Probably this will be more convenient:

![Output](output.png)

## Example

```go
package main

import (
	"io/ioutil"

	"github.com/ztrue/tracerr"
)

func main() {
	if err := read(); err != nil {
		tracerr.PrintSourceColor(err)
	}
}

func read() error {
	return readNonExistent()
}

func readNonExistent() error {
	_, err := ioutil.ReadFile("/tmp/non_existent_file")
	// Add stack trace to existing error, no matter if it's nil.
	return tracerr.Wrap(err)
}
```

Find more executable examples in [examples](examples) dir.

## How to Use

### Import

```go
import "github.com/ztrue/tracerr"
```

### Create New Error

```go
err := tracerr.New("some error")
```

Or:

```go
err := tracerr.Errorf("some error %d", num)
```

### Add Stack Trace to Existing Error

> If `err` is `nil` then it still be `nil` with no stack trace added.

```go
err = tracerr.Wrap(err)
```

### Print Error and Stack Trace

> Stack trace will be printed only if `err` is of type `tracerr.Error`, otherwise just error text will be shown.

This will print error message and stack trace if any:

```go
tracerr.Print(err)
```

This will add source code:

```go
tracerr.PrintSource(err)
```

It's able to set up number of lines of code to display for each frame, which is `6` by default:

```go
tracerr.PrintSource(err, 9)
```

Or to set up number of lines before and after traced line:

```go
tracerr.PrintSource(err, 5, 2)
```

The same, but with color, which is much more useful:

```go
tracerr.PrintSourceColor(err)
```

```go
tracerr.PrintSourceColor(err, 9)
```

```go
tracerr.PrintSourceColor(err, 5, 2)
```

### Save Output to Variable

It's also able to save output to variable instead of printing it, which works the same way:

```go
text := tracerr.Sprint(err)
```

```go
text := tracerr.SprintSource(err)
```

```go
text := tracerr.SprintSource(err, 9)
```

```go
text := tracerr.SprintSource(err, 5, 2)
```

### Get Stack Trace

> Stack trace will be empty if `err` is not an instance of `tracerr.Error`.

```go
frames := tracerr.StackTrace(err)
```

Or if `err` is of type `tracerr.Error`:

```go
frames := err.StackTrace()
```

### Get Original Error

> Unwrapped error will be `nil` if `err` is `nil` and will be the same error if `err` is not an instance of `tracerr.Error`.

```go
err = tracerr.Unwrap(err)
```

Or if `err` is of type `tracerr.Error`:

```go
err = err.Unwrap()
```

## Performance

Stack trace causes a performance overhead, depending on a stack trace depth. This can be insignificant in a number of situations (such as HTTP request handling), however, avoid of adding a stack trace for really hot spots where a high number of errors created frequently, this can be inefficient.

> Benchmarks done on a MacBook Pro 2015 with go 1.11.

Benchmarks for creating a new error with a stack trace of different depth:

```
BenchmarkNew/5    200000    5646 ns/op    976 B/op   4 allocs/op
BenchmarkNew/10   200000   11565 ns/op    976 B/op   4 allocs/op
BenchmarkNew/20    50000   25629 ns/op    976 B/op   4 allocs/op
BenchmarkNew/40    20000   65833 ns/op   2768 B/op   5 allocs/op
```

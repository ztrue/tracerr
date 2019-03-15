# Changelog

## [0.3.0] - 2019-03-15

### Added

- `tracerr.CustomError()` that allows to create error with custom stack trace.

### Changed

- `*tracerr.Error` struct replaced with `tracerr.Error` interface.

## [0.2.1] - 2019-02-16

### Added

- Benchmarks.
- `DefaultCap` variable for performance tuning purposes.

### Changed

- Stack trace performance optimisation.

## [0.2.0] - 2019-02-15

### Added

- `tracerr.Unwrap()` and `Error.Unwrap()` that returns the original error.

## [0.1.2] - 2019-02-12

### Changed

- RWMutex added for files caching, which fixing concurrent cache writing or writing-reading if any.

## [0.1.1] - 2019-02-09

### Added

- Changelog.
- `go.mod` file.
- License.
- Tests with 100% coverage.
- Travis CI.

### Changed

- `Error.Err` and `Error.Frames` properties are now exported.
- `Error.Error()` called on `nil` now returns empty string instead of panics.
- All print and sprint functions called with `nil` error now returns or prints empty string instead of panics.

## [0.1.0] - 2019-02-06

### Added

- Initial version.

# Changelog

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

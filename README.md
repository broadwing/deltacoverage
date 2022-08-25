# deltacoverage
Go package to provide delta coverage between your tests.

## Example

Go project [bench](https://github.com/thiagonache/bench).

```shell
$ deltacoverage TestNonOKStatusRecordedAsFailure
TestNonOKStatusRecordedAsFailure 1.1%
```

The test `TestNonOKStatusRecordedAsFailure` contributes to exclusive 1.1% of the total coverage.

## Motivation

```text
Herb Derby came up with this metric of “delta coverage”. 
You look at your test suite and you measure what coverage each test adds uniquely that no other test provides.
If you have lots of tests with no delta coverage, so that you could just delete them and not lose your 
ability to exercise parts of the system, then you should delete those tests, unless they have some some 
communication purpose.
—Kent Beck, “Is TDD Dead?”
```

Thank you @bitfield for the suggestion.

## Roadmap

- POC (WIP)
- Brute force implementation using `go test` commands
- Borrow `cover.go` code and create `deltagecoverage.go`
- (Dream) Implement as Go CLI feature.

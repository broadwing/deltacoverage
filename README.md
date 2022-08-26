# deltacoverage
Go package to provide delta coverage between your tests.

## Example

```shell
$ deltacoverage
TestParseCoverageResultErrorsIfGivenNoTests 1.0%
TestParseCoverageResultReturnsCorrectValueGivenTestScriptTestCoverResult 3.0%
TestParseCoverageResultReturnsCorrectValueGivenNoTestScriptsTestCoverResult 0.0%
TestParseTestListReturnsZeroTestsIfGivenNoTests 0.0%
TestParseTestListReturnsCorrectValuesGivenTestListResult 0.0%
TestScript 45.0%
```

The test `TestParseCoverageResultErrorsIfGivenNoTests` contributes to exclusive
1.0% of the total coverage while `TestScript` is responsible for 45%.

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

- POC
- Brute force implementation using `go test` commands (WIP)
- Borrow `cover.go` code and create `deltagecoverage.go`
- (Dream) Implement as Go CLI feature.

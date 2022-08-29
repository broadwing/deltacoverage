# deltacoverage

Go package to provide delta coverage between your tests.

## Example

```shell
$ deltacoverage
TestParseCoverageResult_ErrorsIfGivenNoTests 1.0%
TestParseCoverageResult_ReturnsCorrectValueGivenTestScriptTestCoverResult 3.2%
TestParseCoverageResult_ReturnsCorrectValueGivenNoTestScriptsTestCoverResult 0.0%
TestParseListTests_ReturnsZeroTestsIfGivenNoTests 0.0%
TestParseListTests_ReturnsCorrectValuesGivenTestListResult 0.0%
TestTestScript_ 54.2%
```

The test `TestParseCoverageResultErrorsIfGivenNoTests` contributes to exclusive
1.0% of the total coverage while `TestScript` is responsible for 54.2%.

## Motivation

```text
Herb Derby came up with this metric of “delta coverage”. 
You look at your test suite and you measure what coverage each test adds uniquely that no other test provides.
If you have lots of tests with no delta coverage, so that you could just delete them and not lose your 
ability to exercise parts of the system, then you should delete those tests, unless they have some some 
communication purpose.
—Kent Beck, “Is TDD Dead?”
```

Thank you [@bitfield](https://github.com/bitfield) for the suggestion.

## Roadmap

- POC :white_check_mark:
- Brute force implementation using `go test` commands :white_check_mark:
- Borrow `cover.go` code and create `deltagecoverage.go` -> WIP
- (Dream) Implement as Go CLI feature.

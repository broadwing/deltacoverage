# deltacoverage

Go package to provide delta coverage between your tests.

## Installing

```shell
go install github.com/broadwing/deltacoverage/cmd/deltacoverage@latest
```

## Example

```shell
$ deltacoverage
TestCleanup_RemovesOutputPath 0.7%
TestGenerate_CreatesExpectedCoverProfilesGivenCode 7.2%
TestGenerate_ErrorsGivenCodeWithNoTests 1.4%
TestListTests_SetsExpectedTestsNamesGivenCode 0.0%
TestNewCoverProfile_ErrorsIfPathDoesNotExist 0.7%
TestNewCoverProfile_ErrorsIfPathIsNotDirectory 0.7%
TestNewCoverProfile_SetsCorrectPackagePathGivenString 0.0%
TestNewCoverProfile_SetsNonEmptyOutputPath 0.0%
TestParse_ReturnsExpectedNumberStatementsGivenOutputPath 0.0%
TestParse_ReturnsExpectedTestsBranchesGivenOutputPath 0.0%
TestParse_ReturnsExpectedUniqueBranchesGivenOutputPath 0.0%
TestString_ReturnsExpectedDeltaCoverageGivenCoverProfile 10.1%
```

The test `TestCleanup_RemovesOutputPath` contributes to exclusive
0.7% of the total coverage while
`TestString_ReturnsExpectedDeltaCoverageGivenCoverProfile` is responsible for
10.1%.

Although the tests `TestNewCoverProfile_SetsCorrectPackagePathGivenString` and
`TestNewCoverProfile_SetsNonEmptyOutputPath` have no delta coverage or 0.0%, it is okay
since we are testing behaviours and not functions. The tests run the same code
path but test different things. Also, this design allow ease understand
of what is the test supposed to do.

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
- Get same metrics using coverprofiles of each test instead of brute force. It's going to improve performance from O(n2) to O(n):white_check_mark:
- Borrow `cover.go` code and create `deltagecoverage.go` -> Abandoned
- (Dream) Implement as Go CLI feature

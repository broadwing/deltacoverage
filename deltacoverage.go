package deltacoverage

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type CoverProfile struct {
	TotalStatements int
	Branches        map[string]int
	Tests           map[string][]string
}

func (c CoverProfile) String() string {
	output := []string{}
	for testName, ids := range c.Tests {
		perc := 0.0
		for _, id := range ids {
			perc = float64(c.Branches[id] * 100 / 1)
		}
		output = append(output, fmt.Sprintf("%s %.1f%s", testName, perc, "%"))
	}
	sort.Strings(output)
	return strings.Join(output, "\n")
}

func ParseCoverProfile(coverFilePath string, extraFilePaths ...string) (CoverProfile, error) {
	covProfile := CoverProfile{
		Branches: map[string]int{},
		Tests:    map[string][]string{},
	}
	allCoverFilePaths := append([]string{coverFilePath}, extraFilePaths...)
	for _, coverPath := range allCoverFilePaths {
		testName := strings.Split(path.Base(coverPath), ".")[0]
		f, err := os.Open(coverPath)
		if err != nil {
			return CoverProfile{}, err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		// dicard first line `mode: xxx`
		scanner.Scan()
		// line example with headers
		// identifier          statements visited
		// xyz/xyz.go:3.24,5.2 1          1
		for scanner.Scan() {
			items := strings.Fields(scanner.Text())
			nrStmt, err := strconv.Atoi(items[1])
			if err != nil {
				return CoverProfile{}, err
			}
			_, exists := covProfile.Branches[items[0]]
			if !exists {
				covProfile.Branches[items[0]] = nrStmt
				covProfile.TotalStatements += nrStmt
				covProfile.Tests[testName] = append(covProfile.Tests[testName], items[0])
				continue
			}
			// remove references from other tests
			ids := []string{}
			for testName, tests := range covProfile.Tests {
				for _, id := range tests {
					if id == items[0] {
						continue
					}
					ids = append(ids, id)
				}
				covProfile.Tests[testName] = ids
			}
			// insert empty result
			covProfile.Tests[testName] = []string{}
		}
		if err := scanner.Err(); err != nil {
			return CoverProfile{}, err
		}
	}
	return covProfile, nil
}

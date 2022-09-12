package deltacoverage

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/errgo.v2/errors"
)

var ErrMustBeDirectory = errors.New("Must be a directory")

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

func ParseCoverProfile(dirPath string) (CoverProfile, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return CoverProfile{}, err
	}
	if !info.IsDir() {
		return CoverProfile{}, ErrMustBeDirectory
	}
	covProfile := CoverProfile{
		Branches: map[string]int{},
		Tests:    map[string][]string{},
	}
	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		matched, err := filepath.Match("*.coverprofile", filepath.Base(path))
		if err != nil {
			return err
		}
		if matched {
			testName := strings.Split(filepath.Base(path), ".")[0]
			f, err := os.Open(path)
			if err != nil {
				return err
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
					return err
				}
				_, exists := covProfile.Branches[items[0]]
				if !exists {
					covProfile.Branches[items[0]] = nrStmt
					covProfile.TotalStatements += nrStmt
					covProfile.Tests[testName] = append(covProfile.Tests[testName], items[0])
					continue
				}
				// insert empty result
				covProfile.Tests[testName] = []string{}
				// remove references from other tests
				for testName := range covProfile.Tests {
					covProfile.Tests[testName] = []string{}
				}
			}
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return CoverProfile{}, err
	}
	return covProfile, nil
}

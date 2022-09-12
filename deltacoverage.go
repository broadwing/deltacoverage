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
	UniqueBranches  map[string]int
	Tests           map[string][]string
}

func (c CoverProfile) String() string {
	output := []string{}
	for testName, ids := range c.Tests {
		perc := 0.0
		for _, id := range ids {
			statements, exist := c.UniqueBranches[id]
			if !exist {
				continue
			}
			perc = float64(statements) / float64(c.TotalStatements) * 100
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
		UniqueBranches: map[string]int{},
		Tests:          map[string][]string{},
	}
	branchesCount := map[string]int{}
	branchesStmts := map[string]int{}
	err = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".coverprofile" {
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
				branch := items[0]
				nrStmt, err := strconv.Atoi(items[1])
				if err != nil {
					return err
				}
				visited, err := strconv.Atoi(items[2])
				if err != nil {
					return err
				}
				if visited < 1 {
					continue
				}
				_, exists := branchesCount[branch]
				if !exists {
					covProfile.TotalStatements += nrStmt
					branchesStmts[branch] = nrStmt
				}
				branchesCount[branch]++
				covProfile.Tests[testName] = append(covProfile.Tests[testName], branch)
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
	for branch, times := range branchesCount {
		if times > 1 {
			continue
		}
		covProfile.UniqueBranches[branch] = branchesStmts[branch]
	}
	return covProfile, nil
}

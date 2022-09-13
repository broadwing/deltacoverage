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
var ErrMustBeFile = errors.New("Must be a file")

type CoverProfile struct {
	DirPath          string
	NumberStatements int
	UniqueBranches   map[string]int
	Tests            map[string][]string
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
			perc = float64(statements) / float64(c.NumberStatements) * 100
		}
		output = append(output, fmt.Sprintf("%s %.1f%s", testName, perc, "%"))
	}
	sort.Strings(output)
	return strings.Join(output, "\n")
}

func (c *CoverProfile) ParseCoverProfile() error {
	branchesCount := map[string]int{}
	branchesStmts := map[string]int{}
	err := filepath.Walk(c.DirPath, func(path string, info fs.FileInfo, err error) error {
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
					branchesStmts[branch] = nrStmt
				}
				branchesCount[branch]++
				c.Tests[testName] = append(c.Tests[testName], branch)
			}
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for branch, times := range branchesCount {
		if times > 1 {
			continue
		}
		c.UniqueBranches[branch] = branchesStmts[branch]
	}
	return nil
}

func NewCoverProfile(dirPath string) (*CoverProfile, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return &CoverProfile{}, err
	}
	if !info.IsDir() {
		return &CoverProfile{}, ErrMustBeDirectory
	}
	covProf := &CoverProfile{
		DirPath:        dirPath,
		UniqueBranches: map[string]int{},
		Tests:          map[string][]string{},
	}
	files, err := os.ReadDir(covProf.DirPath)
	if err != nil {
		return &CoverProfile{}, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".coverprofile" {
			nrStatsments, err := ParseTotalStatements(dirPath + "/" + file.Name())
			if err != nil {
				return &CoverProfile{}, err
			}
			covProf.NumberStatements = nrStatsments
			break
		}
	}
	err = covProf.ParseCoverProfile()
	if err != nil {
		return &CoverProfile{}, err
	}
	return covProf, nil
}

func ParseTotalStatements(filePath string) (int, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		return 0, ErrMustBeFile
	}
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// dicard first line `mode: xxx`
	scanner.Scan()
	// line example with headers
	// identifier          statements visited
	// xyz/xyz.go:3.24,5.2 1          1
	nrStmt := 0
	for scanner.Scan() {
		testStmt, err := strconv.Atoi(strings.Fields(scanner.Text())[1])
		if err != nil {
			return 0, err
		}
		nrStmt += testStmt
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return nrStmt, nil
}

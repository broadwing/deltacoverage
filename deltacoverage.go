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

type ProfileItem struct {
	Branch     string
	Statements int
	Visited    bool
}

type CoverProfile struct {
	DirPath          string
	NumberStatements int
	UniqueBranches   map[string]int
	Tests            map[string][]string
}

// ParseProfileLine returns a pointer to ProfileItem type for a given string.
// line example with headers
// identifier          statements visited
// xyz/xyz.go:3.24,5.2 1          1
func (c CoverProfile) ParseProfileLine(line string) (*ProfileItem, error) {
	if strings.HasPrefix(line, "mode:") {
		return &ProfileItem{}, nil
	}
	items := strings.Fields(line)
	branch := items[0]
	profItem := &ProfileItem{
		Branch: branch,
	}
	nrStmt, err := strconv.Atoi(items[1])
	if err != nil {
		return &ProfileItem{}, err
	}
	profItem.Statements = nrStmt
	timesVisited, err := strconv.Atoi(items[2])
	if err != nil {
		return &ProfileItem{}, err
	}
	if timesVisited > 0 {
		profItem.Visited = true
	}
	return profItem, err
}

// String prints out the deltacoverage percentage for each test
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

// ParseCoverProfile reads all files from a directory with extension
// .coverprofile, parses it, and populate CoverProfile object
// Design tradeoff to get the total number of statements:
// 1. Parse one file just to sum the total statements
// 2. Recalculate total statements each iteration on the loop
// 3. Add a counter and an if to just calculate the total sum in the first
// iteration
// Current implementation in number 3
func (c *CoverProfile) ParseCoverProfile() error {
	branchesCount := map[string]int{}
	branchesStmts := map[string]int{}
	profilesRead := 0
	err := filepath.Walk(c.DirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".coverprofile" {
			profilesRead++
			testName := strings.Split(filepath.Base(path), ".")[0]
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				profileItem, err := c.ParseProfileLine(line)
				if err != nil {
					return fmt.Errorf("cannot parse profile line %q: %+v", line, err)
				}
				// only sum total statements when reading the first cover profile
				if profilesRead == 1 {
					c.NumberStatements += profileItem.Statements
				}
				if !profileItem.Visited {
					continue
				}
				_, exists := branchesCount[profileItem.Branch]
				if !exists {
					branchesStmts[profileItem.Branch] = profileItem.Statements
				}
				branchesCount[profileItem.Branch]++
				c.Tests[testName] = append(c.Tests[testName], profileItem.Branch)
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
		if times < 2 {
			c.UniqueBranches[branch] = branchesStmts[branch]
		}
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
	err = covProf.ParseCoverProfile()
	if err != nil {
		return &CoverProfile{}, err
	}
	return covProf, nil
}

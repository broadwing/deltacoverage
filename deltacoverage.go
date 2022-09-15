package deltacoverage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/errgo.v2/errors"
)

var ErrMustBeDirectory = errors.New("Must be a directory")
var ErrMustBeFile = errors.New("Must be a file")

type profileItem struct {
	Branch     string
	Statements int
	Visited    bool
}

type CoverProfile struct {
	DirPath          string
	NumberStatements int
	Tests            map[string][]string
	UniqueBranches   map[string]int
}

// parseProfileLine returns a pointer to ProfileItem type for a given string.
// line example with headers
// identifier          statements visited
// xyz/xyz.go:3.24,5.2 1          1
func (c CoverProfile) parseProfileLine(line string) (*profileItem, error) {
	if strings.HasPrefix(line, "mode:") {
		return &profileItem{}, nil
	}
	items := strings.Fields(line)
	branch := items[0]
	profItem := &profileItem{
		Branch: branch,
	}
	nrStmt, err := strconv.Atoi(items[1])
	if err != nil {
		return &profileItem{}, err
	}
	profItem.Statements = nrStmt
	timesVisited, err := strconv.Atoi(items[2])
	if err != nil {
		return &profileItem{}, err
	}
	if timesVisited > 0 {
		profItem.Visited = true
	}
	return profItem, err
}

// String prints out the deltacoverage percentage for each test
func (c CoverProfile) String() string {
	output := []string{}
	if len(c.Tests) == 0 {
		return "No tests found\n"
	}
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
	return fmt.Sprintf("%s\n", strings.Join(output, "\n"))
}

// parseCoverProfile reads all files from a directory with extension
// .coverprofile, parses it, and populate CoverProfile object
// Design tradeoff to get the total number of statements:
// 1. Parse one file just to sum the total statements
// 2. Recalculate total statements each iteration on the loop
// 3. Add a counter and an if to just calculate the total sum in the first
// iteration
// Current implementation is number 3
func (c *CoverProfile) parseCoverProfile() error {
	branchesCount := map[string]int{}
	branchesStmts := map[string]int{}
	profilesRead := 0
	files, err := os.ReadDir(c.DirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".coverprofile" {
			profilesRead++
			testName := strings.Split(filepath.Base(file.Name()), ".")[0]
			f, err := os.Open(c.DirPath + "/" + file.Name())
			if err != nil {
				return err
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				profileItem, err := c.parseProfileLine(line)
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
				_, ok := branchesCount[profileItem.Branch]
				if !ok {
					branchesStmts[profileItem.Branch] = profileItem.Statements
				}
				branchesCount[profileItem.Branch]++
				c.Tests[testName] = append(c.Tests[testName], profileItem.Branch)
			}
			if err := scanner.Err(); err != nil {
				return err
			}
		}
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
	err = covProf.parseCoverProfile()
	if err != nil {
		return &CoverProfile{}, err
	}
	return covProf, nil
}

// npm-blame captures useful informations and common errors from npm node_modules
package npmblame

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gosuri/uitable"
)

const (
	// The package contains executables
	EXEC_ERROR = iota
	// The package contains test files
	TEST_ERROR
	// The package contains benchmarks files
	BENCH_ERROR
	// The package contains jsx files
	JSX_ERROR
	// The package contains ts files
	TS_ERROR
	// The package contains image files
	IMG_ERROR
	// The package contains TravisCI files
	TRAVIS_ERROR
	// The package contains editorconfig, eslint or sass-lint files
	EDITOR_LINT_ERROR
)

// packages contains maps all the pacakage to there given error codes and the
// numeber of times those occured
type NpmPackages map[string]map[int]int

// NewNpmPackages returns a new npm package instance
func NewNpmPackages() NpmPackages {
	return make(NpmPackages)
}

// ExtractPackageName returns the npm package name from a given path
// TODO improve to handle node_modules
func (np NpmPackages) ExtractPackageName(path string) string {
	dir := filepath.Dir(path)
	module := strings.Split(dir, "/")
	if len(module[0]) > 0 {
		return module[0]
	}
	return module[1]

}

// AppendError appends an error to a given package
func (np NpmPackages) AppendError(pkgName string, err int) {
	if len(np[pkgName]) == 0 {
		np[pkgName] = make(map[int]int)
	}
	np[pkgName][err] = np[pkgName][err] + 1
}

// Blame reports on error for a given npm package
func (np NpmPackages) Blame(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	pkg := np.ExtractPackageName(path)
	if pkg == "" || pkg == ".bin" {
		return nil
	}

	if !info.Mode().IsDir() && (info.Mode()&0111) != 0 {
		np.AppendError(pkg, EXEC_ERROR)
	}

	if strings.Contains(path, "test") || strings.Contains(path, "tests") || strings.Contains(path, ".zuul.yml") || strings.Contains(path, "coverage") || strings.Contains(path, ".coveralls.yml") {
		np.AppendError(pkg, TEST_ERROR)
	}

	if strings.Contains(path, "bench") {
		np.AppendError(pkg, BENCH_ERROR)
	}

	if filepath.Ext(path) == ".jsx" {
		np.AppendError(pkg, JSX_ERROR)
	}

	if filepath.Ext(path) == ".ts" {
		np.AppendError(pkg, TS_ERROR)
	}

	if filepath.Ext(path) == ".png" || filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".ico" {
		np.AppendError(pkg, IMG_ERROR)
	}

	if strings.Contains(path, ".travis.yml") {
		np.AppendError(pkg, TRAVIS_ERROR)
	}

	if strings.Contains(path, ".editorconfig") || strings.Contains(path, ".eslintrc") || strings.Contains(path, ".sass-lint.yml") || strings.Contains(path, ".jshintrc") {
		np.AppendError(pkg, EDITOR_LINT_ERROR)
	}

	return nil
}

// totalErrors return the total amount of errors
func (np NpmPackages) TotalErrors(pkgName string) int {
	totalErrors := 0
	for _, err := range np[pkgName] {
		totalErrors += err
	}
	return totalErrors
}

// DisplayErrors pretty prints the npm packages error report
func (np NpmPackages) String() string {
	buf := &bytes.Buffer{}
	var totalErr int

	var keys []string
	for k := range np {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("PACKAGE", "ERRORS", "EXECUTABLE FILE", "TESTS", "BENCH",
		"HAS_JSX", "HAS_TS", "IMAGES", "TRAVIS_FILES", "EDITOR_LINT_FILES")
	for _, name := range keys {
		errors := np[name]

		if len(errors) > 0 {
			pkgErr := np.TotalErrors(name)
			totalErr += 1
			table.AddRow(name, pkgErr, errors[EXEC_ERROR], errors[TEST_ERROR],
				errors[BENCH_ERROR], errors[JSX_ERROR], errors[TS_ERROR],
				errors[IMG_ERROR], errors[TRAVIS_ERROR], errors[EDITOR_LINT_ERROR])
		}
	}

	fmt.Fprintf(buf, "Your node_modules contains %d packages with errors out of %d packages\n\n", totalErr, len(np))
	fmt.Fprintln(buf, table)
	return buf.String()
}

// npm-blame captures useful informations and common errors from npm node_modules
package main

import (
	"fmt"
	"os"
	"path/filepath"
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

// package contains maps all the pacakage to there given error codes and the
// numeber of times those occured
var packages = make(map[string]map[int]int)

// packageName returns the npm package name from a given path
// TODO improve to handle node_modules
func packageName(path string) string {
	dir := filepath.Dir(path)
	module := strings.Split(dir, "/")
	return module[0]
}

// addError appends an error to a given package
func addError(pkgName string, err int) {
	packages[pkgName][err] = packages[pkgName][err] + 1
}

func blame(path string, info os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	pkg := packageName(path)
	if strings.Contains(pkg, ".bin") {
		return nil
	}

	if len(packages[pkg]) == 0 {
		packages[pkg] = make(map[int]int)
	}

	if !info.Mode().IsDir() && (info.Mode()&0111) != 0 {
		addError(pkg, EXEC_ERROR)
	}

	if strings.Contains(path, "test") || strings.Contains(path, "tests") || strings.Contains(path, ".zuul.yml") || strings.Contains(path, "coverage") || strings.Contains(path, ".coveralls.yml") {
		addError(pkg, TEST_ERROR)
	}

	if strings.Contains(path, "bench") {
		addError(pkg, BENCH_ERROR)
	}

	if filepath.Ext(path) == ".jsx" {
		addError(pkg, JSX_ERROR)
	}

	if filepath.Ext(path) == ".ts" {
		addError(pkg, TS_ERROR)
	}

	if filepath.Ext(path) == ".png" {
		addError(pkg, IMG_ERROR)
	}

	if strings.Contains(path, ".travis.yml") {
		addError(pkg, TRAVIS_ERROR)
	}

	if strings.Contains(path, ".editorconfig") || strings.Contains(path, ".eslintrc") || strings.Contains(path, ".sass-lint.yml") || strings.Contains(path, ".jshintrc") {
		addError(pkg, EDITOR_LINT_ERROR)
	}

	return nil
}

// totalErrors return the total amount of errors for a given package
func totalErrors(errors map[int]int) int {
	totalErrors := 0
	for _, err := range errors {
		totalErrors += err
	}
	return totalErrors
}

func main() {
	if err := filepath.Walk(".", blame); err != nil {
		fmt.Println("File system traversing error.", err)
		os.Exit(-1)
	}
	var totalErr int

	table := uitable.New()
	table.MaxColWidth = 50

	table.AddRow("PACKAGE", "ERRORS", "EXECUTABLE FILE", "TESTS", "BENCH",
		"HAS_JSX", "HAS_TS", "IMAGES", "TRAVIS_FILES", "EDITOR_LINT_FILES")
	for name, errors := range packages {
		if len(errors) > 0 {
			pkgErr := totalErrors(errors)
			totalErr += 1
			table.AddRow(name, pkgErr, errors[EXEC_ERROR], errors[TEST_ERROR],
				errors[BENCH_ERROR], errors[JSX_ERROR], errors[TS_ERROR],
				errors[IMG_ERROR], errors[TRAVIS_ERROR], errors[EDITOR_LINT_ERROR])
		}
	}
	fmt.Printf("Your node_modules contains %d packages with errors out of %d packages\n\n", totalErr, len(packages))
	fmt.Println(table)
}

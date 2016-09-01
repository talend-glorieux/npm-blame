// Package npmblame captures useful informations and common errors from npm node_modules
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

// PackageError represents common npm packages errors
type PackageError int

const (
	// ExecError marks a package containing executables
	ExecError PackageError = iota
	// TestError marks a package test files
	TestError
	// BenchError marks a package benchmark files
	BenchError
	// ImageError marks a package images files
	ImageError
	// CIError marks a package containing continous integration files
	CIError
	// DotfileError marks a package lint files
	DotfileError
)

// NpmPackage represents a npm package
type NpmPackage struct {
	BugsURL  string
	Homepage string
	Errors   map[int]int
}

// NpmPackages is a map of all the packages and there given errors
type NpmPackages map[string]map[PackageError]int

// NewNpmPackages returns a new npm package instance
func NewNpmPackages() NpmPackages {
	return make(NpmPackages)
}

// ExtractPackageName returns the npm package name from a given path
// TODO improve to handle nested dependencies
func (np NpmPackages) ExtractPackageName(path string) string {
	dir := filepath.Dir(path)
	//fmt.Println("HAS_NODEMODULES", strings.Contains(dir, "node_modules"), strings.LastIndex(dir, "node_modules"))
	subPkgIndex := strings.LastIndex(dir, "node_modules")
	var module []string
	if subPkgIndex != -1 {
		module = strings.Split(dir[subPkgIndex:], "/")
		if len(module) > 1 {
			return module[1]
		}
	}
	module = strings.Split(dir, "/")
	if len(module[0]) > 0 {
		return module[0]
	}
	return module[1]
}

//ExtractPackageInformations extracts the package information from its package.json
func (np NpmPackages) ExtractPackageInformations() error {
	return nil
}

// AppendError appends an error to a given package
func (np NpmPackages) AppendError(pkgName string, err PackageError) {
	if len(np[pkgName]) == 0 {
		np[pkgName] = make(map[PackageError]int)
	}
	np[pkgName][err] = np[pkgName][err] + 1
}

func (np NpmPackages) checkTests(path string, pkg string) {
	if strings.Contains(path, "test") ||
		strings.Contains(path, "tests") ||
		strings.Contains(path, ".zuul.yml") ||
		strings.Contains(path, "coverage") ||
		strings.Contains(path, ".coveralls.yml") {
		np.AppendError(pkg, TestError)
	}
}

func (np NpmPackages) checkDotFiles(path string, pkg string) {
	if strings.Contains(path, ".editorconfig") ||
		strings.Contains(path, ".eslintrc") ||
		strings.Contains(path, ".sass-lint.yml") ||
		strings.Contains(path, ".jshintrc") {
		np.AppendError(pkg, DotfileError)
	}
}

func (np NpmPackages) checkExecutables(info os.FileInfo, pkg string) {
	if !info.Mode().IsDir() && (info.Mode()&0111) != 0 {
		np.AppendError(pkg, ExecError)
	}
}

func (np NpmPackages) checkImages(path string, pkg string) {
	if filepath.Ext(path) == ".png" ||
		filepath.Ext(path) == ".jpg" ||
		filepath.Ext(path) == ".ico" {
		np.AppendError(pkg, ImageError)
	}
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
	if np[pkg] == nil {

		np[pkg] = make(map[PackageError]int)
	}

	np.checkExecutables(info, pkg)
	np.checkTests(path, pkg)
	np.checkDotFiles(path, pkg)
	np.checkImages(path, pkg)

	if strings.Contains(path, "bench") {
		np.AppendError(pkg, BenchError)
	}

	if strings.Contains(path, ".travis.yml") {
		np.AppendError(pkg, CIError)
	}

	return nil
}

// TotalErrors return the total amount of errors
func (np NpmPackages) TotalErrors(pkgName string) int {
	totalErrors := 0
	for _, err := range np[pkgName] {
		totalErrors += err
	}
	return totalErrors
}

// String returns the printalbe representation of the NpmPackages
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
		"IMAGES", "TRAVIS_FILES", "EDITOR_LINT_FILES")
	for _, name := range keys {
		errors := np[name]

		if len(errors) > 0 {
			pkgErr := np.TotalErrors(name)
			totalErr++
			table.AddRow(name, pkgErr, errors[ExecError], errors[TestError],
				errors[BenchError], errors[ImageError], errors[CIError],
				errors[DotfileError])
		}
	}

	fmt.Fprintf(buf, "Your node_modules contains %d packages with errors out of %d packages\n\n", totalErr, len(np))
	fmt.Fprintln(buf, table)
	return buf.String()
}

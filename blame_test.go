package npmblame

import (
	"testing"

	"github.com/spf13/afero"
)

func TestNewNpmPackages(t *testing.T) {
	np := NewNpmPackages()

	if len(np) > 0 {
		t.Error("NpmPackage is not empty")
	}
}

func TestExtractPackageName(t *testing.T) {
	np := NewNpmPackages()

	t.Run("regular package", func(t *testing.T) {
		p := np.ExtractPackageName("test/regular/path")
		if p != "test" {
			t.Errorf("Wrong package name: expected test got %s", p)
		}
	})

	t.Run("rooted path", func(t *testing.T) {
		p := np.ExtractPackageName("/test/root/path")
		if p != "test" {
			t.Errorf("Wrong package name expected test got %s", p)
		}
	})

	t.Run("root", func(t *testing.T) {
		p := np.ExtractPackageName("/")
		if p != "" {
			t.Error("Expected an empty string got %s", p)
		}
	})

	t.Run("nested package", func(t *testing.T) {
		t.Skip("skipping test while not implemented")
		if p := np.ExtractPackageName("/test/pkg/nested/path"); p != "nested" {
			t.Error("Wrong package name", p)
		}
	})
}

func BenchmarkExtractPackageName(b *testing.B) {
	np := NewNpmPackages()

	for i := 0; i < b.N; i++ {
		np.ExtractPackageName("test")
	}
}

func TestAppendPackage(t *testing.T) {
	np := NewNpmPackages()

	np.AppendError("test", EXEC_ERROR)
	if np["test"] == nil {
		t.Errorf("Error was not appended: %v", np)
	}
	if np["test"][EXEC_ERROR] == 0 {
		t.Errorf("Wrong error was appended: %v", np)
	}
}

func BenchmarkAppendError(b *testing.B) {
	np := NewNpmPackages()
	for i := 0; i < b.N; i++ {
		np.AppendError("test", EXEC_ERROR)
	}
}

func createNodeModulesFolder() (fs afero.Fs, err error) {
	fs = afero.NewMemMapFs()
	err = fs.Mkdir("/pkg", 0600)

	err = fs.Mkdir("/.bin", 0600)
	fs.Create("/.bin/bin")

	// EXEC_ERROR
	fs.Create("/pkg/exec")
	fs.Chmod("/pkg/exec", 0755)

	// TEST_ERROR
	fs.Create("/pkg/test")

	// BENCH_ERROR
	fs.Mkdir("/pkg/bench", 0600)

	// JSX_ERROR
	fs.Create("/pkg/index.jsx")

	// TS_ERROR
	fs.Create("/pkg/index.ts")

	// IMG_ERROR
	fs.Create("/pkg/favicon.ico")
	fs.Create("/pkg/icon.png")
	fs.Create("/pkg/icon.jpg")

	// TRAVIS_ERROR
	fs.Create("/pkg/.travis.yml")

	// EDITOR_LINT_ERROR
	fs.Create("/pkg/.editorconfig")
	fs.Create("/pkg/.eslintrc")
	fs.Create("/pkg/.sass-lint.yml")
	fs.Create("/pkg/.jshintrc")
	return
}

func TestBlame(t *testing.T) {
	fs, err := createNodeModulesFolder()
	if err != nil {
		t.Error("FileSystem error", err)
	}
	np := NewNpmPackages()
	if err := afero.Walk(fs, "/", np.Blame); err != nil {
		t.Error("Walk Error", err)
	}

	t.Run("Exclude Root", func(t *testing.T) {
		if len(np[""]) > 0 {
			t.Error("Root should not be on the package list")
		}
	})

	t.Run("Exclude binaries", func(t *testing.T) {
		if len(np[".bin"]) > 0 {
			t.Error("Main binary folder should be excluded")
		}
	})

	t.Run("EXEC_ERROR", func(t *testing.T) {
		if np["pkg"][EXEC_ERROR] == 0 {
			t.Error("No EXEC_ERROR", np)
		}
	})

	t.Run("TEST_ERROR", func(t *testing.T) {
		if np["pkg"][TEST_ERROR] == 0 {
			t.Error("No TEST_ERROR", np)
		}
	})

	t.Run("BENCH_ERROR", func(t *testing.T) {
		if np["pkg"][BENCH_ERROR] == 0 {
			t.Error("No BENCH_ERROR", np)
		}
	})

	t.Run("JSX_ERROR", func(t *testing.T) {
		if np["pkg"][JSX_ERROR] == 0 {
			t.Error("No JSX_ERROR", np)
		}
	})

	t.Run("TS_ERROR", func(t *testing.T) {
		if np["pkg"][TS_ERROR] == 0 {
			t.Error("No TS_ERROR", np)
		}
	})

	t.Run("IMG_ERROR", func(t *testing.T) {
		if np["pkg"][IMG_ERROR] != 3 {
			t.Error("No IMG_ERROR", np)
		}
	})

	t.Run("TRAVIS_ERROR", func(t *testing.T) {
		if np["pkg"][TRAVIS_ERROR] == 0 {
			t.Error("No TRAVIS_ERROR", np)
		}
	})

	t.Run("EDITOR_LINT_ERROR", func(t *testing.T) {
		if np["pkg"][EDITOR_LINT_ERROR] != 4 {
			t.Error("No EDITOR_LINT_ERROR", np)
		}
	})
}

func BenchmarkBlame(b *testing.B) {
	fs, err := createNodeModulesFolder()
	if err != nil {
		b.Error("FileSystem Error", err)
	}
	np := NewNpmPackages()
	for i := 0; i < b.N; i++ {
		if err := afero.Walk(fs, "/", np.Blame); err != nil {
			b.Error("FileSystem Walk Error", err)
		}
	}
}

func TestTotalErrors(t *testing.T) {
	np := NewNpmPackages()
	np.AppendError("test", EXEC_ERROR)
	np.AppendError("test1", BENCH_ERROR)
	np.AppendError("test2", IMG_ERROR)
	np.AppendError("test2", IMG_ERROR)

	total := np.TotalErrors("test2")

	if total != 2 {
		t.Errorf("Wrong error total: expected 2 got %d.", total)
	}
}

func BenchmarkTotalErrors(b *testing.B) {
	np := NewNpmPackages()
	np.AppendError("test2", IMG_ERROR)
	np.AppendError("test2", IMG_ERROR)
	for i := 0; i < b.N; i++ {
		np.TotalErrors("test2")
	}
}

func TestString(t *testing.T) {
	np := NewNpmPackages()
	np.AppendError("test", EXEC_ERROR)
	if len(np.String()) == 0 {
		t.Errorf("Expected a non empty string")
	}
}

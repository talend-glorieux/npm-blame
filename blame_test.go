package npmblame

import (
	"fmt"
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
			t.Errorf("Expected an empty string got %s", p)
		}
	})

	t.Run("nested package", func(t *testing.T) {
		if p := np.ExtractPackageName("/test/node_modules/nested/path"); p != "nested" {
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

func TestExtractPackageInformation(t *testing.T) {
	np := NewNpmPackages()
	if err := np.ExtractPackageInformations(); err != nil {
		t.Error(err)
	}
}

func TestAppendPackage(t *testing.T) {
	np := NewNpmPackages()

	np.AppendError("test", ExecError)
	if np["test"] == nil {
		t.Errorf("Error was not appended: %v", np)
	}
	if np["test"][ExecError] == 0 {
		t.Errorf("Wrong error was appended: %v", np)
	}
}

func BenchmarkAppendError(b *testing.B) {
	np := NewNpmPackages()
	for i := 0; i < b.N; i++ {
		np.AppendError("test", ExecError)
	}
}

func createNodeModulesFolder() (fs afero.Fs, err error) {
	fs = afero.NewMemMapFs()
	err = fs.Mkdir("/pkg", 0600)

	err = fs.Mkdir("/.bin", 0600)
	fs.Create("/.bin/bin")

	// ExecError
	fs.Create("/pkg/exec")
	fs.Chmod("/pkg/exec", 0755)

	// TestError
	fs.Create("/pkg/test")

	// BenchError
	fs.Mkdir("/pkg/bench", 0600)

	// ImageError
	fs.Create("/pkg/favicon.ico")
	fs.Create("/pkg/icon.png")
	fs.Create("/pkg/icon.jpg")

	// CIError
	fs.Create("/pkg/.travis.yml")

	// DotfileError
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

	t.Run("Walk Error", func(t *testing.T) {
		if err := np.Blame("", nil, fmt.Errorf("")); err == nil {
			t.Error("Blame should stop in case of errors")
		}
	})

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

	t.Run("ExecError", func(t *testing.T) {
		if np["pkg"][ExecError] == 0 {
			t.Error("No ExecError", np)
		}
	})

	t.Run("TestError", func(t *testing.T) {
		if np["pkg"][TestError] == 0 {
			t.Error("No TestError", np)
		}
	})

	t.Run("BenchError", func(t *testing.T) {
		if np["pkg"][BenchError] == 0 {
			t.Error("No BenchError", np)
		}
	})

	t.Run("ImageError", func(t *testing.T) {
		if np["pkg"][ImageError] != 3 {
			t.Error("No ImageError", np)
		}
	})

	t.Run("CIError", func(t *testing.T) {
		if np["pkg"][CIError] == 0 {
			t.Error("No CIError", np)
		}
	})

	t.Run("DotfileError", func(t *testing.T) {
		if np["pkg"][DotfileError] != 4 {
			t.Error("No DotfileError", np)
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
	np.AppendError("test", ExecError)
	np.AppendError("test1", BenchError)
	np.AppendError("test2", ImageError)
	np.AppendError("test2", ImageError)

	total := np.TotalErrors("test2")

	if total != 2 {
		t.Errorf("Wrong error total: expected 2 got %d.", total)
	}
}

func BenchmarkTotalErrors(b *testing.B) {
	np := NewNpmPackages()
	np.AppendError("test2", ImageError)
	np.AppendError("test2", ImageError)
	for i := 0; i < b.N; i++ {
		np.TotalErrors("test2")
	}
}

func TestString(t *testing.T) {
	np := NewNpmPackages()
	np.AppendError("test", ExecError)
	if len(np.String()) == 0 {
		t.Errorf("Expected a non empty string")
	}
}

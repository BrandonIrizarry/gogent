package functions

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// workingDir returns the absolute path of the project's top level
// directory based on the user's home directory.
func workingDir() (string, error) {
	topLevel := "repos/portfolio/gogent"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, topLevel), nil
}

// TestIgnoredFilesMap checks whether the given files are among those
// ignored. It does this by making sure that they're found in the
// corresponding whitelist map. Because of this, if a subtest passes,
// it's either because the file was ignored (for example because its
// parent directory is ignored), or else the file simply doesn't
// exist.
//
// Again, the listing need not be exhaustive, since project
// maintainers are assumed to later add entries to their .gitignore.
func TestIgnoredFilesMap(t *testing.T) {
	tests := []string{
		".env",
		"finished_notes.org",
		"notes.org",
		"dummy/README.txt",
	}

	wdir, err := workingDir()
	if err != nil {
		t.Fatal(err)
	}

	all, err := allFilesMap(wdir, ".")
	if err != nil {
		t.Fatal(err)
	}

	for _, relpath := range tests {
		path := filepath.Join(wdir, relpath)
		if _, isIncluded := all[path]; isIncluded {
			t.Errorf("%s under %s not being ignored", relpath, wdir)
		}
	}
}

// TestAllFilesMap checks whether the indicated files are among
// the complete listing (that is, it doesn't enumerate all the files
// that need be present.)
func TestAllFilesMap(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	tests := []string{
		"internal/functions/getFileContentRecursively.go",
	}

	wdir, err := workingDir()
	if err != nil {
		t.Fatal(err)
	}

	// LLM functions now work with the canonicalized version of
	// the project subdirectory (as opposed to the relative path.)
	all, err := allFilesMap(wdir, wdir)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ALL: %v", all)

	for _, relpath := range tests {
		path := filepath.Join(wdir, relpath)
		_, isListed := all[path]
		if !isListed {
			t.Errorf("file not listed: %s", path)
		}
	}

	falseTests := []string{"Mxyzptlk"}

	for _, relpath := range falseTests {
		path := filepath.Join(wdir, relpath)
		_, isListed := all[path]

		if isListed {
			t.Errorf("false listing: %s", path)
		}
	}

}

package functions

import (
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
		"finished_notes.org",
		"notes.org",
		"dummy/README.txt",
	}

	wdir, err := workingDir()
	if err != nil {
		t.Fatal(err)
	}

	ignored, err := ignoredFilesMap(wdir)
	if err != nil {
		t.Fatal(err)
	}

	for _, relpath := range tests {
		path := filepath.Join(wdir, relpath)
		untracked := pathIsIgnored(ignored, path)

		if !untracked {
			t.Errorf("%s under %s not being ignored", relpath, wdir)
		}
	}
}

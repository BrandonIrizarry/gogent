package functions

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// workingDir returns the absolute path of the project's top level
// directory based on the user's home directory.
func workingDir(topLevel string) (string, error) {
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
	tests := []struct {
		topLevel string
		want     []string
	}{
		{
			"/repos/gogent/",
			[]string{
				".env",
				"finished_notes.org",
				"notes.org",
				"dummy/README.txt",
			},
		},
	}

	for _, tt := range tests {
		wdir, err := workingDir(tt.topLevel)
		if err != nil {
			t.Fatal(err)
		}

		name := fmt.Sprintf("Project under %s", wdir)
		t.Run(name, func(t *testing.T) {
			all, err := allFilesMap(wdir, ".")
			if err != nil {
				t.Fatal(err)
			}

			for _, entry := range tt.want {
				nentry, err := normalizePath(entry, wdir)
				if err != nil {
					t.Fatal(err)
				}

				if _, isIncluded := all[nentry]; isIncluded {
					t.Errorf("%s under %s not being ignored", entry, wdir)
				}
			}
		})
	}
}

// TestAllFilesMap checks whether the indicated files are among
// the complete listing (that is, it doesn't enumerate all the files
// that need be present.)
func TestAllFilesMap(t *testing.T) {
	tests := []struct {
		topLevel string
		want     []string
	}{
		{
			"/repos/gogent",
			[]string{
				"internal/functions/getFileContentRecursively.go",
			},
		},
	}

	for _, tt := range tests {
		wdir, err := workingDir(tt.topLevel)
		if err != nil {
			t.Fatal(err)
		}

		name := fmt.Sprintf("Project under %s", wdir)
		t.Run(name, func(t *testing.T) {
			all, err := allFilesMap(wdir, ".")
			if err != nil {
				t.Fatal(err)
			}

			for _, sample := range tt.want {
				nsample, err := normalizePath(sample, wdir)
				if err != nil {
					t.Error(err)
				}

				_, isListed := all[nsample]
				if !isListed {
					t.Errorf("File not listed: %s", nsample)
				}
			}
		})
	}

	falseTests := []struct {
		topLevel, dontWant string
	}{
		{
			"/repos/gogent",
			"Mxyzptlk",
		},
	}

	for _, ftt := range falseTests {
		wdir, err := workingDir(ftt.topLevel)
		if err != nil {
			t.Fatal(err)
		}

		name := fmt.Sprintf("Project under %s", wdir)
		t.Run(name, func(t *testing.T) {
			all, err := allFilesMap(wdir, ".")
			if err != nil {
				t.Fatal(err)
			}

			_, isListed := all[ftt.dontWant]
			if isListed {
				t.Errorf("False listing: %s", ftt.dontWant)
			}
		})
	}
}

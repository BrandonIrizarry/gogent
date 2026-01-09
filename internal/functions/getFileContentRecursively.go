package functions

import (
	"fmt"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BrandonIrizarry/gogent/internal/baseconfig"
	"google.golang.org/genai"
)

type getFileContentRecursivelyType struct{}

// getFileContentsRecursively reads contents of all files under a
// given directory. A depth parameter must be specified.
var getFileContentRecursively getFileContentRecursivelyType

func (fnobj getFileContentRecursivelyType) Name() string {
	return "getFileContentRecursively"
}

// ignoredFilesMap returns the set of filenames ignored by the project
// per the project's .gitignore file. Each filename is an absolute
// path.
func ignoredFilesMap(workingDir string) (map[string]bool, error) {
	var bld strings.Builder

	cmd := exec.Command("./ignored.sh")
	cmd.Dir = workingDir
	cmd.Stdout = &bld

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Initialize the 'entries' map.
	//
	// Include the .git directory manually, since git ls-files
	// doesn't list it.
	gitDir, err := normalizePath(".git", workingDir)
	if err != nil {
		return nil, fmt.Errorf("not a Git repository: %s: %w", workingDir, err)
	}

	entries := map[string]bool{
		gitDir: true,
	}

	for e := range strings.SplitSeq(bld.String(), "\n") {
		// Splitting creates empty-string entries, which later
		// get confused as referring to the top-level
		// directory. So we must skip them here.
		//
		// FIXME: defend against this where it gets
		// seen by other functions.
		if e == "" {
			continue
		}

		ne, err := normalizePath(e, workingDir)
		if err != nil {
			return nil, err
		}

		entries[ne] = true
	}

	return entries, nil
}

// allFilesMap walks the filesystem starting at dir (an absolute path)
// and returns a set of absolute pathnames corresponding to files
// underneath dir. This function uses ignoreFilesMap to avoid walking
// down certain directories.
func allFilesMap(workingDir, dir string) (map[string]bool, error) {
	allFiles := make(map[string]bool)
	ignored, err := ignoredFilesMap(workingDir)
	if err != nil {
		return nil, err
	}

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		// It's a good idea to check 'd' for a nil value,
		// since it's possible that, for example, 'dir' was
		// malformed by some previous code and therefore the
		// current 'path' argument refers to a file or
		// directory that doesn't exist.
		if d == nil {
			return fmt.Errorf("nil direntry object for path %s", path)
		}

		_, parentIsIgnored := ignored[filepath.Dir(path)]
		if parentIsIgnored {
			log.Printf("Skipping because parent is ignored: %s", path)
			return filepath.SkipDir
		}

		// FIXME: for now we check for "regular files", though
		// I'm not 100% sure this is what will always be
		// sufficient.
		if d.Type().IsRegular() {
			_, fileIsIgnored := ignored[path]
			if fileIsIgnored {
				log.Printf("Skipping because file is ignored: %s", path)
			} else {
				allFiles[path] = true
			}
		}

		return nil
	})

	return allFiles, nil

}

func (fnobj getFileContentRecursivelyType) Function() functionType {
	return func(args map[string]any, baseCfg baseconfig.BaseConfig) *genai.Part {
		dir, err := normalizePath(args["dir"], baseCfg.WorkingDir)
		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		all, err := allFilesMap(baseCfg.WorkingDir, dir)
		if err != nil {
			return ResponseError(fnobj.Name(), err.Error())
		}

		var bld strings.Builder
		for path := range all {
			content, logs, err := fileContent(path, baseCfg.MaxFilesize)
			if err != nil {
				return ResponseError(fnobj.Name(), err.Error())
			}

			for _, lg := range logs {
				log.Println(lg)
			}

			fmt.Fprintf(&bld, "Contents of %s: %s\n\n", path, content)
		}

		return responseOK(fnobj.Name(), bld.String())
	}
}

package workingdir

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

func (cfg localConfig) SelectWorkingDir() (string, error) {
	var err error

	fp := filepicker.New()
	fp.CurrentDirectory, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Configure the file picker to only allow selecting
	// directories. Also, don't display permissions, size, etc.
	fp.FileAllowed = false
	fp.DirAllowed = true
	fp.ShowPermissions = false
	fp.ShowSize = false
	fp.Styles.EmptyDirectory = fp.Styles.EmptyDirectory.SetString("Directory is empty")

	m := cfg.NewModel(fp)

	// FIXME: for now, don't look at any errors from running the
	// file picker.
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(model)

	return mm.selectedDir, nil
}

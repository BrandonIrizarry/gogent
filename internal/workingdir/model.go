package workingdir

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filepicker  filepicker.Model
	selectedDir string

	// Field quitting is necessary for View to return an empty string,
	// so that junk isn't left behind in the terminal after
	// exiting.
	quitting bool
}

func NewModel(fp filepicker.Model) model {
	return model{
		filepicker: fp,
	}
}

// Init implements the tea.Model interface's Init method.
func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

// Update implements the tea.Model interface's Update method.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)
	didSelect, path := m.filepicker.DidSelectFile(msg)

	// We're not using any file-extension filtering features, so
	// we expect 'didSelect' to always be true.
	if !didSelect {
		logger.Printf("didn't select a file;  path=%s", path)
	} else {
		m.selectedDir = path
		logger.Printf("did select a file; path=%s", path)
	}

	logger.Printf("inside Update: %s", m.selectedDir)
	return m, cmd
}

// View implements the tea.Model interface's View method.
func (m model) View() string {
	logger.Printf("inside View: %s", m.selectedDir)
	if m.quitting {
		return ""
	}

	var view strings.Builder
	fmt.Fprintf(&view, "\nPick a directory:\n\n%s\n", m.filepicker.View())

	return view.String()
}

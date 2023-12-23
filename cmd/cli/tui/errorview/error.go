package errorview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
)

// model is the model for an error view.
// The error view is used to display an error if one occurs.
type model struct {
	activeView msgs.View // which view is currently being shown
	err        error     // the error that occurred
}

// InitialModel builds the model.
func InitialModel() tea.Model {
	return model{
		activeView: msgs.ViewBooks,
	}
}

// Update handles a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case msgs.ViewMsg:
		m.activeView = msgs.View(msg)

	case msgs.ErrorMsg:
		m.err = msg.Err
		return m, func() tea.Msg { return msgs.ViewMsg(msgs.ViewError) }
	}

	return m, nil
}

// View renders the model.
func (m model) View() string {
	if m.activeView != msgs.ViewError {
		return ""
	}

	return fmt.Sprintf("Error: %s", m.err)
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return nil
}

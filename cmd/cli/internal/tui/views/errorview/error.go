package errorview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/msgs"
)

// model is the model for an error view.
// The error view is used to display an error if one occurs.
type model struct {
	err   error // the error that occurred
	width int   // the current width of the screen
}

// InitialModel builds the model.
func InitialModel() tea.Model {
	return model{}
}

// Update handles a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case msgs.ErrorMsg:
		m.err = msg.Err
		return m, func() tea.Msg { return msgs.ViewMsg(msgs.ViewError) }

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the model.
func (m model) View() string {
	return lipgloss.
		NewStyle().
		Width(m.width).
		Render(fmt.Sprintf("Error: %s", m.err))
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return nil
}

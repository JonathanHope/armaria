package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonathanhope/armaria/cmd/cli/tui/booksview"
	"github.com/jonathanhope/armaria/cmd/cli/tui/errorview"
)

// model is an app level model that has multiple view models.
// Only one view model is ever active at a time.
type model struct {
	books tea.Model // view to list books
	error tea.Model // view to show an error
}

// Update handles a message.
// The message is passed down to the underlying view models.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var booksCmd tea.Cmd
	m.books, booksCmd = m.books.Update(msg)

	var errorCmd tea.Cmd
	m.error, errorCmd = m.error.Update(msg)

	return m, tea.Batch([]tea.Cmd{booksCmd, errorCmd}...)
}

// View renders the model.
// This renders and every underlying view model.
// However only one underlying view model will actually have content at a time.
func (m model) View() string {
	return m.books.View() + m.error.View()
}

// Init initializes the model.
// This calls Init on the underlying view models.
func (m model) Init() tea.Cmd {
	booksCmd := m.books.Init()
	errorCmd := m.error.Init()

	return tea.Batch([]tea.Cmd{booksCmd, errorCmd}...)
}

// Program is the top level TUI program.
// This is how the TUI is started.
func Program() *tea.Program {
	return tea.NewProgram(model{
		books: booksview.InitialModel(),
		error: errorview.InitialModel(),
	}, tea.WithAltScreen())
}

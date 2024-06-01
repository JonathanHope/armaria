package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/views/booksview"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/views/errorview"
)

// model is an app level model that has multiple view models.
// Only one view model is ever active at a time.
type model struct {
	activeView msgs.View // which view is currently active
	books      tea.Model // view to list books
	error      tea.Model // view to show an error
}

// Update handles a message.
// The message is passed down to the underlying view models.
// The exception to this is keypresses which are only passed to the current view.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var booksCmd tea.Cmd
	var errorCmd tea.Cmd

	if _, ok := msg.(tea.KeyMsg); ok {
		switch m.activeView {
		case msgs.ViewBooks:
			m.books, booksCmd = m.books.Update(msg)
		case msgs.ViewError:
			m.error, errorCmd = m.error.Update(msg)
		}
	} else {
		m.books, booksCmd = m.books.Update(msg)
		m.error, errorCmd = m.error.Update(msg)
	}

	switch msg := msg.(type) {

	case msgs.ViewMsg:
		m.activeView = msgs.View(msg)
	}

	return m, tea.Batch(booksCmd, errorCmd)
}

// View renders the model.
// This renders and every underlying view model.
// Only the currently active view is rendered.
func (m model) View() string {
	switch m.activeView {
	case msgs.ViewBooks:
		return m.books.View()
	case msgs.ViewError:
		return m.error.View()
	}

	return m.books.View() + m.error.View()
}

// Init initializes the model.
// This calls Init on the underlying view models.
func (m model) Init() tea.Cmd {
	booksCmd := m.books.Init()
	errorCmd := m.error.Init()

	return tea.Batch(booksCmd, errorCmd)
}

// Program is the top level TUI program.
// This is how the TUI is started.
func Program() *tea.Program {
	return tea.NewProgram(model{
		activeView: msgs.ViewBooks,
		books:      booksview.InitialModel(),
		error:      errorview.InitialModel(),
	}, tea.WithAltScreen())
}

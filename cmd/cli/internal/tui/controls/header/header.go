package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/utils"
)

// HeaderModel is the model for a header.
// The header displays state information such as breadcrumbs for the selected book.
type HeaderModel struct {
	name        string // name of the header
	title       string // title of the app
	breadcrumbs string // breadcrumbs for the currently selected book
	busy        bool   // if true the writer is busy
	width       int    // max width of the header
}

// Busy returns whether the writer is busy or not.
func (m HeaderModel) Busy() bool {
	return m.busy
}

// SetBusy denotes that the writer is busy.
func (m *HeaderModel) SetBusy() {
	m.busy = true
}

// SetBusy denotes that the writer is free.
func (m *HeaderModel) SetFree() {
	m.busy = false
}

// SetBreadcrumbs sets the bread crumbs displayed in the header.
func (m *HeaderModel) SetBreadcrumbs(breadcrumbs string) {
	m.breadcrumbs = breadcrumbs
}

// Resize changes the size of the header.
func (m *HeaderModel) Resize(width int) {
	m.width = width
}

// InitialModel builds the model.
func InitialModel(name string, title string) HeaderModel {
	return HeaderModel{
		name:  name,
		title: title,
	}
}

// Update handles a message.
func (m HeaderModel) Update(msg tea.Msg) (HeaderModel, tea.Cmd) {
	return m, nil
}

// View renders the model.
func (m HeaderModel) View() string {
	const cellPadding = 1

	cellWidth := m.width / 2
	cellTextWidth := cellWidth - cellPadding*2

	title := m.title
	if m.busy {
		title += " - ⌛"
	}

	rows := [][]string{
		{title, utils.Substr(m.breadcrumbs, cellTextWidth)},
	}

	titleNavStyle := lipgloss.
		NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		Width(cellWidth).
		Padding(0, cellPadding)

	headingTable := table.
		New().
		Border(lipgloss.ThickBorder()).
		BorderTop(true).
		BorderLeft(false).
		BorderRight(false).
		BorderBottom(true).
		BorderColumn(false).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("5"))).
		Width(m.width).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 1 && col == 0:
				return titleNavStyle.Align(lipgloss.Left)
			case row == 1 && col == 1:
				return titleNavStyle.Align(lipgloss.Right)
			}

			return lipgloss.NewStyle()
		}).Rows(rows...)

	return headingTable.Render()
}

// Init initializes the model.
func (m HeaderModel) Init() tea.Cmd {
	return nil
}

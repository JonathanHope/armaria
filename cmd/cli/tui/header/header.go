package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
	"github.com/samber/lo"
)

// model is the model for a header.
// The header displays state information such as breadcrumbs for the selected book.
type model struct {
	name    string   // name of the header
	title   string   // title of the app
	nav     string   // breadcrumbs for the currently selected book
	width   int      // max width of the header
	filters []string // currently active filters
}

// InitialModel builds the model.
func InitialModel(name string, title string) tea.Model {
	return model{
		name:  name,
		title: title,
	}
}

// Update handles a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case msgs.SizeMsg:
		if msg.Name == m.name {
			m.width = msg.Width
		}

	case msgs.NavMsg:
		m.nav = string(msg)

	case msgs.FiltersMsg:
		m.filters = msg

	}

	return m, nil
}

// View renders the model.
func (m model) View() string {
	const cellPadding = 1

	cellWidth := m.width / 2
	cellTextWidth := cellWidth - cellPadding*2

	rows := [][]string{
		{m.title, utils.Substr(m.nav, 0, cellTextWidth)},
	}

	for _, filtersChunk := range lo.Chunk(m.filters, 2) {
		row := make([]string, 2)

		if len(filtersChunk) > 0 {
			row[0] = utils.Substr(m.filters[0], 0, cellTextWidth)
		}

		if len(filtersChunk) > 1 {
			row[1] = utils.Substr(m.filters[1], 0, cellTextWidth)
		}

		rows = append(rows, row)
	}

	titleNavStyle := lipgloss.
		NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		Width(cellWidth).
		Padding(0, cellPadding)

	filterStyle := lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("2")).
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
			case col == 0:
				return filterStyle.Align(lipgloss.Left)
			case col == 1:
				return filterStyle.Align(lipgloss.Right)
			}

			return lipgloss.NewStyle()
		}).Rows(rows...)

	return headingTable.Render()
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return nil
}

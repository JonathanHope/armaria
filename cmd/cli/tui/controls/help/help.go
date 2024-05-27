package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/samber/lo"
)

// Binding describes a keybinding.
type Binding struct {
	Context string // context the binding is valid in
	Key     string // which key
	Help    string // what the key does
}

// HelpModel is the HelpModel for help.
// The help screen shows keybindings.
type HelpModel struct {
	name     string    // name of the help screen
	contexts []string  // the different context to show keybindings for
	bindings []Binding // the different keybindings
	helpMode bool      // if true the help screen is active
}

// HelpMode returns true if the hlep mode is active.
func (m HelpModel) HelpMode() bool {
	return m.helpMode
}

func (m *HelpModel) ShowHelp() {
	m.helpMode = true
}

func (m *HelpModel) HideHelp() {
	m.helpMode = false
}

// InitialModel builds the model.
func InitialModel(name string, contexts []string, bindings []Binding) HelpModel {
	return HelpModel{
		name:     name,
		contexts: contexts,
		bindings: bindings,
	}
}

// Update handles a message.
func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	return m, nil
}

// View renders the model.
func (m HelpModel) View() string {
	headers := []string{""}
	headers = append(headers, m.contexts...)

	keys := lo.Map(m.bindings, func(x Binding, _ int) string {
		return x.Key
	})
	uniqueKeys := lo.Uniq(keys)

	rows := [][]string{}
	for _, key := range uniqueKeys {
		row := make([]string, len(headers))
		row[0] = key

		for i, context := range m.contexts {
			b := lo.Filter(m.bindings, func(x Binding, _ int) bool {
				return x.Key == key && x.Context == context
			})

			if len(b) > 0 {
				row[i+1] = b[0].Help
			}
		}

		rows = append(rows, row)
	}

	const headerRow = 0
	const headerCol = 0
	const cellPadding = 1

	baseStyle := lipgloss.
		NewStyle().
		Padding(0, cellPadding)

	headerStyle := baseStyle.
		Bold(true).
		Foreground(lipgloss.Color("3"))

	helpTable := table.
		New().
		Headers(headers...).
		Border(lipgloss.HiddenBorder()).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == headerCol || row == headerRow {
				return headerStyle
			}
			return baseStyle
		})

	const tablePadding = 1

	return lipgloss.
		NewStyle().
		Padding(0, tablePadding).
		SetString(helpTable.String()).
		Render()
}

// Init initializes the model.
func (m HelpModel) Init() tea.Cmd {
	return nil
}

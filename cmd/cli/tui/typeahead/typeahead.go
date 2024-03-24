package typeahead

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/scrolltable"
	"github.com/jonathanhope/armaria/cmd/cli/tui/textinput"
	"github.com/samber/lo"
)

// typeaheadTable is the type of the underlying table in the typeahead.
type typeaheadTable = scrolltable.ScrolltableModel[msgs.TypeaheadItem]

// TypeaheadModel is the model for a typeahead.
// The typpeahead allows the user to filter a list of things to select from by typing.
type TypeaheadModel struct {
	name            string                   // name of the typeahead
	width           int                      // max width of the typeahead
	typeaheadMode   bool                     // if true typeahead is accepting input
	inputName       string                   // the name of the input in the typeahead
	tableName       string                   // the name of the table in the typeahead
	unfilteredQuery msgs.UnfilteredQueryFn   // returns results when there isn't enough input
	filteredQuery   msgs.FilteredQueryFn     // returns results when there's enough input
	minFilterChars  int                      // the minumum number of chars needed to filter
	operation       string                   // the operation the typeahead is for
	includeInput    bool                     // if true include the current input as an option
	input           textinput.TextInputModel // allows text input
	table           typeaheadTable           // shows the options to select from
}

// TypeaheadMode returns whether the typeahead is accepting input or not.
func (m TypeaheadModel) TypeaheadMode() bool {
	return m.typeaheadMode
}

// Text returns the text currently in the typeahead.
func (m TypeaheadModel) Text() string {
	return m.input.Text()
}

// InitialModel builds the model.
func InitialModel(name string) TypeaheadModel {
	inputName := name + "Input"
	tableName := name + "Table"

	return TypeaheadModel{
		name:      name,
		inputName: inputName,
		tableName: tableName,
		input:     textinput.InitialModel(inputName, ""),
		table: scrolltable.InitialModel[msgs.TypeaheadItem](
			tableName,
			true,
			[]scrolltable.ColumnDefinition[msgs.TypeaheadItem]{
				{
					Mode:   scrolltable.DynamicColumn,
					Header: "",
					RenderCell: func(item msgs.TypeaheadItem) string {
						return item.Label
					},
					Style: func(item msgs.TypeaheadItem, isSelected bool, isHeader bool) lipgloss.Style {
						style := lipgloss.
							NewStyle()

						if isSelected {
							style = style.Bold(true).Underline(true)
						}

						return style
					},
				},
			},
		),
	}
}

// Update handles a message.
func (m TypeaheadModel) Update(msg tea.Msg) (TypeaheadModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case msgs.InputChangedMsg:
		if msg.Name == m.inputName {
			cmds = append(cmds, m.loadItemsCmd())
		}

	case msgs.SizeMsg:
		if msg.Name == m.name {
			m.width = msg.Width

			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msgs.SizeMsg{
				Name:  m.inputName,
				Width: m.width,
			})
			var tableCmd tea.Cmd
			m.table, tableCmd = m.table.Update(msgs.SizeMsg{
				Name:   m.tableName,
				Width:  m.width,
				Height: msg.Height,
			})
			cmds = append(cmds, inputCmd, tableCmd)
		}

	case msgs.TypeaheadModeMsg:
		if msg.Name == m.name {
			m.typeaheadMode = msg.InputMode
			m.unfilteredQuery = msg.UnfilteredQuery
			m.filteredQuery = msg.FilteredQuery
			m.minFilterChars = msg.MinFilterChars
			m.operation = msg.Operation
			m.includeInput = msg.IncludeInput

			if msg.InputMode {
				cmds = append(cmds, func() tea.Msg {
					return msgs.PromptMsg{Name: m.inputName, Prompt: msg.Prompt}
				}, func() tea.Msg {
					return msgs.TextMsg{Name: m.inputName, Text: msg.Text}
				}, func() tea.Msg {
					return msgs.FocusMsg{Name: m.inputName, MaxChars: msg.MaxChars}
				}, m.loadItemsCmd())
			} else {
				cmds = append(cmds, func() tea.Msg {
					return msgs.BlurMsg{Name: m.inputName}
				}, func() tea.Msg {
					return msgs.PromptMsg{Name: m.inputName, Prompt: ""}
				}, func() tea.Msg {
					return msgs.TextMsg{Name: m.inputName, Text: ""}
				})
			}
		}

	case tea.KeyMsg:
		if m.typeaheadMode {
			switch msg.String() {
			case "ctrl+c":
				if m.typeaheadMode {
					return m, tea.Quit
				}

			case "esc":
				cmds = append(cmds, func() tea.Msg {
					return msgs.TypeaheadCancelledMsg{Name: m.name}
				})

			case "enter":
				cmds = append(cmds, func() tea.Msg {
					return msgs.TypeaheadConfirmedMsg{
						Name:      m.name,
						Value:     m.table.Selection(),
						Operation: m.operation,
					}
				})

			case "up":
				var tableCmd tea.Cmd
				m.table, tableCmd = m.table.Update(msg)
				cmds = append(cmds, tableCmd)

			case "down":
				var tableCmd tea.Cmd
				m.table, tableCmd = m.table.Update(msg)
				cmds = append(cmds, tableCmd)

			default:
				var inputCmd tea.Cmd
				m.input, inputCmd = m.input.Update(msg)
				cmds = append(cmds, inputCmd)
			}
		}

	default:
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)

		var tableCmd tea.Cmd
		m.table, tableCmd = m.table.Update(msg)

		cmds = append(cmds, inputCmd, tableCmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m TypeaheadModel) View() string {
	return m.input.View() + "\n\n" + m.table.View()
}

// Init initializes the model.
func (m TypeaheadModel) Init() tea.Cmd {
	return nil
}

// loadItemsCmd loads the available option in the typeahead.
func (m TypeaheadModel) loadItemsCmd() tea.Cmd {
	return func() tea.Msg {
		if len(strings.Split(m.input.Text(), "")) >= m.minFilterChars {
			items, err := m.filteredQuery(m.input.Text())
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}

			numMatch := len(lo.Filter(items, func(item msgs.TypeaheadItem, index int) bool {
				return item.Label == m.input.Text()
			}))

			if m.includeInput && numMatch == 0 {
				items = append(
					[]msgs.TypeaheadItem{{Label: m.input.Text(), Value: m.input.Text(), New: true}},
					items...)
			}

			return msgs.DataMsg[msgs.TypeaheadItem]{Name: m.tableName, Data: items, Move: msgs.DirectionStart}

		} else {
			items, err := m.unfilteredQuery()
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}

			numMatch := len(lo.Filter(items, func(item msgs.TypeaheadItem, index int) bool {
				return item.Label == m.input.Text()
			}))

			if m.includeInput && m.input.Text() != "" && numMatch == 0 {
				items = append(
					[]msgs.TypeaheadItem{{Label: m.input.Text(), Value: m.input.Text(), New: true}},
					items...)
			}

			return msgs.DataMsg[msgs.TypeaheadItem]{Name: m.tableName, Data: items, Move: msgs.DirectionStart}
		}
	}
}

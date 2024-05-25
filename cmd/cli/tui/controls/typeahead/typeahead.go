package typeahead

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/tui/controls/scrolltable"
	"github.com/jonathanhope/armaria/cmd/cli/tui/controls/textinput"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/samber/lo"
)

// TypeaheadItem is the model for an item in the typeahead.
type TypeaheadItem struct {
	Value string // hidden identifier
	Label string // visible text
	New   bool   // if true this is a new item that didn't exist before
}

// UnfilteredQueryFn is a function that returns typeahead items when no filter is active.
type UnfilteredQueryFn func() ([]TypeaheadItem, error)

// FilteredQueryFn is a function that returns the typeahead items when a filter is active.
type FilteredQueryFn func(query string) ([]TypeaheadItem, error)

// typeaheadTable is the type of the underlying table in the typeahead.
type typeaheadTable = scrolltable.ScrolltableModel[TypeaheadItem]

// TypeaheadTogglePayload has the data needed to start or stop the typeahead.
type StartTypeaheadPayload struct {
	UnfilteredQuery UnfilteredQueryFn // returns results when there isn't enough input
	FilteredQuery   FilteredQueryFn   // returns results when there's enough input
	MinFilterChars  int               // the minumum number of chars needed to filter
	IncludeInput    bool              // if true include the current input as an option
	Prompt          string            // the prompt to show
	Text            string            // the text to start the input with
	MaxChars        int               // the maximum number of chars to allow
}

// TypeaheadModel is the model for a typeahead.
// The typpeahead allows the user to filter a list of things to select from by typing.
type TypeaheadModel struct {
	name            string                   // name of the typeahead
	width           int                      // max width of the typeahead
	typeaheadMode   bool                     // if true typeahead is accepting input
	inputName       string                   // the name of the input in the typeahead
	tableName       string                   // the name of the table in the typeahead
	unfilteredQuery UnfilteredQueryFn        // returns results when there isn't enough input
	filteredQuery   FilteredQueryFn          // returns results when there's enough input
	minFilterChars  int                      // the minumum number of chars needed to filter
	includeInput    bool                     // if true include the current input as an option
	input           textinput.TextInputModel // allows text input
	table           typeaheadTable           // shows the options to select from
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
		table: scrolltable.InitialModel[TypeaheadItem](
			tableName,
			true,
			[]scrolltable.ColumnDefinition[TypeaheadItem]{
				{
					Mode:   scrolltable.DynamicColumn,
					Header: "",
					RenderCell: func(item TypeaheadItem) string {
						return item.Label
					},
					Style: func(item TypeaheadItem, isSelected bool, isHeader bool) lipgloss.Style {
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

// TypeaheadMode returns whether the typeahead is accepting input or not.
func (m TypeaheadModel) TypeaheadMode() bool {
	return m.typeaheadMode
}

// TableName returns the name of the underlying table.
func (m TypeaheadModel) TableName() string {
	return m.name + "Table"
}

// TableName returns the name of the underlying table.
func (m TypeaheadModel) InputName() string {
	return m.name + "Input"
}

// Selection returns the currently selected item.
func (m TypeaheadModel) Selection() TypeaheadItem {
	return m.table.Selection()
}

// StartTypeahead will have the typeahead start collecting input
func (m *TypeaheadModel) StartTypeahead(payload StartTypeaheadPayload) tea.Cmd {
	m.typeaheadMode = true
	m.unfilteredQuery = payload.UnfilteredQuery
	m.filteredQuery = payload.FilteredQuery
	m.minFilterChars = payload.MinFilterChars
	m.includeInput = payload.IncludeInput

	m.input.SetPrompt(payload.Prompt)
	m.input.SetText(payload.Text)
	m.input.Focus(payload.MaxChars)

	return m.LoadItemsCmd()
}

// StartTypeahead will have the typeahead stop collecting input
func (m *TypeaheadModel) StopTypeahead() {
	m.typeaheadMode = false
	m.input.SetPrompt("")
	m.input.SetText("")
	m.input.Blur()
}

// Insert inserts runes in front of the cursor.
func (m *TypeaheadModel) Insert(runes []rune) tea.Cmd {
	return m.input.Insert(runes)
}

// Delete deletes the rune in front of the cursor.
func (m *TypeaheadModel) Delete() tea.Cmd {
	return m.input.Delete()
}

// Resize changes the size of the typeahead.
func (m *TypeaheadModel) Resize(width int, height int) {
	m.width = width
	m.table.Resize(width, height)
	m.input.Resize(width)
}

// MoveUp moves the cursor up the table.
func (m *TypeaheadModel) MoveUp() tea.Cmd {
	return m.table.MoveUp()
}

// MoveDown moves the cursor down the table.
func (m *TypeaheadModel) MoveDown() tea.Cmd {
	return m.table.MoveDown()
}

// MoveLeft moves the cursor the left once space.
func (m *TypeaheadModel) MoveLeft() {
	m.input.MoveLeft()
}

// MoveRight moves the cursor to right once space.
func (m *TypeaheadModel) MoveRight() {
	m.input.MoveLeft()
}

// Reload reloads that data in the table.
func (m *TypeaheadModel) Reload(data []TypeaheadItem, move msgs.Direction) tea.Cmd {
	return m.table.Reload(data, move)
}

// Update handles a message.
func (m TypeaheadModel) Update(msg tea.Msg) (TypeaheadModel, tea.Cmd) {
	return m, nil
}

// View renders the model.
func (m TypeaheadModel) View() string {
	return m.input.View() + "\n\n" + m.table.View()
}

// Init initializes the model.
func (m TypeaheadModel) Init() tea.Cmd {
	return nil
}

// LoadItemsCmd loads the available option in the typeahead.
func (m TypeaheadModel) LoadItemsCmd() tea.Cmd {
	return func() tea.Msg {
		if len(strings.Split(m.input.Text(), "")) >= m.minFilterChars {
			items, err := m.filteredQuery(m.input.Text())
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}

			numMatch := len(lo.Filter(items, func(item TypeaheadItem, index int) bool {
				return item.Label == m.input.Text()
			}))

			if m.includeInput && numMatch == 0 {
				items = append(
					[]TypeaheadItem{{Label: m.input.Text(), Value: m.input.Text(), New: true}},
					items...)
			}

			return msgs.DataMsg[TypeaheadItem]{Name: m.tableName, Data: items, Move: msgs.DirectionStart}
		} else {
			items, err := m.unfilteredQuery()
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}

			numMatch := len(lo.Filter(items, func(item TypeaheadItem, index int) bool {
				return item.Label == m.input.Text()
			}))

			if m.includeInput && m.input.Text() != "" && numMatch == 0 {
				items = append(
					[]TypeaheadItem{{Label: m.input.Text(), Value: m.input.Text(), New: true}},
					items...)
			}

			return msgs.DataMsg[TypeaheadItem]{Name: m.tableName, Data: items, Move: msgs.DirectionStart}
		}
	}
}

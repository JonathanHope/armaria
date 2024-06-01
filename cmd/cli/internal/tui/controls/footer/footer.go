package footer

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/controls/textinput"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/utils"
)

const HelpInfoWidth = 7 // width of the help info in the footer

// FooterModel is the model for a header.
// The footer can collect input, and displays informationa about the apps state.
type FooterModel struct {
	name      string                   // name of the footer
	width     int                      // max width of the footer
	inputMode bool                     // if true footer is accepting input
	filters   []string                 // currently applied filters
	inputName string                   // the name of the input in the footer
	input     textinput.TextInputModel // allows text input
}

// Text returns the text in the footers input.
func (m FooterModel) Text() string {
	return m.input.Text()
}

// InputMode returns whether the footer is accepting input or not.
func (m FooterModel) InputMode() bool {
	return m.inputMode
}

// InputName is the name of the underlying input.
func (m FooterModel) InputName() string {
	return m.name + "Input"
}

// Resize changes the size of the footer.
func (m *FooterModel) Resize(width int) {
	m.width = width
	m.input.Resize(width - HelpInfoWidth)
}

// Insert inserts runes in front of the cursor.
func (m *FooterModel) Insert(runes []rune) tea.Cmd {
	return m.input.Insert(runes)
}

// Delete deletes the rune in front of the cursor.
func (m *FooterModel) Delete() tea.Cmd {
	return m.input.Delete()
}

// MoveLeft moves the cursor the left once space.
func (m *FooterModel) MoveLeft() {
	m.input.MoveLeft()
}

// MoveRight moves the cursor to right once space.
func (m *FooterModel) MoveRight() {
	m.input.MoveRight()
}

// StartInputMode switches the footer into accepting input.
func (m *FooterModel) StartInputMode(prompt string, text string, maxChars int) {
	m.inputMode = true
	m.input.SetPrompt(prompt)
	m.input.SetText(text)
	m.input.Focus(maxChars)
}

// StopInputMode switches the footer out of accepting input.
func (m *FooterModel) StopInputMode() {
	m.inputMode = false
	m.input.SetPrompt("")
	m.input.SetText("")
	m.input.Blur()
}

// SetFilters sets the curently applied filters.
func (m *FooterModel) SetFilters(filters []string) {
	m.filters = filters
}

// InitialModel builds the model.
func InitialModel(name string) FooterModel {
	inputName := name + "Input"

	return FooterModel{
		name:      name,
		inputName: inputName,
		input:     textinput.InitialModel(inputName, ""),
	}
}

// Update handles a message.
func (m FooterModel) Update(msg tea.Msg) (FooterModel, tea.Cmd) {
	return m, nil
}

// View renders the model.
func (m FooterModel) View() string {
	help := lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("3")).
		SetString("Help: ?")

	footerStyle := lipgloss.
		NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderTop(true).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Width(m.width).
		BorderForeground(lipgloss.Color("5"))

	var filters string
	if len(m.filters) > 0 {
		filters = strings.Join(m.filters, ", ")
	} else {
		filters = "No filters applied"
	}
	if m.width > 0 {
		filters = utils.Substr(filters, m.width-4)
	}

	filtersStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		SetString(filters).
		Width(m.width).
		Align(lipgloss.Right).
		Padding(0, 2)

	return footerStyle.Render(m.input.View() + help.Render() + "\n" + filtersStyle.Render())
}

// Init initializes the model.
func (m FooterModel) Init() tea.Cmd {
	return nil
}

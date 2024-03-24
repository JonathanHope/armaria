package footer

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/textinput"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case msgs.SizeMsg:
		if msg.Name == m.name {
			m.width = msg.Width

			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msgs.SizeMsg{
				Name:  m.inputName,
				Width: m.width - HelpInfoWidth,
			})
			cmds = append(cmds, inputCmd)
		}

	case msgs.InputModeMsg:
		if msg.Name == m.name {
			m.inputMode = msg.InputMode

			if m.inputMode {
				cmds = append(cmds, m.startInputCmd(msg.Prompt, msg.Text, msg.MaxChars))
			} else {
				cmds = append(cmds, m.endInputCmd())
			}
		}

	case msgs.FiltersMsg:
		if msg.Name == m.name {
			m.filters = msg.Filters
		}

	case tea.KeyMsg:
		if m.inputMode {
			switch msg.String() {
			case "ctrl+c":
				if m.inputMode {
					return m, tea.Quit
				}

			case "esc":
				cmds = append(cmds, func() tea.Msg {
					return msgs.InputCancelledMsg{Name: m.name}
				})

			case "enter":
				if m.input.Text() != "" {
					cmds = append(cmds, func() tea.Msg {
						return msgs.InputConfirmedMsg{Name: m.name}
					})
				}

			default:
				var inputCmd tea.Cmd
				m.input, inputCmd = m.input.Update(msg)
				cmds = append(cmds, inputCmd)
			}
		}

	default:
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		cmds = append(cmds, inputCmd)
	}

	return m, tea.Batch(cmds...)
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

// startInputCmd is a command that switches the footer to input mode.
func (m FooterModel) startInputCmd(prompt string, text string, maxChars int) tea.Cmd {
	return tea.Batch(func() tea.Msg {
		return msgs.PromptMsg{Name: m.inputName, Prompt: prompt}
	}, func() tea.Msg {
		return msgs.TextMsg{Name: m.inputName, Text: text}
	}, func() tea.Msg {
		return msgs.FocusMsg{Name: m.inputName, MaxChars: maxChars}
	})
}

// endInputCmd is a command that switches the footer out of input mode.
func (m FooterModel) endInputCmd() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		return msgs.BlurMsg{Name: m.inputName}
	}, func() tea.Msg {
		return msgs.PromptMsg{Name: m.inputName, Prompt: ""}
	}, func() tea.Msg {
		return msgs.TextMsg{Name: m.inputName, Text: ""}
	})
}

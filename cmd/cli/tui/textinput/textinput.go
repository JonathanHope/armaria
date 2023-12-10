package textinput

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
)

const BlinkSpeed = 600 // how quickly to blink the cursor
const Padding = 1      // how much left and right padding to add

// model is the model for a textinput.
// The textinput allows users to enter and modify text.
type model struct {
	name       string  // name of the text input
	prompt     string  // prompt for the input
	text       string  // the current text being inputted
	width      int     // the width of the text input
	cursor     int     // location of the cursor in the text input
	focus      bool    // whether the text input is focused or not
	blink      bool    // flag that alternates in order to make the cursor blink
	frameStart int     // where the viewable frame of text starts
	sleeper    sleeper // used to sleep
}

// InitialModel builds the model.
func InitialModel(name string, prompt string) model {
	return model{
		name:    name,
		prompt:  prompt,
		sleeper: timeSleeper{},
	}
}

// toEnd moves the cursor to the end of the textinput.
func (m *model) toEnd() {
	m.cursor = len(m.text)
	if m.cursor > m.available()-1 {
		diff := m.cursor - (m.available() - 1)
		m.frameStart += diff
		m.cursor = m.available() - 1
	}
}

// Update handles a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case msgs.FocusMsg:
		if m.name == msg.Name {
			m.focus = true
			m.blink = true
			m.cursor = 0
			m.frameStart = 0
			m.toEnd()
			return m, m.blinkCmd()
		}

	case msgs.BlurMsg:
		m.focus = false
		m.blink = false
		m.cursor = 0
		m.frameStart = 0

	case msgs.BlinkMsg:
		if m.focus && msg.Name == m.name {
			m.blink = !m.blink
			return m, m.blinkCmd()
		}

	case msgs.TextMsg:
		if m.name == msg.Name {
			m.text = msg.Text
			m.toEnd()
		}

	case msgs.PromptMsg:
		if m.name == msg.Name {
			m.prompt = msg.Prompt
		}

	case msgs.SizeMsg:
		if m.name == msg.Name {
			m.width = msg.Width - Padding*2
		}

	case tea.KeyMsg:
		if m.focus {
			switch msg.String() {

			case "backspace":
				// No need to delete a char if the text is empty or the cursor isn't behind a char.
				if m.text != "" && m.cursor != 0 {
					text := strings.Split(m.text, "")

					if m.cursor == 1 {
						// Delete a char at the start of the text.
						m.text = strings.Join(text[1:], "")
					} else if m.cursor+m.frameStart == len(m.text) {
						// Delete a char at the end of the text.
						m.text = strings.Join(text[:len(text)-1], "")
					} else {
						// Delete a char in the middle of the text.
						first := strings.Join(text[:m.cursor+m.frameStart-1], "")
						rest := strings.Join(text[m.cursor+m.frameStart:], "")
						m.text = first + rest
					}

					// Move the frame if needed to keep it full.
					// Otherwise move the cursor back a position.
					if m.cursor+m.frameStart > m.available()-1 && m.frameStart > 0 {
						m.frameStart -= 1
					} else {
						m.cursor -= 1
					}

					return m, m.inputChangedCmd()
				}

			case "left":
				if m.cursor > 0 {
					// Move the cursor back if it's not at the start of the frame.
					m.cursor -= 1
				} else if m.cursor == 0 && m.frameStart > 0 {
					// Move the frame back if the cursor is at the start of the frame and it's possible.
					m.frameStart -= 1
				}

			case "right":
				if m.cursor < m.available()-1 && m.cursor < len(m.text) {
					// Move cursor forward if it's not at the end of the frame.
					m.cursor += 1
				} else if m.cursor == m.available()-1 && m.frameStart+m.cursor < len(m.text) {
					// Move the frame forward if the cursor is at the end of the frame and it's possible.
					m.frameStart += 1
				}

			default:
				if m.cursor == 0 {
					// Insert the char at start of the text.
					m.text = string(msg.Runes) + m.text
				} else if m.cursor+m.frameStart == len(m.text) {
					// Insert the char at the end of the text.
					m.text += string(msg.Runes)
				} else {
					// Insert the char in the middle of the text.
					text := strings.Split(m.text, "")
					first := strings.Join(text[:m.cursor+m.frameStart], "")
					rest := strings.Join(text[m.cursor+m.frameStart:], "")
					m.text = first + string(msg.Runes) + rest
				}

				// Move the cursor forward.
				// If the cursor would move past the end of the frame move the frame forward instead.
				m.cursor += len(msg.Runes)
				if m.cursor > m.available()-1 {
					diff := m.cursor - (m.available() - 1)
					m.frameStart += diff
					m.cursor = m.available() - 1
				}

				return m, m.inputChangedCmd()
			}
		}
	}

	return m, nil
}

// available returns the available space for text.
func (m model) available() int {
	return m.width - len(m.prompt)
}

// View renders the model.
func (m model) View() string {
	promptStyle := lipgloss.
		NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("1"))

	cursor := lipgloss.
		NewStyle().
		Inline(true)
	if m.blink {
		cursor = cursor.Reverse(true)
	}

	s := lipgloss.
		NewStyle().
		Width(m.width).
		Padding(0, Padding)

	if m.width == 0 {
		return s.Render(promptStyle.Render(m.prompt))
	}

	available := m.width - len(m.prompt)
	text := strings.Split(m.text+" ", "")
	if len(text) > available {
		text = text[m.frameStart : m.frameStart+available]
	}

	if m.cursor == 0 {
		// Render the view with the cursor at the start.
		under := strings.Join(text[0:1], "")
		rest := strings.Join(text[1:], "")
		return s.Render(promptStyle.Render(m.prompt) + cursor.Render(under) + rest)
	}

	if m.cursor+m.frameStart == len(m.text) {
		// Render the view with the cursor at the end.
		rest := strings.Join(text[0:len(text)-1], "")
		under := strings.Join(text[len(text)-1:], "")
		return s.Render(promptStyle.Render(m.prompt) + rest + cursor.Render(under))
	}

	// Render the view with the cursor in the middle.
	first := strings.Join(text[0:m.cursor], "")
	under := strings.Join(text[m.cursor:m.cursor+1], "")
	rest := strings.Join(text[m.cursor+1:], "")
	return s.Render(promptStyle.Render(m.prompt) + first + cursor.Render(under) + rest)
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return nil
}

// blinkCmd makes the cursor blink.
func (m *model) blinkCmd() tea.Cmd {
	return func() tea.Msg {
		// By sleeping and then returning another BlinkMsg we can make the cursor blink.
		m.sleeper.sleep(BlinkSpeed * time.Millisecond)
		return msgs.BlinkMsg{Name: m.name}
	}
}

// inputChangedCmd publishes a message with the current text.
func (m model) inputChangedCmd() tea.Cmd {
	return func() tea.Msg {
		return msgs.InputChangedMsg{Name: m.name, Text: m.text}
	}
}

type sleeper interface {
	// sleep pauses execution for the requested duration.
	sleep(time.Duration)
}

// timeSleeper implements sleeper with the time package.
type timeSleeper struct{}

func (s timeSleeper) sleep(d time.Duration) {
	time.Sleep(d)
}

package textinput

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/muesli/reflow/ansi"
)

const BlinkSpeed = 600 // how quickly to blink the cursor
const Padding = 1      // how much left and right padding to add

// TextInputModel is the TextInputModel for a textinput.
// The textinput allows users to enter and modify text.
type TextInputModel struct {
	name        string  // name of the text input
	prompt      string  // prompt for the input
	text        string  // the current text being inputted
	width       int     // the width of the text input
	cursor      int     // location of the cursor in the window
	index       int     // which character is selected in the text
	focus       bool    // whether the text input is focused or not
	maxChars    int     // maximum number of chars to allow
	blink       bool    // flag that alternates in order to make the cursor blink
	windowStart int     // index the window starts at
	windowEnd   int     // index the window ends at
	sleeper     sleeper // used to sleep
}

// InitialModel builds the model.
func InitialModel(name string, prompt string) TextInputModel {
	return TextInputModel{
		name:    name,
		prompt:  prompt,
		sleeper: timeSleeper{},
	}
}

// Text returns the current text in the input.
func (m *TextInputModel) Text() string {
	return m.text
}

// textWithSpace returns the current text with a space at the end.
// This input uses a block cursor so the extra space is needed.
func (m *TextInputModel) textWithSpace() string {
	return m.text + " "
}

// windowWidth returns the available width for text.
// The avaialble width is the overall width less the measured prompt width.
func (m TextInputModel) windowWidth() int {
	width := m.width - ansi.PrintableRuneWidth(m.prompt) - Padding*2
	if width < 0 {
		width = 0
	}

	return width
}

// window returns the currently visible part of the text.
func (m *TextInputModel) window() string {
	if m.textWithSpace() == " " || m.width == 0 {
		return " "
	}

	textRunes := strings.Split(m.textWithSpace(), "")
	start := m.windowStart
	end := m.windowEnd
	if end > len(textRunes) {
		end = len(textRunes)
	}

	windowRunes := textRunes[start:end]
	return strings.Join(windowRunes, "")
}

// initWindow intializes the window after the text changes.
func (m *TextInputModel) initWindow() {
	m.windowStart = 0
	m.windowEnd = m.windowWidth()
	m.cursor = 0
	m.index = 0
	m.moveEnd()
	m.chopLeft()
}

// cursorAtStart returns true if the cursor is at the start of the window.
func (m *TextInputModel) cursorAtStart() bool {
	return m.cursor == 0
}

// cursorAtEnd returns true if the cursor is at the end of the window.
func (m *TextInputModel) cursorAtEnd() bool {
	windowRunes := strings.Split(m.window(), "")
	return m.cursor == len(windowRunes)-1
}

// indexAtStart returns true if the cursor is at the start of the text.
func (m *TextInputModel) indexAtStart() bool {
	return m.index == 0
}

// indexAtEnd returns true if the cursor is at the end of the text.
func (m *TextInputModel) indexAtEnd() bool {
	textRunes := strings.Split(m.textWithSpace(), "")
	return m.index == len(textRunes)-1
}

// chopRight chops runes off the right of the window to make it fit.
func (m *TextInputModel) chopRight() {
	textRunes := strings.Split(m.textWithSpace(), "")
	for ansi.PrintableRuneWidth(m.window()) < m.windowWidth() && m.windowEnd < len(textRunes) {
		m.windowEnd += 1
	}

	for ansi.PrintableRuneWidth(m.window()) > m.windowWidth() {
		m.windowEnd -= 1
	}
}

// chopLeft chops runes off the left of the window to make it fit.
func (m *TextInputModel) chopLeft() {
	for ansi.PrintableRuneWidth(m.window()) < m.windowWidth() && m.windowStart > 0 {
		m.windowStart -= 1
	}

	for ansi.PrintableRuneWidth(m.window()) > m.windowWidth() {
		m.windowStart += 1
	}
}

// moveRight moves the cursor to right once space.
func (m *TextInputModel) moveRight() {
	shift := !m.indexAtEnd() && m.cursorAtEnd()
	previousLength := m.windowEnd - m.windowStart

	if !m.indexAtEnd() {
		m.index += 1
	}

	if !m.cursorAtEnd() {
		m.cursor += 1
	}

	if shift {
		m.windowEnd += 1
	}

	m.chopLeft()

	newLength := m.windowEnd - m.windowStart
	m.cursor += newLength - previousLength
}

// moveLeft moves the cursor the left once space.
func (m *TextInputModel) moveLeft() {
	shift := !m.indexAtStart() && m.cursorAtStart()

	if !m.indexAtStart() {
		m.index -= 1
	}

	if !m.cursorAtStart() {
		m.cursor -= 1
	}

	if shift {
		m.windowStart -= 1
	}

	m.chopRight()
}

// moveEnd moves to the end of the text.
func (m *TextInputModel) moveEnd() {
	for !m.cursorAtEnd() || !m.indexAtEnd() {
		m.moveRight()
	}
}

// moveEnd moves to the start of the text.
func (m *TextInputModel) moveStart() {
	for !m.cursorAtStart() || !m.indexAtStart() {
		m.moveLeft()
	}
}

// delete deletes the rune in front of the cursor.
func (m *TextInputModel) delete() {
	if m.text == "" || m.index == 0 {
		return
	}

	textRunes := strings.Split(m.text, "")

	if m.index == 1 { // Delete a char at the start of the text.
		m.text = strings.Join(textRunes[1:], "")
	} else if m.indexAtEnd() { // Delete a char at the end of the text.
		m.text = strings.Join(textRunes[:len(textRunes)-1], "")
	} else { // Delete a char in the middle of the text.
		first := strings.Join(textRunes[:m.index-1], "")
		rest := strings.Join(textRunes[m.index:], "")
		m.text = first + rest
	}

	m.index -= 1

	if m.windowStart > 0 {
		m.windowStart -= 1
		m.chopRight()
	} else {
		m.cursor -= 1
	}
}

// insert inserts runes in front of the cursor.
func (m *TextInputModel) insert(runes []rune) {
	textRunes := strings.Split(m.text, "")
	cursorAtEnd := m.cursorAtEnd()

	if m.maxChars > 0 && len(textRunes)+len(runes) > m.maxChars {
		return
	}

	if m.indexAtStart() { // Insert the char at start of the text.
		m.text = string(runes) + m.text
	} else if m.indexAtEnd() { // Insert the char at the end of the text.
		m.text += string(runes)
	} else { // Insert the char in the middle of the text.
		first := strings.Join(textRunes[:m.index], "")
		rest := strings.Join(textRunes[m.index:], "")
		m.text = first + string(runes) + rest
	}

	if cursorAtEnd {
		m.moveEnd()
	} else {
		m.cursor += len(runes)
		m.index += len(runes)
		m.chopRight()
	}
}

// Update handles a message.
func (m TextInputModel) Update(msg tea.Msg) (TextInputModel, tea.Cmd) {
	switch msg := msg.(type) {

	case msgs.FocusMsg:
		if m.name == msg.Name {
			m.focus = true
			m.blink = true
			m.maxChars = msg.MaxChars
			m.initWindow()
			return m, m.blinkCmd()
		}

	case msgs.BlurMsg:
		m.focus = false
		m.blink = false
		m.maxChars = 0
		m.initWindow()

	case msgs.BlinkMsg:
		if m.focus && msg.Name == m.name {
			m.blink = !m.blink
			return m, m.blinkCmd()
		}

	case msgs.TextMsg:
		if m.name == msg.Name {
			m.text = msg.Text
			m.initWindow()
		}

	case msgs.PromptMsg:
		if m.name == msg.Name {
			m.prompt = msg.Prompt
			m.initWindow()
		}

	case msgs.SizeMsg:
		if m.name == msg.Name {
			m.width = msg.Width - Padding*2
			m.initWindow()
		}

	case tea.KeyMsg:
		if m.focus {
			switch msg.String() {

			case "backspace":
				m.delete()
				return m, m.inputChangedCmd()

			case "left":
				m.moveLeft()

			case "right":
				m.moveRight()

			default:
				m.insert(msg.Runes)
				return m, m.inputChangedCmd()
			}
		}
	}

	return m, nil
}

// View renders the model.
func (m TextInputModel) View() string {
	promptStyle := lipgloss.
		NewStyle().
		Bold(true).
		Inline(true).
		Foreground(lipgloss.Color("1"))

	cursorStyle := lipgloss.
		NewStyle().
		Inline(true)
	if m.blink {
		cursorStyle = cursorStyle.Reverse(true)
	}

	s := lipgloss.
		NewStyle().
		Width(m.width).
		Padding(0, Padding)

	window := m.window()
	windowRunes := strings.Split(window, "")
	if window == "" {
		return s.Render(promptStyle.Render(m.prompt))
	}

	if m.cursorAtStart() { // Render the view with the cursor at the start.
		under := strings.Join(windowRunes[0:1], "")
		rest := strings.Join(windowRunes[1:], "")
		return s.Render(promptStyle.Render(m.prompt) + cursorStyle.Render(under) + rest)
	}

	if m.cursorAtEnd() { // Render the view with the cursor at the end.
		rest := strings.Join(windowRunes[0:len(windowRunes)-1], "")
		under := strings.Join(windowRunes[len(windowRunes)-1:], "")
		return s.Render(promptStyle.Render(m.prompt) + rest + cursorStyle.Render(under))
	}

	// Render the view with the cursor in the middle.
	first := strings.Join(windowRunes[0:m.cursor], "")
	under := strings.Join(windowRunes[m.cursor:m.cursor+1], "")
	rest := strings.Join(windowRunes[m.cursor+1:], "")
	return s.Render(promptStyle.Render(m.prompt) + first + cursorStyle.Render(under) + rest)

}

// Init initializes the model.
func (m TextInputModel) Init() tea.Cmd {
	return nil
}

// blinkCmd makes the cursor blink.
func (m *TextInputModel) blinkCmd() tea.Cmd {
	return func() tea.Msg {
		// By sleeping and then returning another BlinkMsg we can make the cursor blink.
		m.sleeper.sleep(BlinkSpeed * time.Millisecond)
		return msgs.BlinkMsg{Name: m.name}
	}
}

// inputChangedCmd publishes a message with the current text.
func (m TextInputModel) inputChangedCmd() tea.Cmd {
	return func() tea.Msg {
		return msgs.InputChangedMsg{Name: m.name}
	}
}

type sleeper interface {
	// sleep pauses execution for the requested duration.
	sleep(time.Duration)
}

// timeSleeper implements sleeper with the time package.
type timeSleeper struct{}

// sleep pauses execution of the calling thread for the requested duration.
func (s timeSleeper) sleep(d time.Duration) {
	time.Sleep(d)
}

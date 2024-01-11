package textinput

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

const name = "TestInput"
const prompt = ">"
const width = 4

func TestTextScrollsAsInputAdded(t *testing.T) {
	// enter 1

	gotModel := TextInputModel{
		name:   name,
		prompt: prompt,
		width:  width,
		cursor: 0,
		focus:  true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

	wantModel := TextInputModel{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "1",
		cursor: 1,
		focus:  true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// enter 2

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})

	wantModel = TextInputModel{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "12",
		cursor: 2,
		focus:  true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// enter 3

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})

	wantModel = TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 1,
		focus:      true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestTextScrollsAsInputRemoved(t *testing.T) {
	// delete 3

	gotModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 1,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel := TextInputModel{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "12",
		cursor: 2,
		focus:  true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// delete 2

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel = TextInputModel{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "1",
		cursor: 1,
		focus:  true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// delete 1

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel = TextInputModel{
		name:   name,
		prompt: prompt,
		width:  width,
		cursor: 0,
		focus:  true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// noop

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanScrollToStartOfText(t *testing.T) {
	// move to 3

	gotModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 1,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyLeft})

	wantModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     1,
		frameStart: 1,
		focus:      true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	// move to 2

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyLeft})

	wantModel = TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 1,
		focus:      true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	// move to 1

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyLeft})

	wantModel = TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 0,
		focus:      true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	// noop

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyLeft})

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanScrollToEndOfText(t *testing.T) {
	// move to 2

	gotModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRight})

	wantModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	// move to 3

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyRight})

	wantModel = TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 0,
		focus:      true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	// move to end

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyRight})

	wantModel = TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 1,
		focus:      true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	// noop

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyRight})

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanInsertAtStart(t *testing.T) {
	gotModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}})

	wantModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "0123",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanDeleteAtStart(t *testing.T) {
	gotModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "0123",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanInsertInMiddle(t *testing.T) {
	gotModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "13",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})

	wantModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanDeleteInMiddle(t *testing.T) {
	gotModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel := TextInputModel{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "13",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanChangePrompt(t *testing.T) {
	gotModel := TextInputModel{
		name: name,
	}

	gotModel, gotCmd := gotModel.Update(msgs.PromptMsg{Name: name, Prompt: prompt})

	wantModel := TextInputModel{
		name:   name,
		prompt: prompt,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanChangeText(t *testing.T) {
	const text = "123"

	gotModel := TextInputModel{
		name:  name,
		width: 4,
	}

	gotModel, gotCmd := gotModel.Update(msgs.TextMsg{Name: name, Text: text})

	wantModel := TextInputModel{
		name:   name,
		text:   text,
		width:  4,
		cursor: 3,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanChangeSize(t *testing.T) {
	const width = 3

	gotModel := TextInputModel{
		name: name,
	}

	gotModel, gotCmd := gotModel.Update(msgs.SizeMsg{Name: name, Width: width})

	wantModel := TextInputModel{
		name:  name,
		width: width - Padding*2,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanFocus(t *testing.T) {
	gotModel := TextInputModel{
		name:    name,
		focus:   false,
		blink:   false,
		width:   1,
		sleeper: noopSleeper{},
	}

	gotModel, gotCmd := gotModel.Update(msgs.FocusMsg{Name: name})

	wantModel := TextInputModel{
		name:    name,
		focus:   true,
		blink:   true,
		width:   1,
		sleeper: noopSleeper{},
	}

	wantCmd := func() tea.Msg { return msgs.BlinkMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanBlur(t *testing.T) {
	gotModel := TextInputModel{
		name:    name,
		focus:   true,
		blink:   true,
		sleeper: noopSleeper{},
	}

	gotModel, gotCmd := gotModel.Update(msgs.BlurMsg{Name: name})

	wantModel := TextInputModel{
		name:    name,
		focus:   false,
		blink:   false,
		sleeper: noopSleeper{},
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanBlink(t *testing.T) {
	gotModel := TextInputModel{
		name:    name,
		focus:   true,
		blink:   true,
		sleeper: noopSleeper{},
	}

	gotModel, gotCmd := gotModel.Update(msgs.BlinkMsg{Name: name})

	wantModel := TextInputModel{
		name:    name,
		focus:   true,
		blink:   false,
		sleeper: noopSleeper{},
	}

	wantCmd := func() tea.Msg { return msgs.BlinkMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanLimitMaxChars(t *testing.T) {
	gotModel := TextInputModel{
		name:    name,
		sleeper: noopSleeper{},
		width:   2,
	}

	gotModel, _ = gotModel.Update(msgs.FocusMsg{Name: name, MaxChars: 1})
	gotModel, _ = gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	gotModel, _ = gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

	wantModel := TextInputModel{
		name:     name,
		sleeper:  noopSleeper{},
		width:    2,
		maxChars: 1,
		text:     "1",
		cursor:   1,
		focus:    true,
		blink:    true,
	}

	verifyUpdate(t, gotModel, wantModel, nil, nil)
}

func TestText(t *testing.T) {
	gotModel := TextInputModel{
		text: "123",
	}

	diff := cmp.Diff(gotModel.Text(), "123")
	if diff != "" {
		t.Errorf("Expected and actual text different:\n%s", diff)
	}
}

func verifyUpdate(t *testing.T, gotModel TextInputModel, wantModel TextInputModel, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(TextInputModel{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	utils.CompareCommands(t, gotCmd, wantCmd)
}

type noopSleeper struct{}

func (s noopSleeper) sleep(d time.Duration) {
}

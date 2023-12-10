package textinput

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
)

const name = "TestInput"
const prompt = ">"
const width = 4

func TestTextScrollsAsInputAdded(t *testing.T) {
	// enter 1

	var gotModel tea.Model = model{
		name:   name,
		prompt: prompt,
		width:  width,
		cursor: 0,
		focus:  true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

	wantModel := model{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "1",
		cursor: 1,
		focus:  true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "1"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// enter 2

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})

	wantModel = model{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "12",
		cursor: 2,
		focus:  true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "12"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// enter 3

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})

	wantModel = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 1,
		focus:      true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "123"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestTextScrollsAsInputRemoved(t *testing.T) {
	// delete 3

	var gotModel tea.Model = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 1,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel := model{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "12",
		cursor: 2,
		focus:  true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "12"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// delete 2

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel = model{
		name:   name,
		prompt: prompt,
		width:  width,
		text:   "1",
		cursor: 1,
		focus:  true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "1"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// delete 1

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel = model{
		name:   name,
		prompt: prompt,
		width:  width,
		cursor: 0,
		focus:  true,
	}

	wantCmd = func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: ""} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	// noop

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanScrollToStartOfText(t *testing.T) {
	// move to 3

	var gotModel tea.Model = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 1,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyLeft})

	wantModel := model{
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

	wantModel = model{
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

	wantModel = model{
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

	var gotModel tea.Model = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRight})

	wantModel := model{
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

	wantModel = model{
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

	wantModel = model{
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
	var gotModel tea.Model = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}})

	wantModel := model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "0123",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "0123"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanDeleteAtStart(t *testing.T) {
	var gotModel tea.Model = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "0123",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel := model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     0,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "123"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanInsertInMiddle(t *testing.T) {
	var gotModel tea.Model = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "13",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})

	wantModel := model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "123"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanDeleteInMiddle(t *testing.T) {
	var gotModel tea.Model = model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "123",
		cursor:     2,
		frameStart: 0,
		focus:      true,
	}

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	wantModel := model{
		name:       name,
		prompt:     prompt,
		width:      width,
		text:       "13",
		cursor:     1,
		frameStart: 0,
		focus:      true,
	}

	wantCmd := func() tea.Msg { return msgs.InputChangedMsg{Name: name, Text: "13"} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanChangePrompt(t *testing.T) {
	var gotModel tea.Model = model{
		name: name,
	}

	gotModel, gotCmd := gotModel.Update(msgs.PromptMsg{Name: name, Prompt: prompt})

	wantModel := model{
		name:   name,
		prompt: prompt,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanChangeText(t *testing.T) {
	const text = "123"

	var gotModel tea.Model = model{
		name:  name,
		width: 4,
	}

	gotModel, gotCmd := gotModel.Update(msgs.TextMsg{Name: name, Text: text})

	wantModel := model{
		name:   name,
		text:   text,
		width:  4,
		cursor: 3,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanChangeSize(t *testing.T) {
	const width = 3

	var gotModel tea.Model = model{
		name: name,
	}

	gotModel, gotCmd := gotModel.Update(msgs.SizeMsg{Name: name, Width: width})

	wantModel := model{
		name:  name,
		width: width - Padding*2,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanFocus(t *testing.T) {
	var gotModel tea.Model = model{
		name:    name,
		focus:   false,
		blink:   false,
		width:   1,
		sleeper: noopSleeper{},
	}

	gotModel, gotCmd := gotModel.Update(msgs.FocusMsg{Name: name})

	wantModel := model{
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
	var gotModel tea.Model = model{
		name:    name,
		focus:   true,
		blink:   true,
		sleeper: noopSleeper{},
	}

	gotModel, gotCmd := gotModel.Update(msgs.BlurMsg{Name: name})

	wantModel := model{
		name:    name,
		focus:   false,
		blink:   false,
		sleeper: noopSleeper{},
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanBlink(t *testing.T) {
	var gotModel tea.Model = model{
		name:    name,
		focus:   true,
		blink:   true,
		sleeper: noopSleeper{},
	}

	gotModel, gotCmd := gotModel.Update(msgs.BlinkMsg{Name: name})

	wantModel := model{
		name:    name,
		focus:   true,
		blink:   false,
		sleeper: noopSleeper{},
	}

	wantCmd := func() tea.Msg { return msgs.BlinkMsg{Name: name} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func verifyUpdate(t *testing.T, gotModel tea.Model, wantModel tea.Model, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(model{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	if gotCmd == nil || wantCmd == nil {
		if gotCmd != nil || wantCmd != nil {
			t.Errorf("Expected and actual cmds different: one is nil and one is non-nil")
		}

		return
	}

	cmdDiff := cmp.Diff(gotCmd(), wantCmd())
	if modelDiff != "" {
		t.Errorf("Expected and actual cmds different:\n%s", cmdDiff)
	}
}

type noopSleeper struct{}

func (s noopSleeper) sleep(d time.Duration) {
}

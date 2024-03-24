package textinput

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"testing"
)

const name = "TestInput"
const prompt = ">"

func TestText(t *testing.T) {
	model := TextInputModel{
		text: "test",
	}

	diff := cmp.Diff(model.Text(), "test")
	if diff != "" {
		t.Errorf("Expected and actual texts different")
	}
}

func TestInsert(t *testing.T) {
	validate := func(model TextInputModel, window string) {
		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual windows different:\n%s", windowDiff)
		}
	}

	// start

	model := TextInputModel{
		name:  name,
		text:  "🐂🐜",
		width: 12,
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'🦊'}})

	validate(model, "🦊🐂🐜 ")

	// end

	model = TextInputModel{
		name:  name,
		text:  "🦊🐂",
		width: 12,
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'🐜'}})

	validate(model, "🦊🐂🐜 ")

	// middle

	model = TextInputModel{
		name:  name,
		text:  "🦊🐜",
		width: 12,
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'🐂'}})

	validate(model, "🦊🐂🐜 ")
}

func TestInsertMovesWindow(t *testing.T) {
	validate := func(model TextInputModel, index int, window string) {
		atEndDiff := cmp.Diff(model.cursorAtEnd(), true)
		if atEndDiff != "" {
			t.Errorf("Expected and actual cursor at ends different:\n%s", atEndDiff)
		}

		cursorDiff := cmp.Diff(model.cursor, 1)
		if cursorDiff != "" {
			t.Errorf("Expected and actual cursors different:\n%s", cursorDiff)
		}

		indexDiff := cmp.Diff(model.index, index)
		if indexDiff != "" {
			t.Errorf("Expected and actual indexes different:\n%s", indexDiff)
		}

		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual windows different:\n%s", windowDiff)
		}
	}

	model := TextInputModel{
		name:   name,
		width:  6,
		prompt: prompt,
		text:   "🦊",
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	validate(model, 1, "🦊 ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	validate(model, 2, "c ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'🦊'}})
	validate(model, 3, "🦊 ")
}

func TestDelete(t *testing.T) {
	validate := func(model TextInputModel, window string) {
		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual windows different:\n%s", windowDiff)
		}
	}

	// start

	model := TextInputModel{
		name:  name,
		text:  "🦊🐂🐜",
		width: 12,
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	validate(model, "🐂🐜 ")

	// end

	model = TextInputModel{
		name:  name,
		text:  "🦊🐂🐜",
		width: 12,
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	validate(model, "🦊🐂 ")

	// middle

	model = TextInputModel{
		name:  name,
		text:  "🦊🐂🐜",
		width: 12,
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	validate(model, "🦊🐜 ")
}

func TestDeleteMovesWindow(t *testing.T) {
	validate := func(model TextInputModel, window string) {
		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual windows different:\n%s", windowDiff)
		}
	}

	model := TextInputModel{
		name:  name,
		text:  "🦊c🦊c🦊c",
		width: 6,
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	validate(model, "🦊c ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	validate(model, "c🦊 ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	validate(model, "🦊c ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	validate(model, "c🦊 ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	validate(model, "🦊c ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	validate(model, "🦊 ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	validate(model, " ")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	validate(model, " ")
}

func TestMoveRight(t *testing.T) {
	validate := func(model TextInputModel, window string, cursor int, index int) {
		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual windows different:\n%s", windowDiff)
		}

		cursorDiff := cmp.Diff(model.cursor, cursor)
		if cursorDiff != "" {
			t.Errorf("Expected and actual cursors different:\n%s", cursorDiff)
		}

		indexDiff := cmp.Diff(model.index, index)
		if indexDiff != "" {
			t.Errorf("Expected and actual indexes different:\n%s", indexDiff)
		}
	}

	model := TextInputModel{
		name:   name,
		width:  6,
		prompt: prompt,
		text:   "a🦊b🐂c🐜",
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "a🦊", 0, 0)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "a🦊", 1, 1)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🦊b", 1, 2)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "b🐂", 1, 3)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐂c", 1, 4)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "c🐜", 1, 5)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐜 ", 1, 6)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐜 ", 1, 6)
}

func TestMoveRightVariation2(t *testing.T) {
	validate := func(model TextInputModel, window string, cursor int, index int) {
		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual at windows different:\n%s", windowDiff)
		}

		cursorDiff := cmp.Diff(model.cursor, cursor)
		if cursorDiff != "" {
			t.Errorf("Expected and actual at cursors different:\n%s", cursorDiff)
		}

		indexDiff := cmp.Diff(model.index, index)
		if indexDiff != "" {
			t.Errorf("Expected and actual at indexes different:\n%s", indexDiff)
		}
	}

	model := TextInputModel{
		name:   name,
		width:  7,
		prompt: prompt,
		text:   "🦊🐂abcd🐜🐕",
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🦊🐂", 0, 0)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🦊🐂", 1, 1)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐂a", 1, 2)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐂ab", 2, 3)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "abc", 2, 4)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "abcd", 3, 5)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "cd🐜", 2, 6)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐜🐕", 1, 7)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐕 ", 1, 8)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	validate(model, "🐕 ", 1, 8)
}

func TestMoveLeft(t *testing.T) {
	validate := func(model TextInputModel, window string, cursor int, index int) {
		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual windows different:\n%s", windowDiff)
		}

		cursorDiff := cmp.Diff(model.cursor, cursor)
		if cursorDiff != "" {
			t.Errorf("Expected and actual cursors different:\n%s", cursorDiff)
		}

		indexDiff := cmp.Diff(model.index, index)
		if indexDiff != "" {
			t.Errorf("Expected and actual indexes different:\n%s", indexDiff)
		}
	}

	model := TextInputModel{
		name:   name,
		width:  6,
		prompt: prompt,
		text:   "a🦊b🐂c🐜",
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	validate(model, "🐜 ", 1, 6)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🐜 ", 0, 5)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "c🐜", 0, 4)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🐂c", 0, 3)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "b🐂", 0, 2)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🦊b", 0, 1)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "a🦊", 0, 0)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "a🦊", 0, 0)
}

func TestMoveLeftVariation2(t *testing.T) {
	validate := func(model TextInputModel, window string, cursor int, index int) {
		windowDiff := cmp.Diff(model.window(), window)
		if windowDiff != "" {
			t.Errorf("Expected and actual windows different:\n%s", windowDiff)
		}

		cursorDiff := cmp.Diff(model.cursor, cursor)
		if cursorDiff != "" {
			t.Errorf("Expected and actual cursors different:\n%s", cursorDiff)
		}

		indexDiff := cmp.Diff(model.index, index)
		if indexDiff != "" {
			t.Errorf("Expected and actual indexes different:\n%s", indexDiff)
		}
	}

	model := TextInputModel{
		name:   name,
		width:  7,
		prompt: prompt,
		text:   "🦊🐂abcd🐜🐕",
	}

	model, _ = model.Update(msgs.FocusMsg{Name: name})
	validate(model, "🐕 ", 1, 8)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🐕 ", 0, 7)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🐜🐕", 0, 6)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "d🐜", 0, 5)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "cd🐜", 0, 4)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "bcd", 0, 3)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "abcd", 0, 2)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🐂ab", 0, 1)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🦊🐂", 0, 0)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	validate(model, "🦊🐂", 0, 0)
}

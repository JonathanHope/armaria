package textinput

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

	model.Focus(0)
	model.MoveLeft()
	model.MoveLeft()
	model.Insert([]rune{'🦊'})

	validate(model, "🦊🐂🐜 ")

	// end

	model = TextInputModel{
		name:  name,
		text:  "🦊🐂",
		width: 12,
	}

	model.Focus(0)
	model.Insert([]rune{'🐜'})

	validate(model, "🦊🐂🐜 ")

	// middle

	model = TextInputModel{
		name:  name,
		text:  "🦊🐜",
		width: 12,
	}

	model.Focus(0)
	model.MoveLeft()
	model.Insert([]rune{'🐂'})

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

	model.Focus(0)
	validate(model, 1, "🦊 ")

	model.Insert([]rune{'c'})
	validate(model, 2, "c ")

	model.Insert([]rune{'🦊'})
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

	model.Focus(0)
	model.MoveLeft()
	model.MoveLeft()
	model.Delete()

	validate(model, "🐂🐜 ")

	// end

	model = TextInputModel{
		name:  name,
		text:  "🦊🐂🐜",
		width: 12,
	}

	model.Focus(0)
	model.Delete()

	validate(model, "🦊🐂 ")

	// middle

	model = TextInputModel{
		name:  name,
		text:  "🦊🐂🐜",
		width: 12,
	}

	model.Focus(0)
	model.MoveLeft()
	model.Delete()

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

	model.Focus(0)
	validate(model, "🦊c ")

	model.Delete()
	validate(model, "c🦊 ")

	model.Delete()
	validate(model, "🦊c ")

	model.Delete()
	validate(model, "c🦊 ")

	model.Delete()
	validate(model, "🦊c ")

	model.Delete()
	validate(model, "🦊 ")

	model.Delete()
	validate(model, " ")

	model.Delete()
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

	model.Focus(0)
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	validate(model, "a🦊", 0, 0)

	model.MoveRight()
	validate(model, "a🦊", 1, 1)

	model.MoveRight()
	validate(model, "🦊b", 1, 2)

	model.MoveRight()
	validate(model, "b🐂", 1, 3)

	model.MoveRight()
	validate(model, "🐂c", 1, 4)

	model.MoveRight()
	validate(model, "c🐜", 1, 5)

	model.MoveRight()
	validate(model, "🐜 ", 1, 6)

	model.MoveRight()
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

	model.Focus(0)
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	model.MoveLeft()
	validate(model, "🦊🐂", 0, 0)

	model.MoveRight()
	validate(model, "🦊🐂", 1, 1)

	model.MoveRight()
	validate(model, "🐂a", 1, 2)

	model.MoveRight()
	validate(model, "🐂ab", 2, 3)

	model.MoveRight()
	validate(model, "abc", 2, 4)

	model.MoveRight()
	validate(model, "abcd", 3, 5)

	model.MoveRight()
	validate(model, "cd🐜", 2, 6)

	model.MoveRight()
	validate(model, "🐜🐕", 1, 7)

	model.MoveRight()
	validate(model, "🐕 ", 1, 8)

	model.MoveRight()
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

	model.Focus(0)
	validate(model, "🐜 ", 1, 6)

	model.MoveLeft()
	validate(model, "🐜 ", 0, 5)

	model.MoveLeft()
	validate(model, "c🐜", 0, 4)

	model.MoveLeft()
	validate(model, "🐂c", 0, 3)

	model.MoveLeft()
	validate(model, "b🐂", 0, 2)

	model.MoveLeft()
	validate(model, "🦊b", 0, 1)

	model.MoveLeft()
	validate(model, "a🦊", 0, 0)

	model.MoveLeft()
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

	model.Focus(0)
	validate(model, "🐕 ", 1, 8)

	model.MoveLeft()
	validate(model, "🐕 ", 0, 7)

	model.MoveLeft()
	validate(model, "🐜🐕", 0, 6)

	model.MoveLeft()
	validate(model, "d🐜", 0, 5)

	model.MoveLeft()
	validate(model, "cd🐜", 0, 4)

	model.MoveLeft()
	validate(model, "bcd", 0, 3)

	model.MoveLeft()
	validate(model, "abcd", 0, 2)

	model.MoveLeft()
	validate(model, "🐂ab", 0, 1)

	model.MoveLeft()
	validate(model, "🦊🐂", 0, 0)

	model.MoveLeft()
	validate(model, "🦊🐂", 0, 0)
}

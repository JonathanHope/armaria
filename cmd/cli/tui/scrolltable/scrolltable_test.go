package scrolltable

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
)

func TestCanUpdateData(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	var gotModel tea.Model = model[TestDatum]{
		cursor: 0,
		height: height,
	}
	gotModel, gotCmd := gotModel.Update(msgs.DataMsg[TestDatum]{Data: data})

	wantModel := model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}

	wantCmd := func() tea.Msg { return msgs.SelectionChangedMsg[TestDatum]{} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanScrollDown(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	var gotModel tea.Model = model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	wantModel := model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 1,
	}

	wantCmd := func() tea.Msg { return msgs.SelectionChangedMsg[TestDatum]{} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanScrollUp(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	var gotModel tea.Model = model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 1,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	wantModel := model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg { return msgs.SelectionChangedMsg[TestDatum]{} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanMoveDown(t *testing.T) {
	const height = Reserved + 2
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	var gotModel tea.Model = model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	wantModel := model[TestDatum]{
		cursor:     1,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg { return msgs.SelectionChangedMsg[TestDatum]{} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanMoveUp(t *testing.T) {
	const height = Reserved + 2
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	var gotModel tea.Model = model[TestDatum]{
		cursor:     1,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	wantModel := model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg { return msgs.SelectionChangedMsg[TestDatum]{} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanScrollIfFrameEmpty(t *testing.T) {
	const height = Reserved
	data := []TestDatum{}

	var gotModel tea.Model = model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	wantModel := model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func FrameSizeChangesWithHeight(t *testing.T) {
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	var gotModel tea.Model = model[TestDatum]{
		cursor:     0,
		height:     Reserved + 1,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(msgs.SizeMsg{Height: Reserved + 2})

	wantModel := model[TestDatum]{
		cursor:     0,
		height:     Reserved + 1,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg { return msgs.SelectionChangedMsg[TestDatum]{} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func FrameCannotBeLargerThanData(t *testing.T) {
	const height = Reserved + 2
	data := []TestDatum{
		{ID: "1"},
	}

	var gotModel tea.Model = model[TestDatum]{
		cursor: 0,
		height: height,
	}
	gotModel, gotCmd := gotModel.Update(msgs.DataMsg[TestDatum]{Data: data})

	wantModel := model[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg { return msgs.SelectionChangedMsg[TestDatum]{} }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

type TestDatum struct {
	ID string
}

func verifyUpdate(t *testing.T, gotModel tea.Model, wantModel tea.Model, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(model[TestDatum]{})
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

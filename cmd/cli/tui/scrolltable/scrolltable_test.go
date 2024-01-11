package scrolltable

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

func TestCanUpdateData(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	gotModel := ScrolltableModel[TestDatum]{
		cursor: 0,
		height: height,
	}
	gotModel, gotCmd := gotModel.Update(msgs.DataMsg[TestDatum]{Data: data})

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg {
				return msgs.SelectionChangedMsg[TestDatum]{
					Selection: TestDatum{ID: "1"},
				}
			},
			func() tea.Msg { return msgs.FreeMsg{} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanUpdateDataMoveUp(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		frameStart: 1,
	}
	gotModel, gotCmd := gotModel.Update(msgs.DataMsg[TestDatum]{Data: data, Move: msgs.DirectionUp})

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg {
				return msgs.SelectionChangedMsg[TestDatum]{
					Selection: TestDatum{ID: "1"},
				}
			},
			func() tea.Msg { return msgs.FreeMsg{} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanUpdateDataMoveDown(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	gotModel := ScrolltableModel[TestDatum]{
		cursor: 0,
		height: height,
	}
	gotModel, gotCmd := gotModel.Update(msgs.DataMsg[TestDatum]{Data: data, Move: msgs.DirectionDown})

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 1,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg {
				return msgs.SelectionChangedMsg[TestDatum]{
					Selection: TestDatum{ID: "2"},
				}
			},
			func() tea.Msg { return msgs.FreeMsg{} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanUpdateDataMoveStart(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		frameStart: 1,
	}
	gotModel, gotCmd := gotModel.Update(msgs.DataMsg[TestDatum]{Data: data, Move: msgs.DirectionStart})

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg {
				return msgs.SelectionChangedMsg[TestDatum]{
					Selection: TestDatum{ID: "1"},
				}
			},
			func() tea.Msg { return msgs.FreeMsg{} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanScrollDown(t *testing.T) {
	const height = Reserved + 1
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 1,
	}

	wantCmd := func() tea.Msg {
		return msgs.SelectionChangedMsg[TestDatum]{Selection: TestDatum{ID: "2"}}
	}

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

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 1,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg {
		return msgs.SelectionChangedMsg[TestDatum]{Selection: TestDatum{ID: "1"}}
	}

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

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     1,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg {
		return msgs.SelectionChangedMsg[TestDatum]{Selection: TestDatum{ID: "2"}}
	}

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

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     1,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg {
		return msgs.SelectionChangedMsg[TestDatum]{Selection: TestDatum{ID: "1"}}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanScrollIfFrameEmpty(t *testing.T) {
	const height = Reserved
	data := []TestDatum{}

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyDown}))

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyUp}))

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestFrameSizeChangesWithHeight(t *testing.T) {
	data := []TestDatum{
		{ID: "1"},
		{ID: "2"},
	}

	gotModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     Reserved + 1,
		data:       data,
		frameStart: 0,
	}
	gotModel, gotCmd := gotModel.Update(msgs.SizeMsg{Height: Reserved + 2})

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     Reserved + 2,
		data:       data,
		frameStart: 0,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestFrameCannotBeLargerThanData(t *testing.T) {
	const height = Reserved + 2
	data := []TestDatum{
		{ID: "1"},
	}

	gotModel := ScrolltableModel[TestDatum]{
		cursor: 0,
		height: height,
	}
	gotModel, gotCmd := gotModel.Update(msgs.DataMsg[TestDatum]{Data: data})

	wantModel := ScrolltableModel[TestDatum]{
		cursor:     0,
		height:     height,
		data:       data,
		frameStart: 0,
	}
	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg {
				return msgs.SelectionChangedMsg[TestDatum]{
					Selection: TestDatum{ID: "1"},
				}
			},
			func() tea.Msg { return msgs.FreeMsg{} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestEmpty(t *testing.T) {
	gotModel := ScrolltableModel[TestDatum]{}

	diff := cmp.Diff(gotModel.Empty(), true)
	if diff != "" {
		t.Errorf("Expected and actual empty different:\n%s", diff)
	}

	gotModel = ScrolltableModel[TestDatum]{
		height: 5,
		data: []TestDatum{
			{ID: "1"},
		},
	}

	diff = cmp.Diff(gotModel.Empty(), false)
	if diff != "" {
		t.Errorf("Expected and actual empty different:\n%s", diff)
	}
}

func TestSelection(t *testing.T) {
	gotModel := ScrolltableModel[TestDatum]{}

	diff := cmp.Diff(gotModel.Selection(), TestDatum{})
	if diff != "" {
		t.Errorf("Expected and actual selections different:\n%s", diff)
	}

	gotModel = ScrolltableModel[TestDatum]{
		height: 5,
		data: []TestDatum{
			{ID: "1"},
		},
	}

	diff = cmp.Diff(gotModel.Selection(), TestDatum{ID: "1"})
	if diff != "" {
		t.Errorf("Expected and actual selections different:\n%s", diff)
	}
}

func TestIndex(t *testing.T) {
	gotModel := ScrolltableModel[TestDatum]{
		frameStart: 1,
		cursor:     1,
	}

	diff := cmp.Diff(gotModel.Index(), 2)
	if diff != "" {
		t.Errorf("Expected and actual index different:\n%s", diff)
	}
}

func TestFrame(t *testing.T) {
	gotModel := ScrolltableModel[TestDatum]{}

	diff := cmp.Diff(gotModel.Frame(), []TestDatum(nil))
	if diff != "" {
		t.Errorf("Expected and actual frames different:\n%s", diff)
	}

	gotModel = ScrolltableModel[TestDatum]{
		height: 5,
		data: []TestDatum{
			{ID: "1"},
		},
	}

	diff = cmp.Diff(gotModel.Frame(), []TestDatum{{ID: "1"}})
	if diff != "" {
		t.Errorf("Expected and actual frames different:\n%s", diff)
	}
}

func TestData(t *testing.T) {
	gotModel := ScrolltableModel[TestDatum]{
		data: []TestDatum{{ID: "1"}},
	}

	diff := cmp.Diff(gotModel.Data(), []TestDatum{{ID: "1"}})
	if diff != "" {
		t.Errorf("Expected and actual data different:\n%s", diff)
	}
}

type TestDatum struct {
	ID string
}

func verifyUpdate(t *testing.T, gotModel ScrolltableModel[TestDatum], wantModel ScrolltableModel[TestDatum], gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(ScrolltableModel[TestDatum]{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	utils.CompareCommands(t, gotCmd, wantCmd)
}

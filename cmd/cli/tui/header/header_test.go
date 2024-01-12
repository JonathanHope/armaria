package header

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

const Name = "header"

func TestCanUpdateWidth(t *testing.T) {
	gotModel := HeaderModel{
		name: Name,
	}
	gotModel, gotCmd := gotModel.Update(msgs.SizeMsg{Name: Name, Width: 1})

	wantModel := HeaderModel{
		name:  Name,
		width: 1,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd)
}

func TestCanUpdateNav(t *testing.T) {
	gotModel := HeaderModel{}
	gotModel, gotCmd := gotModel.Update(msgs.NavMsg("nav"))

	wantModel := HeaderModel{
		nav: "nav",
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd)
}

func TestCanMarkBusy(t *testing.T) {
	gotModel := HeaderModel{}
	gotModel, gotCmd := gotModel.Update(msgs.BusyMsg{})

	wantModel := HeaderModel{
		busy: true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd)
}

func TestCanMarkFree(t *testing.T) {
	gotModel := HeaderModel{
		busy: true,
	}
	gotModel, gotCmd := gotModel.Update(msgs.FreeMsg{})

	wantModel := HeaderModel{
		busy: false,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd)
}

func TestBusy(t *testing.T) {
	gotModel := HeaderModel{
		busy: true,
	}

	modelDiff := cmp.Diff(gotModel.Busy(), true)
	if modelDiff != "" {
		t.Errorf("Expected and actual busy different:\n%s", modelDiff)
	}
}

func verifyUpdate(t *testing.T, gotModel HeaderModel, wantModel HeaderModel, gotCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(HeaderModel{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	utils.CompareCommands(t, gotCmd, nil)
}

package header

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
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

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanUpdateNav(t *testing.T) {
	gotModel := HeaderModel{}
	gotModel, gotCmd := gotModel.Update(msgs.NavMsg("nav"))

	wantModel := HeaderModel{
		nav: "nav",
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func verifyUpdate(t *testing.T, gotModel HeaderModel, wantModel HeaderModel, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(HeaderModel{})
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

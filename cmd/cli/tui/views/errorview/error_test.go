package errorview

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

func TestHandlesErrorMessage(t *testing.T) {
	err := errors.New("test error")
	var gotModel tea.Model = model{}

	gotModel, gotCmd := gotModel.Update(msgs.ErrorMsg{Err: err})

	wantModel := model{
		err: err,
	}

	wantCmd := func() tea.Msg { return msgs.ViewMsg(msgs.ViewError) }

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func verifyUpdate(t *testing.T, gotModel tea.Model, wantModel tea.Model, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	equateErrors := cmpopts.EquateErrors()
	unexported := cmp.AllowUnexported(model{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported, equateErrors)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	utils.CompareCommands(t, gotCmd, wantCmd)
}

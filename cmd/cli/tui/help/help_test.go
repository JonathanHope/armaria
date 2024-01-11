package help

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

const Name = "help"

func TestCanShowHelp(t *testing.T) {
	gotModel := HelpModel{
		name: Name,
	}
	gotModel, gotCmd := gotModel.Update(msgs.ShowHelpMsg{Name: Name})

	wantModel := HelpModel{
		name:     Name,
		helpMode: true,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanHideHelp(t *testing.T) {
	gotModel := HelpModel{
		name:     Name,
		helpMode: true,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	wantModel := HelpModel{
		name:     Name,
		helpMode: false,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)

	gotModel, gotCmd = gotModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestHelpMode(t *testing.T) {
	gotModel := HelpModel{
		name:     Name,
		helpMode: true,
	}

	modelDiff := cmp.Diff(gotModel.HelpMode(), true)
	if modelDiff != "" {
		t.Errorf("Expected and actual help modes different:\n%s", modelDiff)
	}
}

func verifyUpdate(t *testing.T, gotModel HelpModel, wantModel HelpModel, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(HelpModel{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	utils.CompareCommands(t, gotCmd, wantCmd)
}

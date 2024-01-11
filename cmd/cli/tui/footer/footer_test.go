package footer

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

const Name = "footer"

func TestCanUpdateWidth(t *testing.T) {
	gotModel := FooterModel{
		name: Name,
	}
	gotModel, gotCmd := gotModel.Update(msgs.SizeMsg{Name: Name, Width: 1})

	wantModel := FooterModel{
		name:  Name,
		width: 1,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanStartInputMode(t *testing.T) {
	gotModel := FooterModel{
		name: Name,
	}
	gotModel, gotCmd := gotModel.Update(msgs.InputModeMsg{
		Name:      Name,
		InputMode: true,
		Prompt:    "prompt",
		Text:      "text",
		MaxChars:  5,
	})

	wantModel := FooterModel{
		name:      Name,
		inputMode: true,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.PromptMsg{Name: TextInputName, Prompt: "prompt"} },
			func() tea.Msg { return msgs.TextMsg{Name: TextInputName, Text: "text"} },
			func() tea.Msg { return msgs.FocusMsg{Name: TextInputName, MaxChars: 5} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanEndInputMode(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputMode: true,
	}
	gotModel, gotCmd := gotModel.Update(msgs.InputModeMsg{Name: Name, InputMode: false})

	wantModel := FooterModel{
		name:      Name,
		inputMode: false,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.BlurMsg{Name: TextInputName} },
			func() tea.Msg { return msgs.PromptMsg{Name: TextInputName} },
			func() tea.Msg { return msgs.TextMsg{Name: TextInputName} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanSetFilters(t *testing.T) {
	gotModel := FooterModel{
		name: Name,
	}
	gotModel, gotCmd := gotModel.Update(msgs.FiltersMsg{Name: Name, Filters: []string{"one"}})

	wantModel := FooterModel{
		name:    Name,
		filters: []string{"one"},
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanCancelInput(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputMode: true,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

	wantModel := FooterModel{
		name:      Name,
		inputMode: true,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.InputCancelledMsg{} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanConfirmInput(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputMode: true,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyEnter})

	wantModel := FooterModel{
		name:      Name,
		inputMode: true,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.InputConfirmedMsg{} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func verifyUpdate(t *testing.T, gotModel FooterModel, wantModel FooterModel, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(FooterModel{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported, cmpopts.IgnoreFields(FooterModel{}, "input"))
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	utils.CompareCommands(t, gotCmd, wantCmd)
}

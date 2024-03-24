package footer

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/textinput"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

const Name = "footer"

func TestCanUpdateWidth(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
	}
	gotModel, gotCmd := gotModel.Update(msgs.SizeMsg{Name: Name, Width: 1})

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		width:     1,
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanStartInputMode(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
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
		inputName: Name + "Input",
		inputMode: true,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.PromptMsg{Name: Name + "Input", Prompt: "prompt"} },
			func() tea.Msg { return msgs.TextMsg{Name: Name + "Input", Text: "text"} },
			func() tea.Msg { return msgs.FocusMsg{Name: Name + "Input", MaxChars: 5} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanEndInputMode(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: true,
	}
	gotModel, gotCmd := gotModel.Update(msgs.InputModeMsg{Name: Name, InputMode: false})

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: false,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.BlurMsg{Name: Name + "Input"} },
			func() tea.Msg { return msgs.PromptMsg{Name: Name + "Input"} },
			func() tea.Msg { return msgs.TextMsg{Name: Name + "Input"} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanSetFilters(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
	}
	gotModel, gotCmd := gotModel.Update(msgs.FiltersMsg{Name: Name, Filters: []string{"one"}})

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		filters:   []string{"one"},
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, nil)
}

func TestCanCancelInput(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: true,
	}
	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyEsc})

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: true,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.InputCancelledMsg{Name: Name} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanConfirmInput(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: true,
		input:     textinput.InitialModel(Name+"Input", "> "),
	}

	gotModel.input, _ = gotModel.input.Update(msgs.SizeMsg{
		Name:  Name + "Input",
		Width: 12,
	})

	gotModel.input, _ = gotModel.input.Update(msgs.TextMsg{
		Name: Name + "Input",
		Text: "text",
	})

	gotModel, gotCmd := gotModel.Update(tea.KeyMsg{Type: tea.KeyEnter})

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: true,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.InputConfirmedMsg{Name: Name} },
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

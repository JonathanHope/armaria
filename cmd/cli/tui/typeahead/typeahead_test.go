package typeahead

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
)

const name = "typeaheaed"
const inputName = "input"
const tableName = "table"
const operation = "operation"
const prompt = ">"
const text = "text"
const maxChars = 5

func TestTypeaheadMode(t *testing.T) {
	model := TypeaheadModel{
		typeaheadMode: true,
	}

	diff := cmp.Diff(model.TypeaheadMode(), true)
	if diff != "" {
		t.Errorf("Expected and actual typeaheadmode different")
	}
}

func TestCanSwitchToTypeaheadMode(t *testing.T) {
	gotModel := TypeaheadModel{
		name:      name,
		inputName: inputName,
		tableName: tableName,
	}

	gotModel, cmd := gotModel.Update(msgs.TypeaheadModeMsg{
		Name:           name,
		InputMode:      true,
		MinFilterChars: 3,
		Operation:      operation,
		Prompt:         prompt,
		Text:           "text",
		MaxChars:       5,
		UnfilteredQuery: func() ([]msgs.TypeaheadItem, error) {
			return []msgs.TypeaheadItem{}, nil
		},
		FilteredQuery: func(query string) ([]msgs.TypeaheadItem, error) {
			return []msgs.TypeaheadItem{}, nil
		},
	})

	wantModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
		operation:      operation,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.PromptMsg{Name: inputName, Prompt: prompt} },
			func() tea.Msg { return msgs.TextMsg{Name: inputName, Text: text} },
			func() tea.Msg { return msgs.FocusMsg{Name: inputName, MaxChars: maxChars} },
			func() tea.Msg {
				return msgs.DataMsg[msgs.TypeaheadItem]{
					Name: tableName,
					Data: []msgs.TypeaheadItem{},
					Move: msgs.DirectionStart,
				}
			},
		}
	}

	verifyUpdate(t, gotModel, wantModel, cmd, wantCmd)
}

func TestCanSwitchFromTypeaheadMode(t *testing.T) {
	gotModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
		operation:      operation,
	}

	gotModel, cmd := gotModel.Update(msgs.TypeaheadModeMsg{
		Name:      name,
		InputMode: false,
	})

	wantModel := TypeaheadModel{
		name:          name,
		inputName:     inputName,
		tableName:     tableName,
		typeaheadMode: false,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.BlurMsg{Name: inputName} },
			func() tea.Msg { return msgs.PromptMsg{Name: inputName, Prompt: ""} },
			func() tea.Msg { return msgs.TextMsg{Name: inputName, Text: ""} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, cmd, wantCmd)
}

func TestCanCancelTypeahead(t *testing.T) {
	gotModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
		operation:      operation,
	}

	gotModel, cmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEsc}))

	wantModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
		operation:      operation,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.TypeaheadCancelledMsg{Name: name} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, cmd, wantCmd)
}

func TestCanConfirmTypeahead(t *testing.T) {
	gotModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
		operation:      operation,
	}

	gotModel, cmd := gotModel.Update(tea.KeyMsg(tea.Key{Type: tea.KeyEnter}))

	wantModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
		operation:      operation,
	}

	wantCmd := func() tea.Msg {
		return tea.BatchMsg{
			func() tea.Msg { return msgs.TypeaheadConfirmedMsg{Name: name, Operation: operation} },
		}
	}

	verifyUpdate(t, gotModel, wantModel, cmd, wantCmd)
}

func verifyUpdate(t *testing.T, gotModel TypeaheadModel, wantModel TypeaheadModel, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	unexported := cmp.AllowUnexported(TypeaheadModel{})
	modelDiff := cmp.Diff(
		gotModel,
		wantModel,
		unexported,
		cmpopts.IgnoreFields(TypeaheadModel{}, "input"),
		cmpopts.IgnoreFields(TypeaheadModel{}, "table"),
		cmpopts.IgnoreFields(TypeaheadModel{}, "filteredQuery"),
		cmpopts.IgnoreFields(TypeaheadModel{}, "unfilteredQuery"),
	)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}

	utils.CompareCommands(t, gotCmd, wantCmd)
}

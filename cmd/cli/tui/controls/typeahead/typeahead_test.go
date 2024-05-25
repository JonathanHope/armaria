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
	gotModel.input.Resize(15)

	gotCmd := gotModel.StartTypeahead(StartTypeaheadPayload{
		MinFilterChars: 3,
		Prompt:         prompt,
		Text:           "text",
		MaxChars:       5,
		UnfilteredQuery: func() ([]TypeaheadItem, error) {
			return []TypeaheadItem{}, nil
		},
		FilteredQuery: func(query string) ([]TypeaheadItem, error) {
			return []TypeaheadItem{}, nil
		},
	})

	wantModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
	}

	wantCmd := func() tea.Msg {
		return msgs.DataMsg[TypeaheadItem]{
			Name: tableName,
			Data: []TypeaheadItem{},
			Move: msgs.DirectionStart,
		}
	}

	verifyUpdate(t, gotModel, wantModel, gotCmd, wantCmd)
}

func TestCanSwitchFromTypeaheadMode(t *testing.T) {
	gotModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  true,
		minFilterChars: 3,
	}
	gotModel.input.Resize(15)

	gotModel.StopTypeahead()

	wantModel := TypeaheadModel{
		name:           name,
		inputName:      inputName,
		tableName:      tableName,
		typeaheadMode:  false,
		minFilterChars: 3,
	}

	verifyUpdate(t, gotModel, wantModel, nil, nil)
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

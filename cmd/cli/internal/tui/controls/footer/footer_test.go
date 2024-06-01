package footer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/controls/textinput"
)

const Name = "footer"

func TestCanUpdateWidth(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		input:     textinput.TextInputModel{},
	}
	gotModel.Resize(15)

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		width:     15,
	}

	verifyUpdate(t, gotModel, wantModel)
}

func TestCanStartInputMode(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
	}
	gotModel.Resize(100)
	gotModel.StartInputMode("prompt", "text", 15)

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: true,
		width:     100,
	}

	verifyUpdate(t, gotModel, wantModel)
}

func TestCanEndInputMode(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: true,
	}
	gotModel.Resize(100)
	gotModel.StopInputMode()

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		inputMode: false,
		width:     100,
	}

	verifyUpdate(t, gotModel, wantModel)
}

func TestCanSetFilters(t *testing.T) {
	gotModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
	}
	gotModel.SetFilters([]string{"one"})

	wantModel := FooterModel{
		name:      Name,
		inputName: Name + "Input",
		filters:   []string{"one"},
	}

	verifyUpdate(t, gotModel, wantModel)
}

func verifyUpdate(t *testing.T, gotModel FooterModel, wantModel FooterModel) {
	unexported := cmp.AllowUnexported(FooterModel{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported, cmpopts.IgnoreFields(FooterModel{}, "input"))
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}
}

package header

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const Name = "header"

func TestCanUpdateWidth(t *testing.T) {
	gotModel := HeaderModel{
		name: Name,
	}
	gotModel.Resize(1)

	wantModel := HeaderModel{
		name:  Name,
		width: 1,
	}

	verifyUpdate(t, gotModel, wantModel)
}

func TestCanUpdateNav(t *testing.T) {
	gotModel := HeaderModel{}
	gotModel.SetBreadcrumbs("breadcrumbs")

	wantModel := HeaderModel{
		breadcrumbs: "breadcrumbs",
	}

	verifyUpdate(t, gotModel, wantModel)
}

func TestCanMarkBusy(t *testing.T) {
	gotModel := HeaderModel{}
	gotModel.SetBusy()

	wantModel := HeaderModel{
		busy: true,
	}

	verifyUpdate(t, gotModel, wantModel)
}

func TestCanMarkFree(t *testing.T) {
	gotModel := HeaderModel{
		busy: true,
	}
	gotModel.SetFree()

	wantModel := HeaderModel{
		busy: false,
	}

	verifyUpdate(t, gotModel, wantModel)
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

func verifyUpdate(t *testing.T, gotModel HeaderModel, wantModel HeaderModel) {
	unexported := cmp.AllowUnexported(HeaderModel{})
	modelDiff := cmp.Diff(gotModel, wantModel, unexported)
	if modelDiff != "" {
		t.Errorf("Expected and actual models different:\n%s", modelDiff)
	}
}

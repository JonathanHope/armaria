package utils

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
)

func CompareCommands(t *testing.T, gotCmd tea.Cmd, wantCmd tea.Cmd) {
	if gotCmd == nil || wantCmd == nil {
		if gotCmd != nil || wantCmd != nil {
			t.Errorf("Expected and actual cmds different: one is nil and one is non-nil")
		}
	} else if gotBatchCmd, ok := gotCmd().(tea.BatchMsg); ok {
		if wantBatchCmd, ok := wantCmd().(tea.BatchMsg); ok {
			gotBatches := []tea.BatchMsg{gotBatchCmd}
			wantBatches := []tea.BatchMsg{wantBatchCmd}
			var gotResults []tea.Msg
			var wantResults []tea.Msg

			for len(gotBatches) > 0 || len(wantBatches) > 0 {
				if len(gotBatches) > 0 {
					var gotBatch tea.BatchMsg
					gotBatch, gotBatches = gotBatches[0], gotBatches[1:]

					for _, got := range gotBatch {
						result := got()

						if batchResult, ok := result.(tea.BatchMsg); ok {
							gotBatches = append(gotBatches, batchResult)
						} else {
							gotResults = append(gotResults, result)
						}
					}
				}
				if len(wantBatches) > 0 {
					var wantBatch tea.BatchMsg
					wantBatch, wantBatches = wantBatches[0], wantBatches[1:]

					for _, want := range wantBatch {
						result := want()

						if batchResult, ok := result.(tea.BatchMsg); ok {
							wantBatches = append(wantBatches, batchResult)
						} else {
							wantResults = append(wantResults, result)
						}
					}
				}
			}

			cmdDiff := cmp.Diff(gotResults, wantResults)
			if cmdDiff != "" {
				t.Errorf("Expected and actual cmds different:\n%s", cmdDiff)
			}
		} else {
			t.Errorf("Expected and actual cmds different: one was batch one wasn't")
		}
	} else {
		cmdDiff := cmp.Diff(gotCmd(), wantCmd())
		if cmdDiff != "" {
			t.Errorf("Expected and actual cmds different:\n%s", cmdDiff)
		}
	}
}

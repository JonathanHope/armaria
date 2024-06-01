package messaging

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jonathanhope/armaria/internal/null"
)

func TestSendReceiveLoop(t *testing.T) {
	payload := AddBookPayload{
		URL:         "https://test.com",
		Name:        null.NullStringFromPtr(nil),
		Description: null.NullStringFromPtr(nil),
		ParentID:    null.NullStringFromPtr(nil),
		Tags:        nil,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expected := NativeMessage{
		Kind:    "add-book",
		Payload: string(json),
	}

	buffer := bytes.NewBuffer(nil)

	err = expected.SendMessage(buffer)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	actual, err := ReceiveMessage(buffer)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(expected, actual)
	if diff != "" {
		t.Errorf("Expected and actual messages different:\n%s", diff)
	}
}

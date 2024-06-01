package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jonathanhope/armaria/cmd/cli/internal/messaging"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg"
)

func TestRemoveFolder(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armaria.DefaultAddFolderOptions()
	options.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindRemoveFolder, messaging.RemoveFolderPayload{
		DB: null.NullStringFrom(db),
		ID: folder.ID,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindVoid, messaging.VoidPayload{})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

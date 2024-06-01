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

func TestAddFolder(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	got, err := nativeMessageLoop(messaging.MessageKindAddFolder, messaging.AddFolderPayload{
		DB:   null.NullStringFrom(db),
		Name: "Blogs",
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	id, err := getLastInsertedID(db, []string{})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:       id,
			Name:     "Blogs",
			IsFolder: true,
			Tags:     []string{},
		},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

func TestAddFolderToFolder(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armaria.DefaultAddFolderOptions()
	options.WithDB(db)
	folder, err := armaria.AddFolder("Programming", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindAddFolder, messaging.AddFolderPayload{
		DB:       null.NullStringFrom(db),
		Name:     "Blogs",
		ParentID: null.NullStringFrom(folder.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	id, err := getLastInsertedID(db, []string{folder.ID})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:         id,
			Name:       "Blogs",
			IsFolder:   true,
			Tags:       []string{},
			ParentID:   null.NullStringFrom(folder.ID),
			ParentName: null.NullStringFrom("Programming"),
		},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

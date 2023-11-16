package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jonathanhope/armaria/cmd/host/internal/messaging"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg/api"
)

func TestUpdateFolder(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddFolderOptions()
	options.WithDB(db)
	folder, err := armariaapi.AddFolder("Blogs", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateFolder, messaging.UpdateFolderPayload{
		DB:   null.NullStringFrom(db),
		ID:   folder.ID,
		Name: null.NullStringFrom("Programming"),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:       folder.ID,
			Name:     "Programming",
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

func TestUpdateFolderParentID(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddFolderOptions()
	options.WithDB(db)

	programming, err := armariaapi.AddFolder("Programming", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	blogs, err := armariaapi.AddFolder("Blogs", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateFolder, messaging.UpdateFolderPayload{
		DB:       null.NullStringFrom(db),
		ID:       blogs.ID,
		ParentID: null.NullStringFrom(programming.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:         blogs.ID,
			Name:       "Blogs",
			IsFolder:   true,
			Tags:       []string{},
			ParentID:   null.NullStringFrom(programming.ID),
			ParentName: null.NullStringFrom(programming.Name),
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

func TestUpdateFolderRemoveParentID(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddFolderOptions()
	options.WithDB(db)

	programming, err := armariaapi.AddFolder("Programming", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	options.WithParentID(programming.ID)

	blogs, err := armariaapi.AddFolder("Blogs", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateFolder, messaging.UpdateFolderPayload{
		DB:             null.NullStringFrom(db),
		ID:             blogs.ID,
		RemoveParentID: true,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:       blogs.ID,
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

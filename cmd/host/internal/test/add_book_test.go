package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jonathanhope/armaria/cmd/host/internal/messaging"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg"
)

func TestAddBook(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	got, err := nativeMessageLoop(messaging.MessageKindAddBook, messaging.AddBookPayload{
		DB:  null.NullStringFrom(db),
		URL: "https://jho.pe",
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
			ID:   id,
			URL:  null.NullStringFrom("https://jho.pe"),
			Name: "https://jho.pe",
			Tags: []string{},
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

func TestAddBookToFolder(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armaria.DefaultAddFolderOptions()
	options.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindAddBook, messaging.AddBookPayload{
		DB:       null.NullStringFrom(db),
		URL:      "https://jho.pe",
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
			URL:        null.NullStringFrom("https://jho.pe"),
			Name:       "https://jho.pe",
			Tags:       []string{},
			ParentID:   null.NullStringFrom(folder.ID),
			ParentName: null.NullStringFrom("Blogs"),
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

func TestAddBookWithName(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	got, err := nativeMessageLoop(messaging.MessageKindAddBook, messaging.AddBookPayload{
		DB:   null.NullStringFrom(db),
		URL:  "https://jho.pe",
		Name: null.NullStringFrom("The Flat Field"),
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
			ID:   id,
			URL:  null.NullStringFrom("https://jho.pe"),
			Name: "The Flat Field",
			Tags: []string{},
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

func TestAddBookWithDescription(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	got, err := nativeMessageLoop(messaging.MessageKindAddBook, messaging.AddBookPayload{
		DB:          null.NullStringFrom(db),
		URL:         "https://jho.pe",
		Description: null.NullStringFrom("The Blog of Jonathan Hope."),
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
			ID:          id,
			URL:         null.NullStringFrom("https://jho.pe"),
			Name:        "https://jho.pe",
			Description: null.NullStringFrom("The Blog of Jonathan Hope."),
			Tags:        []string{},
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

func TestAddBookWithTags(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	got, err := nativeMessageLoop(messaging.MessageKindAddBook, messaging.AddBookPayload{
		DB:   null.NullStringFrom(db),
		URL:  "https://jho.pe",
		Tags: []string{"blog", "programming"},
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
			ID:   id,
			URL:  null.NullStringFrom("https://jho.pe"),
			Name: "https://jho.pe",
			Tags: []string{"blog", "programming"},
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

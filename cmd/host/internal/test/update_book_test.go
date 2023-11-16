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

func TestUpdateBookURL(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddBookOptions()
	options.WithDB(db)
	book, err := armariaapi.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateBook, messaging.UpdateBookPayload{
		DB:  null.NullStringFrom(db),
		ID:  book.ID,
		URL: null.NullStringFrom("https://theflatfield.net"),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:   book.ID,
			URL:  null.NullStringFrom("https://theflatfield.net"),
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

func TestUpdateBookName(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddBookOptions()
	options.WithDB(db)
	book, err := armariaapi.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateBook, messaging.UpdateBookPayload{
		DB:   null.NullStringFrom(db),
		ID:   book.ID,
		Name: null.NullStringFrom("The Flat Field"),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:   book.ID,
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

func TestUpdateBookDescription(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddBookOptions()
	options.WithDB(db)
	book, err := armariaapi.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateBook, messaging.UpdateBookPayload{
		DB:          null.NullStringFrom(db),
		ID:          book.ID,
		Description: null.NullStringFrom("The blog of Jonathan Hope."),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:          book.ID,
			URL:         null.NullStringFrom("https://jho.pe"),
			Name:        "https://jho.pe",
			Tags:        []string{},
			Description: null.NullStringFrom("The blog of Jonathan Hope."),
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

func TestUpdateBookRemoveDescription(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddBookOptions()
	options.WithDB(db)
	options.WithDescription("The blog of Jonathan Hope.")
	book, err := armariaapi.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateBook, messaging.UpdateBookPayload{
		DB:                null.NullStringFrom(db),
		ID:                book.ID,
		RemoveDescription: true,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:   book.ID,
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

func TestUpdateBookParentID(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armariaapi.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armariaapi.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	options := armariaapi.DefaultAddBookOptions()
	options.WithDB(db)
	book, err := armariaapi.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateBook, messaging.UpdateBookPayload{
		DB:       null.NullStringFrom(db),
		ID:       book.ID,
		ParentID: null.NullStringFrom(folder.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:         book.ID,
			URL:        null.NullStringFrom("https://jho.pe"),
			Name:       "https://jho.pe",
			Tags:       []string{},
			ParentID:   null.NullStringFrom(folder.ID),
			ParentName: null.NullStringFrom(folder.Name),
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

func TestUpdateBookRemoveParentID(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armariaapi.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armariaapi.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	options := armariaapi.DefaultAddBookOptions()
	options.WithDB(db)
	options.WithParentID(folder.ID)
	book, err := armariaapi.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindUpdateBook, messaging.UpdateBookPayload{
		DB:             null.NullStringFrom(db),
		ID:             book.ID,
		RemoveParentID: true,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBook, messaging.BookPayload{
		Book: messaging.BookDTO{
			ID:   book.ID,
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

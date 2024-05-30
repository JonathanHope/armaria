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

func TestRemoveTags(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armaria.DefaultAddBookOptions()
	options.WithDB(db)
	options.WithTags([]string{"blog", "programming"})
	book, err := armaria.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindRemoveTags, messaging.RemoveTagsPayload{
		DB:   null.NullStringFrom(db),
		ID:   book.ID,
		Tags: []string{"blog", "programming"},
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

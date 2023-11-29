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

func TestAddTags(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armariaapi.DefaultAddBookOptions()
	options.WithDB(db)
	book, err := armariaapi.AddBook("https://jho.pe", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindAddTags, messaging.AddTagsPayload{
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

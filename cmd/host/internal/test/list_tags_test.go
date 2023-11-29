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
	"github.com/jonathanhope/armaria/pkg/model"
)

func TestListTags(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	bookOptions := armariaapi.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithTags([]string{"blog", "programming"})
	_, err := armariaapi.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListTags, messaging.ListTagsPayload{
		DB: null.NullStringFrom(db),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindTags, messaging.TagsPayload{
		Tags: []string{"blog", "programming"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

func TestListTagsAsc(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	bookOptions := armariaapi.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithTags([]string{"blog", "programming"})
	_, err := armariaapi.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListTags, messaging.ListTagsPayload{
		DB:        null.NullStringFrom(db),
		Direction: string(armaria.DirectionAsc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindTags, messaging.TagsPayload{
		Tags: []string{"blog", "programming"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

func TestListTagsDesc(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	bookOptions := armariaapi.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithTags([]string{"blog", "programming"})
	_, err := armariaapi.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListTags, messaging.ListTagsPayload{
		DB:        null.NullStringFrom(db),
		Direction: string(armaria.DirectionDesc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindTags, messaging.TagsPayload{
		Tags: []string{"programming", "blog"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

func TestListTagsWithLimit(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	bookOptions := armariaapi.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithTags([]string{"blog", "programming"})
	_, err := armariaapi.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListTags, messaging.ListTagsPayload{
		DB:    null.NullStringFrom(db),
		First: null.NullInt64From(1),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindTags, messaging.TagsPayload{
		Tags: []string{"blog"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

func TestListTagsWithAfter(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	bookOptions := armariaapi.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithTags([]string{"blog", "programming"})
	_, err := armariaapi.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListTags, messaging.ListTagsPayload{
		DB:    null.NullStringFrom(db),
		After: null.NullStringFrom("blog"),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindTags, messaging.TagsPayload{
		Tags: []string{"programming"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

func TestListTagsWithQuery(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	bookOptions := armariaapi.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithTags([]string{"blog", "programming"})
	_, err := armariaapi.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListTags, messaging.ListTagsPayload{
		DB:    null.NullStringFrom(db),
		Query: null.NullStringFrom("gram"),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindTags, messaging.TagsPayload{
		Tags: []string{"programming"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual message different:\n%s", diff)
	}
}

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
	"github.com/samber/lo"
)

func TestUpdateFolder(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	options := armaria.DefaultAddFolderOptions()
	options.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", options)
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

	options := armaria.DefaultAddFolderOptions()
	options.WithDB(db)

	programming, err := armaria.AddFolder("Programming", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	blogs, err := armaria.AddFolder("Blogs", options)
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

	options := armaria.DefaultAddFolderOptions()
	options.WithDB(db)

	programming, err := armaria.AddFolder("Programming", options)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	options.WithParentID(programming.ID)

	blogs, err := armaria.AddFolder("Blogs", options)
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

func TestUpdateFolderOrderStart(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	addOptions := armaria.DefaultAddFolderOptions()
	addOptions.WithDB(db)
	folder1, err := armaria.AddFolder("one", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	folder2, err := armaria.AddFolder("two", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	folder3, err := armaria.AddFolder("three", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = nativeMessageLoop(messaging.MessageKindUpdateFolder, messaging.UpdateFolderPayload{
		DB:       null.NullStringFrom(db),
		ID:       folder2.ID,
		NextBook: null.NullStringFrom(folder1.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	listOptions := armaria.DefaultListBooksOptions()
	listOptions.WithFolders(true)
	listOptions.WithDB(db)
	books, err := armaria.ListBooks(listOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got := lo.Map(books, func(x armaria.Book, _ int) string {
		return x.ID
	})

	want := []string{
		folder2.ID,
		folder1.ID,
		folder3.ID,
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual orders different:\n%s", diff)
	}
}

func TestUpdateFolderOrderEnd(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	addOptions := armaria.DefaultAddFolderOptions()
	addOptions.WithDB(db)
	folder1, err := armaria.AddFolder("one", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	folder2, err := armaria.AddFolder("two", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	folder3, err := armaria.AddFolder("three", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = nativeMessageLoop(messaging.MessageKindUpdateFolder, messaging.UpdateFolderPayload{
		DB:           null.NullStringFrom(db),
		ID:           folder2.ID,
		PreviousBook: null.NullStringFrom(folder3.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	listOptions := armaria.DefaultListBooksOptions()
	listOptions.WithFolders(true)
	listOptions.WithDB(db)
	books, err := armaria.ListBooks(listOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got := lo.Map(books, func(x armaria.Book, _ int) string {
		return x.ID
	})

	want := []string{
		folder1.ID,
		folder3.ID,
		folder2.ID,
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual orders different:\n%s", diff)
	}
}

func TestUpdateFolderOrderBetween(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	addOptions := armaria.DefaultAddFolderOptions()
	addOptions.WithDB(db)
	folder1, err := armaria.AddFolder("one", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	folder2, err := armaria.AddFolder("two", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	folder3, err := armaria.AddFolder("three", addOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = nativeMessageLoop(messaging.MessageKindUpdateFolder, messaging.UpdateFolderPayload{
		DB:           null.NullStringFrom(db),
		ID:           folder3.ID,
		PreviousBook: null.NullStringFrom(folder1.ID),
		NextBook:     null.NullStringFrom(folder2.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	listOptions := armaria.DefaultListBooksOptions()
	listOptions.WithBooks(true)
	listOptions.WithDB(db)
	books, err := armaria.ListBooks(listOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got := lo.Map(books, func(x armaria.Book, _ int) string {
		return x.ID
	})

	want := []string{
		folder1.ID,
		folder3.ID,
		folder2.ID,
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("Expected and actual orders different:\n%s", diff)
	}
}

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

func TestListBooksFolders(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:       folder.ID,
				Name:     "Blogs",
				IsFolder: true,
				Tags:     []string{},
			},
			{
				ID:   book.ID,
				URL:  null.NullStringFrom("https://jho.pe"),
				Name: "https://jho.pe",
				Tags: []string{},
			},
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

func TestListBooks(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	_, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   false,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:   book.ID,
				URL:  null.NullStringFrom("https://jho.pe"),
				Name: "https://jho.pe",
				Tags: []string{},
			},
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

func TestListFolder(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	_, err = armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: false,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:       folder.ID,
				Name:     "Blogs",
				IsFolder: true,
				Tags:     []string{},
			},
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

func TestListBooksFoldersWithLimit(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	_, err = armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		First:            null.NullInt64From(1),
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:       folder.ID,
				Name:     "Blogs",
				IsFolder: true,
				Tags:     []string{},
			},
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

func TestListBooksFoldersNamAsc(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:       folder.ID,
				Name:     "Blogs",
				IsFolder: true,
				Tags:     []string{},
			},
			{
				ID:   book.ID,
				URL:  null.NullStringFrom("https://jho.pe"),
				Name: "https://jho.pe",
				Tags: []string{},
			},
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

func TestListBooksFoldersNamDesc(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionDesc),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:   book.ID,
				URL:  null.NullStringFrom("https://jho.pe"),
				Name: "https://jho.pe",
				Tags: []string{},
			},
			{
				ID:       folder.ID,
				Name:     "Blogs",
				IsFolder: true,
				Tags:     []string{},
			},
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

func TestListBooksFoldersWithAfter(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
		After:            null.NullStringFrom(folder.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:   book.ID,
				URL:  null.NullStringFrom("https://jho.pe"),
				Name: "https://jho.pe",
				Tags: []string{},
			},
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

func TestListBooksFoldersWithParentID(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithParentID(folder.ID)
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
		ParentID:         null.NullStringFrom(folder.ID),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:         book.ID,
				URL:        null.NullStringFrom("https://jho.pe"),
				Name:       "https://jho.pe",
				Tags:       []string{},
				ParentID:   null.NullStringFrom(folder.ID),
				ParentName: null.NullStringFrom(folder.Name),
			},
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

func TestListBooksFoldersWithoutParentID(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	folder, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithParentID(folder.ID)
	_, err = armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
		WithoutParentID:  true,
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:       folder.ID,
				Name:     "Blogs",
				IsFolder: true,
				Tags:     []string{},
			},
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

func TestListBooksFoldersWithQuery(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	_, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
		Query:            null.NullStringFrom("jho"),
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:   book.ID,
				URL:  null.NullStringFrom("https://jho.pe"),
				Name: "https://jho.pe",
				Tags: []string{},
			},
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

func TestListBooksFoldersWithTags(t *testing.T) {
	db := fmt.Sprintf("%s.db", uuid.New().String())
	defer func() { os.Remove(db) }()

	folderOptions := armaria.DefaultAddFolderOptions()
	folderOptions.WithDB(db)
	_, err := armaria.AddFolder("Blogs", folderOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	bookOptions := armaria.DefaultAddBookOptions()
	bookOptions.WithDB(db)
	bookOptions.WithTags([]string{"blog", "programming"})
	book, err := armaria.AddBook("https://jho.pe", bookOptions)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := nativeMessageLoop(messaging.MessageKindListBooks, messaging.ListBooksPayload{
		DB:               null.NullStringFrom(db),
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            string(armaria.OrderName),
		Direction:        string(armaria.DirectionAsc),
		Tags:             []string{"blog"},
	})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want, err := messaging.PayloadToMessage(messaging.MessageKindBooks, messaging.BooksPayload{
		Books: []messaging.BookDTO{
			{
				ID:   book.ID,
				URL:  null.NullStringFrom("https://jho.pe"),
				Name: "https://jho.pe",
				Tags: []string{"blog", "programming"},
			},
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

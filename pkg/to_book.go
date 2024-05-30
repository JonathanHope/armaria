package armaria

import (
	"strings"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/samber/lo"
)

// toBooks converts a slice of BookDTO to a slice of Books.
func toBooks(books []db.BookDTO) []Book {
	return lo.Map(books, func(book db.BookDTO, _ int) Book {
		return toBook(book)
	})
}

// toBook converts a BookDTO to a Book.
func toBook(book db.BookDTO) Book {
	return Book{
		ID:          book.ID,
		URL:         null.PtrFromNullString(book.URL),
		Name:        book.Name,
		Description: null.PtrFromNullString(book.Description),
		ParentID:    null.PtrFromNullString(book.ParentID),
		IsFolder:    book.IsFolder,
		Order:       book.Order,
		ParentName:  null.PtrFromNullString(book.ParentName),
		Tags:        parseTags(book.Tags),
	}
}

// parseTags parses the tags coming back from the database.
func parseTags(tags string) []string {
	if tags == "" {
		return make([]string, 0)
	}

	return strings.Split(tags, ",")
}

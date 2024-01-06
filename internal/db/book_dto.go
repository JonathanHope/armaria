package db

import (
	"strings"

	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg/model"
)

// bookDTO is a DTO to stuff DB results into.
type bookDTO struct {
	ID          string          `db:"id"`
	URL         null.NullString `db:"url"`
	Name        string          `db:"name"`
	Description null.NullString `db:"description"`
	ParentID    null.NullString `db:"parent_id"`
	IsFolder    bool            `db:"is_folder"`
	Order       string          `db:"order"`
	ParentName  null.NullString `db:"parent_name"`
	Tags        string          `db:"tags"`
}

// toBook converts a bookDTO to a Book.
func (book bookDTO) toBook() armaria.Book {
	return armaria.Book{
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

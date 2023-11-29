package messaging

import (
	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/null"
)

// BookDTO is a bookmark or folder that can be marshalled into JSON.
type BookDTO struct {
	ID          string          `json:"id"`
	URL         null.NullString `json:"url"`
	Name        string          `json:"name"`
	Description null.NullString `json:"description"`
	ParentID    null.NullString `json:"parentId"`
	IsFolder    bool            `json:"isFolder"`
	ParentName  null.NullString `json:"parentName"`
	Tags        []string        `json:"tags"`
}

// bookMapper maps a Book to a BookDTO.
func bookMapper(book armaria.Book) BookDTO {
	return BookDTO{
		ID:          book.ID,
		URL:         null.NullStringFromPtr(book.URL),
		Name:        book.Name,
		Description: null.NullStringFromPtr(book.Description),
		ParentID:    null.NullStringFromPtr(book.ParentID),
		IsFolder:    book.IsFolder,
		ParentName:  null.NullStringFromPtr(book.ParentName),
		Tags:        book.Tags,
	}
}

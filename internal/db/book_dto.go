package db

import (
	"github.com/jonathanhope/armaria/internal/null"
)

// BookDTO is a DTO to stuff DB results into.
type BookDTO struct {
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

package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg/model"
)

// getBookOptions are the optional arguments for GetBook.
type getBookOptions struct {
	DB null.NullString
}

// DefaultGetBookOptions are the default options for GetBook.
func DefaultGetBookOptions() *getBookOptions {
	return &getBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *getBookOptions) WithDB(db string) *getBookOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// GetBook gets a bookmark in the bookmarks database.
func GetBook(id string, options *getBookOptions) (book armaria.Book, err error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return armaria.Book{}, fmt.Errorf("error getting config while getting bookmark: %w", err)
	}

	return db.QueryWithDB(options.DB, config.DB, func(tx db.Transaction) (armaria.Book, error) {
		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:       id,
			IncludeBooks:   true,
			IncludeFolders: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting bookmarks while getting bookmark: %w", err)
		}

		if len(books) == 0 {
			return armaria.Book{}, armaria.ErrNotFound
		}

		return books[0], nil
	})
}

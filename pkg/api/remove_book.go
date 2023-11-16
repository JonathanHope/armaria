package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
)

// removeBookOptions are the optional arguments for RemoveBook.
type removeBookOptions struct {
	DB null.NullString
}

// DefaultRemoveBookOptions are the default options for RemoveBook.
func DefaultRemoveBookOptions() removeBookOptions {
	return removeBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *removeBookOptions) WithDB(db string) {
	o.DB = null.NullStringFrom(db)
}

// RemoveBook removes a bookmark from the bookmarks database.
func RemoveBook(id string, options removeBookOptions) (err error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return fmt.Errorf("error getting config while removing bookmark: %w", err)
	}

	return db.ExecWithTransaction(options.DB, config.DB, func(tx db.Transaction) error {
		if err := validate.BookID(tx, id); err != nil {
			return fmt.Errorf("bookmark ID validation failed while removing bookmark: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return fmt.Errorf("error getting bookmarks while removing bookmark: %w", err)
		}
		book := books[0]

		if err = db.UnlinkTags(tx, book.ID, book.Tags); err != nil {
			return fmt.Errorf("error unlinking tags while removing bookmark: %w", err)
		}

		if err = db.RemoveBook(tx, book.ID); err != nil {
			return fmt.Errorf("error while removing bookmark: %w", err)
		}

		if err = db.CleanOrphanedTags(tx, book.Tags); err != nil {
			return fmt.Errorf("error cleaning orphaned tags while removing bookmark: %w", err)
		}

		return nil
	})
}

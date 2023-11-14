package lib

import (
	"fmt"
)

// removeBookOptions are the optional arguments for RemoveBook.
type removeBookOptions struct {
	db NullString
}

// DefaultRemoveBookOptions are the default options for RemoveBook.
func DefaultRemoveBookOptions() removeBookOptions {
	return removeBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *removeBookOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// RemoveBook removes a bookmark from the bookmarks database.
func RemoveBook(id string, options removeBookOptions) (err error) {
	return execWithTransaction(options.db, connectDB, func(tx transaction) error {
		if err := validateBookID(tx, id); err != nil {
			return fmt.Errorf("bookmark ID validation failed while removing bookmark: %w", err)
		}

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return fmt.Errorf("error getting bookmarks while removing bookmark: %w", err)
		}
		book := books[0]

		if err = unlinkTagsDB(tx, book.ID, book.Tags); err != nil {
			return fmt.Errorf("error unlinking tags while removing bookmark: %w", err)
		}

		if err = removeBookDB(tx, book.ID); err != nil {
			return fmt.Errorf("error while removing bookmark: %w", err)
		}

		if err = cleanOrphanedTagsDB(tx, book.Tags); err != nil {
			return fmt.Errorf("error cleaning orphaned tags while removing bookmark: %w", err)
		}

		return nil
	})
}

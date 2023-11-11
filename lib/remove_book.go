package lib

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
			return err
		}

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return err
		}
		book := books[0]

		if err = unlinkTagsDB(tx, book.ID, book.Tags); err != nil {
			return err
		}

		if err = removeBookDB(tx, book.ID); err != nil {
			return err
		}

		if err = cleanOrphanedTagsDB(tx, book.Tags); err != nil {
			return err
		}

		return nil
	})
}

package lib

import (
	"fmt"
)

// addFolderOptions are the optional arguments for AddFolder.
type addFolderOptions struct {
	db       NullString
	parentID NullString
}

// DefaultAddFolderOptions are the default options for AddFolder.
func DefaultAddFolderOptions() addFolderOptions {
	return addFolderOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addFolderOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// WithParentID sets the folders' parent ID.
func (o *addFolderOptions) WithParentID(parentID string) {
	o.parentID = NullStringFrom(parentID)
}

// AddFolder adds a folder to the bookmarks database.
func AddFolder(name string, options addFolderOptions) (Book, error) {
	return queryWithTransaction(options.db, connectDB, func(tx transaction) (Book, error) {
		var book Book

		if err := validateName(NullStringFrom(name)); err != nil {
			return book, fmt.Errorf("name validation failed while adding folder: %w", err)
		}

		if err := validateParentID(tx, options.parentID); err != nil {
			return book, fmt.Errorf("parent ID validation failed while adding folder: %w", err)
		}

		id, err := addFolderDB(tx, name, options.parentID)
		if err != nil {
			return book, fmt.Errorf("error while adding folder: %w", err)
		}

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:       id,
			includeFolders: true,
		})
		if err != nil {
			return book, fmt.Errorf("error getting folders while adding folder: %w", err)
		}

		return books[0], nil
	})
}

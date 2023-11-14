package lib

import (
	"fmt"
)

// updateFolderOptions are the optional arguments for UpdateFolder.
type updateFolderOptions struct {
	db       NullString
	name     NullString
	parentID NullString
}

// DefaultUpdateFolderOptions are the default options for UpdateFolder.
func DefaultUpdateFolderOptions() updateFolderOptions {
	return updateFolderOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *updateFolderOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// WithName updates the name of a folder.
func (o *updateFolderOptions) WithName(name string) {
	o.name = NullStringFrom(name)
}

// WithParentID updates the parentID of a folder.
func (o *updateFolderOptions) WithParentID(parentID string) {
	o.parentID = NullStringFrom(parentID)
}

// WithoutParentID removes the parent ID of a folder.
func (o *updateFolderOptions) WithoutParentID() {
	o.parentID = NullStringFromPtr(nil)
}

// UpdateFolder updates a folder in the bookmarks database.
func UpdateFolder(id string, options updateFolderOptions) (Book, error) {
	return queryWithTransaction(options.db, connectDB, func(tx transaction) (Book, error) {
		var book Book

		if err := validateParentID(tx, NullStringFrom(id)); err != nil {
			return book, fmt.Errorf("bookmark ID validation failed while updating folder: %w", err)
		}

		if !options.name.Dirty && !options.parentID.Dirty {
			return book, ErrNoUpdate
		}

		if options.name.Dirty {
			if err := validateName(options.name); err != nil {
				return book, fmt.Errorf("name validation failed while updating folder: %w", err)
			}
		}

		if options.parentID.Dirty {
			if err := validateParentID(tx, options.parentID); err != nil {
				return book, fmt.Errorf("parent ID validation failed while updating folder: %w", err)
			}
		}

		if err := updateFolderDB(tx, id, updateFolderDBArgs{
			name:     options.name,
			parentID: options.parentID,
		}); err != nil {
			return book, fmt.Errorf("error while updating folder: %w", err)
		}

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:       id,
			includeFolders: true,
		})
		if err != nil {
			return book, fmt.Errorf("error geting bookmarks while updating folder: %w", err)
		}

		return books[0], nil
	})
}

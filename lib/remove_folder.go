package lib

import (
	"fmt"

	"github.com/samber/lo"
)

// removeFolderOptions are the optional arguments for RemoveFolder.
type removeFolderOptions struct {
	db NullString
}

// DefaultRemoveFolderOptions are the default options for RemoveFolder.
func DefaultRemoveFolderOptions() removeFolderOptions {
	return removeFolderOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *removeFolderOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// RemoveFolder removes a folder from the bookmarks database.
func RemoveFolder(id string, options removeFolderOptions) error {
	return execWithTransaction(options.db, connectDB, func(tx transaction) error {
		if err := validateParentID(tx, NullStringFrom(id)); err != nil {
			return fmt.Errorf("parent ID validation failed while removing folder: %w", err)
		}

		bookOrFolders, err := getParentAndChildren(tx, id)
		if err != nil {
			return fmt.Errorf("error getting folder and children while removing folder: %w", err)
		}

		for _, bookOrFolder := range lo.Reverse(bookOrFolders) {
			if !bookOrFolder.IsFolder {
				if err = unlinkTagsDB(tx, bookOrFolder.ID, bookOrFolder.Tags); err != nil {
					return fmt.Errorf("error unlinking tags while removing folder: %w", err)
				}

				if err = removeBookDB(tx, bookOrFolder.ID); err != nil {
					return fmt.Errorf("error remmoving bookmark while removing folder: %w", err)
				}

				if err = cleanOrphanedTagsDB(tx, bookOrFolder.Tags); err != nil {
					return fmt.Errorf("error cleaning orphaned tags while removing folder: %w", err)
				}
			} else {
				if err = removeFolderDB(tx, bookOrFolder.ID); err != nil {
					return fmt.Errorf("error while removing folder: %w", err)
				}
			}
		}

		return nil
	})
}

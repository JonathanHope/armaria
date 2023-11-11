package lib

import (
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
			return err
		}

		bookOrFolders, err := getParentAndChildren(tx, id)
		if err != nil {
			return err
		}

		for _, bookOrFolder := range lo.Reverse(bookOrFolders) {
			if !bookOrFolder.IsFolder {
				if err = unlinkTagsDB(tx, bookOrFolder.ID, bookOrFolder.Tags); err != nil {
					return err
				}

				if err = removeBookDB(tx, bookOrFolder.ID); err != nil {
					return err
				}

				if err = cleanOrphanedTagsDB(tx, bookOrFolder.Tags); err != nil {
					return err
				}
			} else {
				if err = removeFolderDB(tx, bookOrFolder.ID); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

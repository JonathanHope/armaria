package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
	"github.com/jonathanhope/armaria/pkg/model"
	"github.com/samber/lo"
)

// removeFolderOptions are the optional arguments for RemoveFolder.
type removeFolderOptions struct {
	DB null.NullString
}

// DefaultRemoveFolderOptions are the default options for RemoveFolder.
func DefaultRemoveFolderOptions() *removeFolderOptions {
	return &removeFolderOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *removeFolderOptions) WithDB(db string) *removeFolderOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// RemoveFolder removes a folder from the bookmarks database.
func RemoveFolder(id string, options *removeFolderOptions) error {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return fmt.Errorf("error getting config while removing folder: %w", err)
	}

	return db.ExecWithTransaction(options.DB, config.DB, func(tx db.Transaction) error {
		if err := validate.ParentID(tx, null.NullStringFrom(id)); err != nil {
			return fmt.Errorf("parent ID validation failed while removing folder: %w", err)
		}

		bookOrFolders, err := db.GetParentAndChildren(tx, id)
		if err != nil {
			return fmt.Errorf("error getting folder and children while removing folder: %w", err)
		}

		for _, bookOrFolder := range lo.Reverse(bookOrFolders) {
			if !bookOrFolder.IsFolder {
				if err = db.UnlinkTags(tx, bookOrFolder.ID, bookOrFolder.Tags); err != nil {
					return fmt.Errorf("error unlinking tags while removing folder: %w", err)
				}

				if err = db.RemoveBook(tx, bookOrFolder.ID); err != nil {
					return fmt.Errorf("error remmoving bookmark while removing folder: %w", err)
				}

				if err = db.CleanOrphanedTags(tx, bookOrFolder.Tags); err != nil {
					return fmt.Errorf("error cleaning orphaned tags while removing folder: %w", err)
				}
			} else {
				if err = db.RemoveFolder(tx, bookOrFolder.ID); err != nil {
					return fmt.Errorf("error while removing folder: %w", err)
				}
			}
		}

		return nil
	})
}

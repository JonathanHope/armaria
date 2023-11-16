package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
)

// updateFolderOptions are the optional arguments for UpdateFolder.
type updateFolderOptions struct {
	DB       null.NullString
	Name     null.NullString
	ParentID null.NullString
}

// DefaultUpdateFolderOptions are the default options for UpdateFolder.
func DefaultUpdateFolderOptions() updateFolderOptions {
	return updateFolderOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *updateFolderOptions) WithDB(db string) {
	o.DB = null.NullStringFrom(db)
}

// WithName updates the name of a folder.
func (o *updateFolderOptions) WithName(name string) {
	o.Name = null.NullStringFrom(name)
}

// WithParentID updates the parentID of a folder.
func (o *updateFolderOptions) WithParentID(parentID string) {
	o.ParentID = null.NullStringFrom(parentID)
}

// WithoutParentID removes the parent ID of a folder.
func (o *updateFolderOptions) WithoutParentID() {
	o.ParentID = null.NullStringFromPtr(nil)
}

// UpdateFolder updates a folder in the bookmarks database.
func UpdateFolder(id string, options updateFolderOptions) (armaria.Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return armaria.Book{}, fmt.Errorf("error getting config while updating folder: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (armaria.Book, error) {
		if err := validate.ParentID(tx, null.NullStringFrom(id)); err != nil {
			return armaria.Book{}, fmt.Errorf("bookmark ID validation failed while updating folder: %w", err)
		}

		if !options.Name.Dirty && !options.ParentID.Dirty {
			return armaria.Book{}, armaria.ErrNoUpdate
		}

		if options.Name.Dirty {
			if err := validate.Name(options.Name); err != nil {
				return armaria.Book{}, fmt.Errorf("name validation failed while updating folder: %w", err)
			}
		}

		if options.ParentID.Dirty {
			if err := validate.ParentID(tx, options.ParentID); err != nil {
				return armaria.Book{}, fmt.Errorf("parent ID validation failed while updating folder: %w", err)
			}
		}

		if err := db.UpdateFolder(tx, id, db.UpdateFolderArgs{
			Name:     options.Name,
			ParentID: options.ParentID,
		}); err != nil {
			return armaria.Book{}, fmt.Errorf("error while updating folder: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:       id,
			IncludeFolders: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error geting bookmarks while updating folder: %w", err)
		}

		return books[0], nil
	})
}

package armaria

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/order"
)

// updateFolderOptions are the optional arguments for UpdateFolder.
type updateFolderOptions struct {
	DB           null.NullString
	Name         null.NullString
	ParentID     null.NullString
	PreviousBook null.NullString
	NextBook     null.NullString
}

// DefaultUpdateFolderOptions are the default options for UpdateFolder.
func DefaultUpdateFolderOptions() *updateFolderOptions {
	return &updateFolderOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *updateFolderOptions) WithDB(db string) *updateFolderOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// WithName updates the name of a folder.
func (o *updateFolderOptions) WithName(name string) *updateFolderOptions {
	o.Name = null.NullStringFrom(name)
	return o
}

// WithParentID updates the parentID of a folder.
func (o *updateFolderOptions) WithParentID(parentID string) *updateFolderOptions {
	o.ParentID = null.NullStringFrom(parentID)
	return o
}

// WithoutParentID removes the parent ID of a folder.
func (o *updateFolderOptions) WithoutParentID() *updateFolderOptions {
	o.ParentID = null.NullStringFromPtr(nil)
	return o
}

// WithOrderBefore moves the bookmark to be before the provided book.
func (o *updateFolderOptions) WithOrderBefore(id string) *updateFolderOptions {
	o.NextBook = null.NullStringFrom(id)
	return o
}

// WithOrderAfter moves the bookmark to be after the provided book.
func (o *updateFolderOptions) WithOrderAfter(id string) *updateFolderOptions {
	o.PreviousBook = null.NullStringFrom(id)
	return o
}

// UpdateFolder updates a folder in the bookmarks database.
func UpdateFolder(id string, options *updateFolderOptions) (Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return Book{}, fmt.Errorf("error getting config while updating folder: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (Book, error) {
		if err := validateParentID(tx, null.NullStringFrom(id)); err != nil {
			return Book{}, fmt.Errorf("bookmark ID validation failed while updating folder: %w", err)
		}

		if !options.Name.Dirty && !options.ParentID.Dirty && !options.PreviousBook.Dirty && !options.NextBook.Dirty {
			return Book{}, ErrNoUpdate
		}

		if options.Name.Dirty {
			if err := validateName(options.Name); err != nil {
				return Book{}, fmt.Errorf("name validation failed while updating folder: %w", err)
			}
		}

		if options.ParentID.Dirty {
			if err := validateParentID(tx, options.ParentID); err != nil {
				return Book{}, fmt.Errorf("parent ID validation failed while updating folder: %w", err)
			}
		}

		current, err := validateOrdering(tx, options.PreviousBook, options.NextBook)
		if err != nil {
			return Book{}, fmt.Errorf("ordering validation failed while updating bookmark: %w", err)
		}

		if current == "" && options.ParentID.Dirty {
			previous, err := db.MaxOrder(tx, options.ParentID)
			if err != nil {
				return Book{}, fmt.Errorf("error getting max order while adding bookmark: %w", err)
			}

			if previous == "" {
				current, err = order.Initial()
				if err != nil {
					return Book{}, fmt.Errorf("error getting current order while adding bookmark: %w", err)
				}
			} else {
				current, err = order.End(previous)
				if err != nil {
					return Book{}, fmt.Errorf("error getting current order while adding bookmark: %w", err)
				}
			}
		}

		if err := db.UpdateFolder(tx, id, db.UpdateFolderArgs{
			Name:     options.Name,
			ParentID: options.ParentID,
			Order:    current,
		}); err != nil {
			return Book{}, fmt.Errorf("error while updating folder: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:       id,
			IncludeFolders: true,
		})
		if err != nil {
			return Book{}, fmt.Errorf("error geting bookmarks while updating folder: %w", err)
		}

		return toBook(books[0]), nil
	})
}

package armaria

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/order"
)

// addFolderOptions are the optional arguments for AddFolder.
type addFolderOptions struct {
	DB       null.NullString
	ParentID null.NullString
}

// DefaultAddFolderOptions are the default options for AddFolder.
func DefaultAddFolderOptions() *addFolderOptions {
	return &addFolderOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addFolderOptions) WithDB(db string) *addFolderOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// WithParentID sets the folders' parent ID.
func (o *addFolderOptions) WithParentID(parentID string) *addFolderOptions {
	o.ParentID = null.NullStringFrom(parentID)
	return o
}

// AddFolder adds a folder to the bookmarks database.
func AddFolder(name string, options *addFolderOptions) (Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return Book{}, fmt.Errorf("error getting config while adding folder: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (Book, error) {
		if err := validateName(null.NullStringFrom(name)); err != nil {
			return Book{}, fmt.Errorf("name validation failed while adding folder: %w", err)
		}

		if err := validateParentID(tx, options.ParentID); err != nil {
			return Book{}, fmt.Errorf("parent ID validation failed while adding folder: %w", err)
		}

		previous, err := db.MaxOrder(tx, options.ParentID)
		if err != nil {
			return Book{}, fmt.Errorf("error getting max order while adding bookmark: %w", err)
		}

		var current string
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

		id, err := db.AddFolder(tx, name, options.ParentID, current)
		if err != nil {
			return Book{}, fmt.Errorf("error while adding folder: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:       id,
			IncludeFolders: true,
		})
		if err != nil {
			return Book{}, fmt.Errorf("error getting folders while adding folder: %w", err)
		}

		return toBook(books[0]), nil
	})
}

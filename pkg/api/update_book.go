package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
	"github.com/jonathanhope/armaria/pkg/model"
)

// updateBookOptions are the optional arguments for UpdateBook.
type updateBookOptions struct {
	DB          null.NullString
	Name        null.NullString
	URL         null.NullString
	Description null.NullString
	ParentID    null.NullString
}

// DefaultUpdateBookOptions are the default options for UpdateBook.
func DefaultUpdateBookOptions() *updateBookOptions {
	return &updateBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *updateBookOptions) WithDB(db string) *updateBookOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// WithName updates the name of a bookmark.
func (o *updateBookOptions) WithName(name string) *updateBookOptions {
	o.Name = null.NullStringFrom(name)
	return o
}

// WithURL updates the URL of a bookmark.
func (o *updateBookOptions) WithURL(url string) *updateBookOptions {
	o.URL = null.NullStringFrom(url)
	return o
}

// WithDescription updates the description of a bookmark.
func (o *updateBookOptions) WithDescription(description string) *updateBookOptions {
	o.Description = null.NullStringFrom(description)
	return o
}

// WithParentID updates the parentID of a bookmark.
func (o *updateBookOptions) WithParentID(parentID string) *updateBookOptions {
	o.ParentID = null.NullStringFrom(parentID)
	return o
}

// WithoutDescription removes the description of a bookmark.
func (o *updateBookOptions) WithoutDescription() *updateBookOptions {
	o.Description = null.NullStringFromPtr(nil)
	return o
}

// WithoutParentID removes the parent ID of a bookmark.
func (o *updateBookOptions) WithoutParentID() *updateBookOptions {
	o.ParentID = null.NullStringFromPtr(nil)
	return o
}

// UpdateBook updates a bookmark in the bookmarks database.
func UpdateBook(id string, options *updateBookOptions) (armaria.Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return armaria.Book{}, fmt.Errorf("error getting config while updating bookmark: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (armaria.Book, error) {
		if err := validate.BookID(tx, id); err != nil {
			return armaria.Book{}, fmt.Errorf("bookmark ID validation failed while updating bookmark: %w", err)
		}

		if !options.Name.Dirty && !options.URL.Dirty && !options.Description.Dirty && !options.ParentID.Dirty {
			return armaria.Book{}, armaria.ErrNoUpdate
		}

		if options.Name.Dirty {
			if err := validate.Name(options.Name); err != nil {
				return armaria.Book{}, fmt.Errorf("name validation failed while updating bookmark: %w", err)
			}
		}

		if options.URL.Dirty {
			if err := validate.URL(options.URL); err != nil {
				return armaria.Book{}, fmt.Errorf("URL validation failed while updating bookmark: %w", err)
			}
		}

		if options.Description.Dirty {
			if err := validate.Description(options.Description); err != nil {
				return armaria.Book{}, fmt.Errorf("description validation failed while updating bookmark: %w", err)
			}
		}

		if options.ParentID.Dirty {
			if err := validate.ParentID(tx, options.ParentID); err != nil {
				return armaria.Book{}, fmt.Errorf("parent ID validation failed while updating bookmark: %w", err)
			}
		}

		if err := db.UpdateBook(tx, id, db.UpdateBookArgs{
			Name:        options.Name,
			URL:         options.URL,
			Description: options.Description,
			ParentID:    options.ParentID,
		}); err != nil {
			return armaria.Book{}, fmt.Errorf("error while updating bookmark: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting bookmarks while updating bookmark: %w", err)
		}

		return books[0], nil
	})
}

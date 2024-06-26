package armaria

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/order"
)

// updateBookOptions are the optional arguments for UpdateBook.
type updateBookOptions struct {
	DB           null.NullString
	Name         null.NullString
	URL          null.NullString
	Description  null.NullString
	ParentID     null.NullString
	PreviousBook null.NullString
	NextBook     null.NullString
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

// WithOrderBefore moves the bookmark to be before the provided book.
func (o *updateBookOptions) WithOrderBefore(id string) *updateBookOptions {
	o.NextBook = null.NullStringFrom(id)
	return o
}

// WithOrderAfter moves the bookmark to be after the provided book.
func (o *updateBookOptions) WithOrderAfter(id string) *updateBookOptions {
	o.PreviousBook = null.NullStringFrom(id)
	return o
}

// UpdateBook updates a bookmark in the bookmarks database.
func UpdateBook(id string, options *updateBookOptions) (Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return Book{}, fmt.Errorf("error getting config while updating bookmark: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (Book, error) {
		if err := validateBookID(tx, id); err != nil {
			return Book{}, fmt.Errorf("bookmark ID validation failed while updating bookmark: %w", err)
		}

		if !options.Name.Dirty && !options.URL.Dirty && !options.Description.Dirty && !options.ParentID.Dirty && !options.PreviousBook.Dirty && !options.NextBook.Dirty {
			return Book{}, ErrNoUpdate
		}

		if options.Name.Dirty {
			if err := validateName(options.Name); err != nil {
				return Book{}, fmt.Errorf("name validation failed while updating bookmark: %w", err)
			}
		}

		if options.URL.Dirty {
			if err := validateURL(options.URL); err != nil {
				return Book{}, fmt.Errorf("URL validation failed while updating bookmark: %w", err)
			}
		}

		if options.Description.Dirty {
			if err := validateDescription(options.Description); err != nil {
				return Book{}, fmt.Errorf("description validation failed while updating bookmark: %w", err)
			}
		}

		if options.ParentID.Dirty {
			if err := validateParentID(tx, options.ParentID); err != nil {
				return Book{}, fmt.Errorf("parent ID validation failed while updating bookmark: %w", err)
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

		if err := db.UpdateBook(tx, id, db.UpdateBookArgs{
			Name:        options.Name,
			URL:         options.URL,
			Description: options.Description,
			ParentID:    options.ParentID,
			Order:       current,
		}); err != nil {
			return Book{}, fmt.Errorf("error while updating bookmark: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return Book{}, fmt.Errorf("error getting bookmarks while updating bookmark: %w", err)
		}

		return toBook(books[0]), nil
	})
}

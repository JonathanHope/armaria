package armaria

import (
	"fmt"

	"errors"
	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/order"
	"github.com/samber/lo"
)

// addBookOptions are the optional arguments for AddBook.
type addBookOptions struct {
	DB          null.NullString
	Name        null.NullString
	Description null.NullString
	ParentID    null.NullString
	Tags        []string
}

// DefaultAddBookOptions are the default options for AddBook.
func DefaultAddBookOptions() *addBookOptions {
	return &addBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addBookOptions) WithDB(db string) *addBookOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// WithName sets the bookmark's name.
func (o *addBookOptions) WithName(name string) *addBookOptions {
	o.Name = null.NullStringFrom(name)
	return o
}

// WithDescription sets the bookmark's description.
func (o *addBookOptions) WithDescription(description string) *addBookOptions {
	o.Description = null.NullStringFrom(description)
	return o
}

// WithParentID sets the bookmark's parent ID.
func (o *addBookOptions) WithParentID(parentID string) *addBookOptions {
	o.ParentID = null.NullStringFrom(parentID)
	return o
}

// WithTags sets the bookmark's tags.
func (o *addBookOptions) WithTags(tags []string) *addBookOptions {
	o.Tags = tags
	return o
}

// AddBook adds a bookmark to the bookmarks database.
func AddBook(url string, options *addBookOptions) (Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return Book{}, fmt.Errorf("error getting config while adding bookmark: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (Book, error) {
		// Default name to URL if not provided.
		if !options.Name.Valid {
			options.Name = null.NullStringFrom(url)
		}

		if err := validateURL(null.NullStringFrom(url)); err != nil {
			return Book{}, fmt.Errorf("URL validation failed while adding bookmark: %w", err)
		}

		if err := validateName(options.Name); err != nil {
			return Book{}, fmt.Errorf("name validation failed while adding bookmark: %w", err)
		}

		if err := validateDescription(options.Description); err != nil {
			return Book{}, fmt.Errorf("description validation failed while adding bookmark: %w", err)
		}

		if err := validateParentID(tx, options.ParentID); err != nil {
			return Book{}, fmt.Errorf("parent ID validation failed while adding bookmark: %w", err)
		}

		if err := validateTags(options.Tags, make([]string, 0)); err != nil {
			return Book{}, fmt.Errorf("tags validation failed while adding bookmark: %w", err)
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

		id, err := db.AddBook(tx, url, options.Name.String, options.Description, options.ParentID, current)
		if err != nil {
			return Book{}, fmt.Errorf("error while adding bookmark: %w", err)
		}

		existingTags, err := db.GetTags(tx, db.GetTagsArgs{
			TagsFilter: options.Tags,
		})
		if err != nil {
			return Book{}, fmt.Errorf("error getting tags while adding bookmark: %w", err)
		}

		tagsToAdd, _ := lo.Difference(options.Tags, existingTags)
		if err = db.AddTags(tx, tagsToAdd); err != nil {
			return Book{}, fmt.Errorf("error adding tags while adding bookmark: %w", err)
		}

		if err = db.LinkTags(tx, id, options.Tags); err != nil {
			return Book{}, fmt.Errorf("error linking tags while adding bookmark: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return Book{}, fmt.Errorf("error getting bookmarks while adding bookmark: %w", err)
		}

		return toBook(books[0]), nil
	})
}

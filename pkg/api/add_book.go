package armariaapi

import (
	"fmt"

	"errors"
	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
	"github.com/jonathanhope/armaria/pkg/model"
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
func DefaultAddBookOptions() addBookOptions {
	return addBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addBookOptions) WithDB(db string) {
	o.DB = null.NullStringFrom(db)
}

// WithName sets the bookmark's name.
func (o *addBookOptions) WithName(name string) {
	o.Name = null.NullStringFrom(name)
}

// WithDescription sets the bookmark's description.
func (o *addBookOptions) WithDescription(description string) {
	o.Description = null.NullStringFrom(description)
}

// WithParentID sets the bookmark's parent ID.
func (o *addBookOptions) WithParentID(parentID string) {
	o.ParentID = null.NullStringFrom(parentID)
}

// WithTags sets the bookmark's tags.
func (o *addBookOptions) WithTags(tags []string) {
	o.Tags = tags
}

// AddBook adds a bookmark to the bookmarks database.
func AddBook(url string, options addBookOptions) (armaria.Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return armaria.Book{}, fmt.Errorf("error getting config while adding bookmark: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (armaria.Book, error) {
		// Default name to URL if not provided.
		if !options.Name.Valid {
			options.Name = null.NullStringFrom(url)
		}

		if err := validate.URL(null.NullStringFrom(url)); err != nil {
			return armaria.Book{}, fmt.Errorf("URL validation failed while adding bookmark: %w", err)
		}

		if err := validate.Name(options.Name); err != nil {
			return armaria.Book{}, fmt.Errorf("name validation failed while adding bookmark: %w", err)
		}

		if err := validate.Description(options.Description); err != nil {
			return armaria.Book{}, fmt.Errorf("description validation failed while adding bookmark: %w", err)
		}

		if err := validate.ParentID(tx, options.ParentID); err != nil {
			return armaria.Book{}, fmt.Errorf("parent ID validation failed while adding bookmark: %w", err)
		}

		if err := validate.Tags(options.Tags, make([]string, 0)); err != nil {
			return armaria.Book{}, fmt.Errorf("tags validation failed while adding bookmark: %w", err)
		}

		id, err := db.AddBook(tx, url, options.Name.String, options.Description, options.ParentID)
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error while adding bookmark: %w", err)
		}

		existingTags, err := db.GetTags(tx, db.GetTagsArgs{
			TagsFilter: options.Tags,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting tags while adding bookmark: %w", err)
		}

		tagsToAdd, _ := lo.Difference(options.Tags, existingTags)
		if err = db.AddTags(tx, tagsToAdd); err != nil {
			return armaria.Book{}, fmt.Errorf("error adding tags while adding bookmark: %w", err)
		}

		if err = db.LinkTags(tx, id, options.Tags); err != nil {
			return armaria.Book{}, fmt.Errorf("error linking tags while adding bookmark: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting bookmarks while adding bookmark: %w", err)
		}

		return books[0], nil
	})
}

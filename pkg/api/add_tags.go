package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
	"github.com/samber/lo"
)

// addTagsOptions are the optional arguments for AddTags.
type addTagsOptions struct {
	DB null.NullString
}

// DefaultAddTagsOptions are the default options for AddTags.
func DefaultAddTagsOptions() addTagsOptions {
	return addTagsOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addTagsOptions) WithDB(db string) {
	o.DB = null.NullStringFrom(db)
}

// AddTags adds tags to a bookmark in the bookmarks database.
func AddTags(id string, tags []string, options addTagsOptions) (armaria.Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return armaria.Book{}, fmt.Errorf("error getting config while adding tag: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (armaria.Book, error) {
		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting tags while adding tags: %w", err)
		}

		if len(books) != 1 || books[0].IsFolder {
			return armaria.Book{}, armaria.ErrBookNotFound
		}

		if err := validate.Tags(tags, books[0].Tags); err != nil {
			return armaria.Book{}, fmt.Errorf("tags validation failed while adding tags: %w", err)
		}

		existingTags, err := db.GetTags(tx, db.GetTagsArgs{
			TagsFilter: tags,
		})
		if err != nil {
			return armaria.Book{}, err
		}

		tagsToAdd, _ := lo.Difference(tags, existingTags)
		if err = db.AddTags(tx, tagsToAdd); err != nil {
			return armaria.Book{}, fmt.Errorf("error while adding tags: %w", err)
		}

		if err = db.LinkTags(tx, id, tags); err != nil {
			return armaria.Book{}, fmt.Errorf("error linking tags while adding tags: %w", err)
		}

		books, err = db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting bookmarks while adding tags: %w", err)
		}

		return books[0], nil
	})
}

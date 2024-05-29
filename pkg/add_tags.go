package armaria

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/samber/lo"
)

// addTagsOptions are the optional arguments for AddTags.
type addTagsOptions struct {
	DB null.NullString
}

// DefaultAddTagsOptions are the default options for AddTags.
func DefaultAddTagsOptions() *addTagsOptions {
	return &addTagsOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addTagsOptions) WithDB(db string) *addTagsOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// AddTags adds tags to a bookmark in the bookmarks database.
func AddTags(id string, tags []string, options *addTagsOptions) (Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return Book{}, fmt.Errorf("error getting config while adding tag: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (Book, error) {
		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return Book{}, fmt.Errorf("error getting tags while adding tags: %w", err)
		}

		if len(books) != 1 || books[0].IsFolder {
			return Book{}, ErrBookNotFound
		}

		book := toBook(books[0])

		if err := validateTags(tags, book.Tags); err != nil {
			return Book{}, fmt.Errorf("tags validation failed while adding tags: %w", err)
		}

		existingTags, err := db.GetTags(tx, db.GetTagsArgs{
			TagsFilter: tags,
		})
		if err != nil {
			return Book{}, err
		}

		tagsToAdd, _ := lo.Difference(tags, existingTags)
		if err = db.AddTags(tx, tagsToAdd); err != nil {
			return Book{}, fmt.Errorf("error while adding tags: %w", err)
		}

		if err = db.LinkTags(tx, id, tags); err != nil {
			return Book{}, fmt.Errorf("error linking tags while adding tags: %w", err)
		}

		books, err = db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return Book{}, fmt.Errorf("error getting bookmarks while adding tags: %w", err)
		}

		return toBook(books[0]), nil
	})
}

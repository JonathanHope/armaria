package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg/model"
	"github.com/samber/lo"
)

// removeTagsOptions are the optional arguments for RemoveTags.
type removeTagsOptions struct {
	DB null.NullString
}

// DefaultRemoveTagsOptions are the default options for RemoveTags.
func DefaultRemoveTagsOptions() removeTagsOptions {
	return removeTagsOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *removeTagsOptions) WithDB(db string) {
	o.DB = null.NullStringFrom(db)
}

// RemoveTags removes tags from a bookmark in the bookmarks database.
func RemoveTags(id string, tags []string, options removeTagsOptions) (armaria.Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return armaria.Book{}, fmt.Errorf("error getting config while removing tags: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) (armaria.Book, error) {
		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting bookmarks while removing tags: %w", err)
		}

		if len(books) != 1 || books[0].IsFolder {
			return armaria.Book{}, armaria.ErrBookNotFound
		}

		for _, tag := range tags {
			if !lo.Contains(books[0].Tags, tag) {
				return armaria.Book{}, armaria.ErrTagNotFound
			}
		}

		if err = db.UnlinkTags(tx, books[0].ID, tags); err != nil {
			return armaria.Book{}, fmt.Errorf("error unlinking tags while removing tags: %w", err)
		}

		if err = db.CleanOrphanedTags(tx, tags); err != nil {
			return armaria.Book{}, fmt.Errorf("error cleaning orphaned tags while removing tags: %w", err)
		}

		books, err = db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:     id,
			IncludeBooks: true,
		})
		if err != nil {
			return armaria.Book{}, fmt.Errorf("error getting bookmarks while removing tags: %w", err)
		}

		return books[0], nil
	})
}

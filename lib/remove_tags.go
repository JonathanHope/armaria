package lib

import (
	"fmt"

	"github.com/samber/lo"
)

// removeTagsOptions are the optional arguments for RemoveTags.
type removeTagsOptions struct {
	db NullString
}

// DefaultRemoveTagsOptions are the default options for RemoveTags.
func DefaultRemoveTagsOptions() removeTagsOptions {
	return removeTagsOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *removeTagsOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// RemoveTags removes tags from a bookmark in the bookmarks database.
func RemoveTags(id string, tags []string, options removeTagsOptions) (Book, error) {
	return queryWithTransaction(options.db, connectDB, func(tx transaction) (Book, error) {
		var book Book

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return book, fmt.Errorf("error getting bookmarks while removing tags: %w", err)
		}

		if len(books) != 1 || books[0].IsFolder {
			return book, ErrBookNotFound
		}

		for _, tag := range tags {
			if !lo.Contains(books[0].Tags, tag) {
				return book, ErrTagNotFound
			}
		}

		if err = unlinkTagsDB(tx, books[0].ID, tags); err != nil {
			return book, fmt.Errorf("error unlinking tags while removing tags: %w", err)
		}

		if err = cleanOrphanedTagsDB(tx, tags); err != nil {
			return book, fmt.Errorf("error cleaning orphaned tags while removing tags: %w", err)
		}

		books, err = getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return book, fmt.Errorf("error getting bookmarks while removing tags: %w", err)
		}

		return books[0], nil
	})
}

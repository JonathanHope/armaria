package lib

import (
	"github.com/samber/lo"
)

// addTagsOptions are the optional arguments for AddTags.
type addTagsOptions struct {
	db NullString
}

// DefaultAddTagsOptions are the default options for AddTags.
func DefaultAddTagsOptions() addTagsOptions {
	return addTagsOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addTagsOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// AddTags adds tags to a bookmark in the bookmarks database.
func AddTags(id string, tags []string, options addTagsOptions) (Book, error) {
	return queryWithTransaction(options.db, connectDB, func(tx transaction) (Book, error) {
		var book Book

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return book, err
		}

		if len(books) != 1 || books[0].IsFolder {
			return book, ErrBookNotFound
		}

		if err := validateTags(tags, books[0].Tags); err != nil {
			return book, err
		}

		existingTags, err := getTagsDB(tx, getTagsDBArgs{
			tagsFilter: tags,
		})
		if err != nil {
			return book, err
		}

		tagsToAdd, _ := lo.Difference(tags, existingTags)
		if err = addTagsDB(tx, tagsToAdd); err != nil {
			return book, err
		}

		if err = linkTagsDB(tx, id, tags); err != nil {
			return book, err
		}

		books, err = getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return book, err
		}

		return books[0], nil
	})
}

package lib

import (
	"fmt"

	"github.com/samber/lo"
)

// addBookOptions are the optional arguments for AddBook.
type addBookOptions struct {
	db          NullString
	name        NullString
	description NullString
	parentID    NullString
	tags        []string
}

// DefaultAddBookOptions are the default options for AddBook.
func DefaultAddBookOptions() addBookOptions {
	return addBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *addBookOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// WithName sets the bookmark's name.
func (o *addBookOptions) WithName(name string) {
	o.name = NullStringFrom(name)
}

// WithDescription sets the bookmark's description.
func (o *addBookOptions) WithDescription(description string) {
	o.description = NullStringFrom(description)
}

// WithParentID sets the bookmark's parent ID.
func (o *addBookOptions) WithParentID(parentID string) {
	o.parentID = NullStringFrom(parentID)
}

// WithTags sets the bookmark's tags.
func (o *addBookOptions) WithTags(tags []string) {
	o.tags = tags
}

// AddBook adds a bookmark to the bookmarks database.
func AddBook(url string, options addBookOptions) (Book, error) {
	return queryWithTransaction(options.db, connectDB, func(tx transaction) (Book, error) {
		var book Book

		// Default name to URL if not provided.
		if !options.name.Valid {
			options.name = NullStringFrom(url)
		}

		if err := validateURL(NullStringFrom(url)); err != nil {
			return book, fmt.Errorf("URL validation failed while adding bookmark: %w", err)
		}

		if err := validateName(options.name); err != nil {
			return book, fmt.Errorf("name validation failed while adding bookmark: %w", err)
		}

		if err := validateDescription(options.description); err != nil {
			return book, fmt.Errorf("description validation failed while adding bookmark: %w", err)
		}

		if err := validateParentID(tx, options.parentID); err != nil {
			return book, fmt.Errorf("parent ID validation failed while adding bookmark: %w", err)
		}

		if err := validateTags(options.tags, make([]string, 0)); err != nil {
			return book, fmt.Errorf("tags validation failed while adding bookmark: %w", err)
		}

		id, err := addBookDB(tx, url, options.name.String, options.description, options.parentID)
		if err != nil {
			return book, fmt.Errorf("error while adding bookmark: %w", err)
		}

		existingTags, err := getTagsDB(tx, getTagsDBArgs{
			tagsFilter: options.tags,
		})
		if err != nil {
			return book, fmt.Errorf("error getting tags while adding bookmark: %w", err)
		}

		tagsToAdd, _ := lo.Difference(options.tags, existingTags)
		if err = addTagsDB(tx, tagsToAdd); err != nil {
			return book, fmt.Errorf("error adding tags while adding bookmark: %w", err)
		}

		if err = linkTagsDB(tx, id, options.tags); err != nil {
			return book, fmt.Errorf("error linking tags while adding bookmark: %w", err)
		}

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return book, fmt.Errorf("error getting bookmarks while adding bookmark: %w", err)
		}

		return books[0], nil
	})
}

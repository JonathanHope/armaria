package lib

import (
	"fmt"
)

// listBooksOptions are the optional arguments for ListBooks.
type listBooksOptions struct {
	db               NullString
	includeBookmarks bool
	includeFolders   bool
	parentID         NullString
	query            NullString
	tags             []string
	after            NullString
	order            Order
	direction        Direction
	first            NullInt64
}

// DefaultListBooksOptions are the default options for ListBooks.
func DefaultListBooksOptions() listBooksOptions {
	return listBooksOptions{
		includeBookmarks: true,
		includeFolders:   true,
		order:            OrderModified,
		direction:        DirectionAsc,
	}
}

// WithDB sets the location of the bookmarks database.
func (o *listBooksOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// WithBooks sets whether to include bookmark results.
func (o *listBooksOptions) WithBooks(include bool) {
	o.includeBookmarks = include
}

// WithIncludeFolder sets whether to include folder results.
func (o *listBooksOptions) WithFolders(include bool) {
	o.includeFolders = include
}

// WithParentID filters by parent ID.
func (o *listBooksOptions) WithParentID(parentID string) {
	o.parentID = NullStringFrom(parentID)
}

// WithQuery searches on name, URL, description.
func (o *listBooksOptions) WithQuery(query string) {
	o.query = NullStringFrom(query)
}

// WithTags filters by tag.
func (o *listBooksOptions) WithTags(tags []string) {
	o.tags = tags
}

// WithAfter returns results after an ID.
func (o *listBooksOptions) WithAfter(after string) {
	o.after = NullStringFrom(after)
}

// WithOrder sets the column to order on.
func (o *listBooksOptions) WithOrder(order Order) {
	o.order = order
}

// WithDirection sets the direction to order by.
func (o *listBooksOptions) WithDirection(direction Direction) {
	o.direction = direction
}

// withFirst sets the max number of results to return.
func (o *listBooksOptions) WithFirst(first int64) {
	o.first = NullInt64From(first)
}

// WithoutParentID removes the parent ID of a bookmark.
func (o *listBooksOptions) WithoutParentID() {
	o.parentID = NullStringFromPtr(nil)
}

// ListBooks lists bookmarks and folders in the bookmarks database.
func ListBooks(options listBooksOptions) ([]Book, error) {
	return queryWithDB(options.db, connectDB, func(tx transaction) ([]Book, error) {
		books := make([]Book, 0)

		if !options.includeBookmarks && !options.includeFolders {
			return books, nil
		}

		if err := validateFirst(options.first); err != nil {
			return books, fmt.Errorf("first validation failed while listing bookmarks: %w", err)
		}

		if err := validateDirection(options.direction); err != nil {
			return books, fmt.Errorf("direction validation failed while listing bookmarks: %w", err)
		}

		if err := validateOrder(options.order); err != nil {
			return books, fmt.Errorf("order validation failed while listing bookmarks: %w", err)
		}

		if err := validateQuery(options.query); err != nil {
			return books, fmt.Errorf("query validation failed while listing bookmarks: %w", err)
		}

		books, err := getBooksDB(tx, getBooksDBArgs{
			includeBooks:   options.includeBookmarks,
			includeFolders: options.includeFolders,
			parentID:       options.parentID,
			query:          options.query,
			tags:           options.tags,
			after:          options.after,
			order:          options.order,
			direction:      options.direction,
			first:          options.first,
		})
		if err != nil {
			return books, fmt.Errorf("error while listing bookmarks: %w", err)
		}

		return books, nil
	})
}

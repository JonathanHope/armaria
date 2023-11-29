package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
	"github.com/jonathanhope/armaria/pkg/model"
)

// listBookOptions are the optional arguments for ListBooks.
type listBookOptions struct {
	DB               null.NullString
	IncludeBookmarks bool
	IncludeFolders   bool
	ParentID         null.NullString
	Query            null.NullString
	Tags             []string
	After            null.NullString
	Order            armaria.Order
	Direction        armaria.Direction
	First            null.NullInt64
}

// DefaultListBooksOptions are the default options for ListBooks.
func DefaultListBooksOptions() listBookOptions {
	return listBookOptions{
		IncludeBookmarks: true,
		IncludeFolders:   true,
		Order:            armaria.OrderModified,
		Direction:        armaria.DirectionAsc,
	}
}

// WithDB sets the location of the bookmarks database.
func (o *listBookOptions) WithDB(db string) {
	o.DB = null.NullStringFrom(db)
}

// WithBooks sets whether to include bookmark results.
func (o *listBookOptions) WithBooks(include bool) {
	o.IncludeBookmarks = include
}

// WithIncludeFolder sets whether to include folder results.
func (o *listBookOptions) WithFolders(include bool) {
	o.IncludeFolders = include
}

// WithParentID filters by parent ID.
func (o *listBookOptions) WithParentID(parentID string) {
	o.ParentID = null.NullStringFrom(parentID)
}

// WithQuery searches on name, URL, description.
func (o *listBookOptions) WithQuery(query string) {
	o.Query = null.NullStringFrom(query)
}

// WithTags filters by tag.
func (o *listBookOptions) WithTags(tags []string) {
	o.Tags = tags
}

// WithAfter returns results after an ID.
func (o *listBookOptions) WithAfter(after string) {
	o.After = null.NullStringFrom(after)
}

// WithOrder sets the column to order on.
func (o *listBookOptions) WithOrder(order armaria.Order) {
	o.Order = order
}

// WithDirection sets the direction to order by.
func (o *listBookOptions) WithDirection(direction armaria.Direction) {
	o.Direction = direction
}

// withFirst sets the max number of results to return.
func (o *listBookOptions) WithFirst(first int64) {
	o.First = null.NullInt64From(first)
}

// WithoutParentID removes the parent ID of a bookmark.
func (o *listBookOptions) WithoutParentID() {
	o.ParentID = null.NullStringFromPtr(nil)
}

// ListBooks lists bookmarks and folders in the bookmarks database.
func ListBooks(options listBookOptions) ([]armaria.Book, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return nil, fmt.Errorf("error getting config while listing bookmarks: %w", err)
	}

	return db.QueryWithDB(options.DB, config.DB, func(tx db.Transaction) ([]armaria.Book, error) {
		if !options.IncludeBookmarks && !options.IncludeFolders {
			return nil, nil
		}

		if err := validate.First(options.First); err != nil {
			return nil, fmt.Errorf("first validation failed while listing bookmarks: %w", err)
		}

		if err := validate.Direction(options.Direction); err != nil {
			return nil, fmt.Errorf("direction validation failed while listing bookmarks: %w", err)
		}

		if err := validate.Order(options.Order); err != nil {
			return nil, fmt.Errorf("order validation failed while listing bookmarks: %w", err)
		}

		if err := validate.Query(options.Query); err != nil {
			return nil, fmt.Errorf("query validation failed while listing bookmarks: %w", err)
		}

		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IncludeBooks:   options.IncludeBookmarks,
			IncludeFolders: options.IncludeFolders,
			ParentID:       options.ParentID,
			Query:          options.Query,
			Tags:           options.Tags,
			After:          options.After,
			Order:          options.Order,
			Direction:      options.Direction,
			First:          options.First,
		})
		if err != nil {
			return books, fmt.Errorf("error while listing bookmarks: %w", err)
		}

		return books, nil
	})
}

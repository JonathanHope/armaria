package armariaapi

import (
	"errors"
	"fmt"

	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/validate"
)

// listTagsOptions are the optional arguments for ListTags.
type listTagsOptions struct {
	DB        null.NullString
	Query     null.NullString
	After     null.NullString
	Direction armaria.Direction
	First     null.NullInt64
}

// DefaultListTagsOptions are the default options for ListTags.
func DefaultListTagsOptions() listTagsOptions {
	return listTagsOptions{
		Direction: armaria.DirectionAsc,
	}
}

// WithDB sets the location of the bookmarks database.
func (o *listTagsOptions) WithDB(db string) {
	o.DB = null.NullStringFrom(db)
}

// WithQuery searches on tags.
func (o *listTagsOptions) WithQuery(query string) {
	o.Query = null.NullStringFrom(query)
}

// WithAfter returns results after an ID.
func (o *listTagsOptions) WithAfter(after string) {
	o.After = null.NullStringFrom(after)
}

// WithDirection sets the direction to order by.
func (o *listTagsOptions) WithDirection(direction armaria.Direction) {
	o.Direction = direction
}

// withFirst sets the max number of results to return.
func (o *listTagsOptions) WithFirst(first int64) {
	o.First = null.NullInt64From(first)
}

// ListTags lists tags in the bookmarks database.
func ListTags(options listTagsOptions) ([]string, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		return nil, fmt.Errorf("error getting config while listing tags: %w", err)
	}

	return db.QueryWithDB(options.DB, config.DB, func(tx db.Transaction) ([]string, error) {
		tags := make([]string, 0)

		if err := validate.First(options.First); err != nil {
			return tags, fmt.Errorf("first validation failed while listing tags: %w", err)
		}

		if err := validate.Direction(options.Direction); err != nil {
			return tags, fmt.Errorf("direction validation failed while listing tags: %w", err)
		}

		if err := validate.Query(options.Query); err != nil {
			return tags, fmt.Errorf("query validation failed while listing tags: %w", err)
		}

		tags, err := db.GetTags(tx, db.GetTagsArgs{
			Query:     options.Query,
			After:     options.After,
			Direction: options.Direction,
			First:     options.First,
		})
		if err != nil {
			return tags, fmt.Errorf("error while listing tags: %w", err)
		}

		return tags, nil
	})
}

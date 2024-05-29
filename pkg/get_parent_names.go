package armaria

import (
	"fmt"

	"errors"
	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/samber/lo"
)

// getParentNameOptions are the optional arguments for GetParentNames.
type getParentNameOptions struct {
	DB null.NullString
}

// DefaultGetParentNameOptions are the default options for GetParentNames.
func DefaultGetParentNameOptions() *getParentNameOptions {
	return &getParentNameOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *getParentNameOptions) WithDB(db string) *getParentNameOptions {
	o.DB = null.NullStringFrom(db)
	return o
}

// GetParentNames gets the parent names of a bookmark or folder.
func GetParentNames(ID string, options *getParentNameOptions) ([]string, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return nil, fmt.Errorf("error getting config while getting parent names: %w", err)
	}

	return db.QueryWithTransaction(options.DB, config.DB, func(tx db.Transaction) ([]string, error) {
		names, err := db.GetBookFolderParents(tx, ID)
		if err != nil {
			return nil, fmt.Errorf("error getting parent: %w", err)
		}

		if len(names) == 0 {
			return nil, ErrNotFound
		}

		names = lo.Reverse(names)
		return names, nil
	})
}

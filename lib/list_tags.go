package lib

// listTagsOptions are the optional arguments for ListTags.
type listTagsOptions struct {
	db        NullString
	query     NullString
	after     NullString
	direction Direction
	first     NullInt64
}

// DefaultListTagsOptions are the default options for ListTags.
func DefaultListTagsOptions() listTagsOptions {
	return listTagsOptions{
		direction: DirectionAsc,
	}
}

// WithDB sets the location of the bookmarks database.
func (o *listTagsOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// WithQuery searches on tags.
func (o *listTagsOptions) WithQuery(query string) {
	o.query = NullStringFrom(query)
}

// WithAfter returns results after an ID.
func (o *listTagsOptions) WithAfter(after string) {
	o.after = NullStringFrom(after)
}

// WithDirection sets the direction to order by.
func (o *listTagsOptions) WithDirection(direction Direction) {
	o.direction = direction
}

// withFirst sets the max number of results to return.
func (o *listTagsOptions) WithFirst(first int64) {
	o.first = NullInt64From(first)
}

// ListTags lists tags in the bookmarks database.
func ListTags(options listTagsOptions) ([]string, error) {
	return queryWithDB(options.db, connectDB, func(tx transaction) ([]string, error) {
		tags := make([]string, 0)

		if err := validateFirst(options.first); err != nil {
			return tags, err
		}

		if err := validateDirection(options.direction); err != nil {
			return tags, err
		}

		if err := validateQuery(options.query); err != nil {
			return tags, err
		}

		return getTagsDB(tx, getTagsDBArgs{
			query:     options.query,
			after:     options.after,
			direction: options.direction,
			first:     options.first,
		})
	})
}

package lib

// updateBookOptions are the optional arguments for UpdateBook.
type updateBookOptions struct {
	db          NullString
	name        NullString
	url         NullString
	description NullString
	parentID    NullString
}

// DefaultUpdateBookOptions are the default options for UpdateBook.
func DefaultUpdateBookOptions() updateBookOptions {
	return updateBookOptions{}
}

// WithDB sets the location of the bookmarks database.
func (o *updateBookOptions) WithDB(db string) {
	o.db = NullStringFrom(db)
}

// WithName updates the name of a bookmark.
func (o *updateBookOptions) WithName(name string) {
	o.name = NullStringFrom(name)
}

// WithURL updates the URL of a bookmark.
func (o *updateBookOptions) WithURL(url string) {
	o.url = NullStringFrom(url)
}

// WithDescription updates the description of a bookmark.
func (o *updateBookOptions) WithDescription(description string) {
	o.description = NullStringFrom(description)
}

// WithParentID updates the parentID of a bookmark.
func (o *updateBookOptions) WithParentID(parentID string) {
	o.parentID = NullStringFrom(parentID)
}

// WithoutDescription removes the description of a bookmark.
func (o *updateBookOptions) WithoutDescription() {
	o.description = NullStringFromPtr(nil)
}

// WithoutParentID removes the parent ID of a bookmark.
func (o *updateBookOptions) WithoutParentID() {
	o.parentID = NullStringFromPtr(nil)
}

// UpdateBook updates a bookmark in the bookmarks database.
func UpdateBook(id string, options updateBookOptions) (Book, error) {
	return queryWithTransaction(options.db, connectDB, func(tx transaction) (Book, error) {
		var book Book

		if err := validateBookID(tx, id); err != nil {
			return book, err
		}

		if !options.name.Dirty && !options.url.Dirty && !options.description.Dirty && !options.parentID.Dirty {
			return book, ErrNoUpdate
		}

		if options.name.Dirty {
			if err := validateName(options.name); err != nil {
				return book, err
			}
		}

		if options.url.Dirty {
			if err := validateURL(options.url); err != nil {
				return book, err
			}
		}

		if options.description.Dirty {
			if err := validateDescription(options.description); err != nil {
				return book, err
			}
		}

		if options.parentID.Dirty {
			if err := validateParentID(tx, options.parentID); err != nil {
				return book, err
			}
		}

		if err := updateBookDB(tx, id, updateBookDBArgs{
			name:        options.name,
			url:         options.url,
			description: options.description,
			parentID:    options.parentID,
		}); err != nil {
			return book, err
		}

		books, err := getBooksDB(tx, getBooksDBArgs{
			idFilter:     id,
			includeBooks: true,
		})
		if err != nil {
			return book, err
		}

		return books[0], nil
	})
}

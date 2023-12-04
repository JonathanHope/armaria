package validate

import (
	"fmt"
	"regexp"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg/model"
	"github.com/samber/lo"
)

// This file contains validators that should be run on user inputs before any DB calls are made.

const MaxNameLength = 2048
const MaxURLLength = 2048
const MaxDescriptionLength = 4096
const MaxTagLength = 128
const MaxTags = 24
const MinQueryLength = 3

// URL validates a URL value.
// It's only used for bookmarks and must have a length >= 1 and <= 2048.
func URL(URL null.NullString) error {
	if !URL.Valid || URL.String == "" {
		return armaria.ErrURLTooShort
	}

	if len(URL.String) > MaxURLLength {
		return armaria.ErrURLTooLong
	}

	return nil
}

// Name validates a name value.
// It's required and must have a length >= 1 and <= 1024.
func Name(name null.NullString) error {
	if !name.Valid || name.String == "" {
		return armaria.ErrNameTooShort
	}

	if len(name.String) > MaxNameLength {
		return armaria.ErrNameTooLong
	}

	return nil
}

// Description validates a description value.
// IT's optional and must have a length <= 5096.
func Description(description null.NullString) error {
	if !description.Valid {
		return nil
	}

	if description.String == "" {
		return armaria.ErrDescriptionTooShort
	}

	if len(description.String) > MaxDescriptionLength {
		return armaria.ErrDescriptionTooLong
	}

	return nil
}

// Tags validates a tags value.
// They must be unique.
// Each must have a length >= 1 and <= 128.
// Each must be comprised of the chars A-Z a-z 0-9 - _
func Tags(tags []string, existingTags []string) error {
	if len(tags)+len(existingTags) > 24 {
		return armaria.ErrTooManyTags
	}

	if len(tags) != len(lo.Uniq(tags)) {
		return armaria.ErrDuplicateTag
	}

	r, err := regexp.Compile(`^[a-zA-Z0-9\-_]*$`)
	if err != nil {
		return fmt.Errorf("error compiling regex while validating tags: %w", err)
	}

	for _, tag := range tags {
		if tag == "" {
			return armaria.ErrTagTooShort
		}

		if len(tag) > MaxTagLength {
			return armaria.ErrTagTooLong
		}

		matched := r.MatchString(tag)
		if !matched {
			return armaria.ErrTagInvalidChar
		}

		if lo.Contains(existingTags, tag) {
			return armaria.ErrDuplicateTag
		}
	}

	return nil
}

// First validates a limit value.
// It's optional, but if it is provided it must be > 0.
func First(limit null.NullInt64) error {
	if !limit.Valid {
		return nil
	}

	if limit.Int64 <= 0 {
		return armaria.ErrFirstTooSmall
	}

	return nil
}

// Order validates an order value.
// It must be modified or name.
func Order(order armaria.Order) error {
	if order != armaria.OrderModified && order != armaria.OrderName {
		return armaria.ErrInvalidOrder
	}

	return nil
}

// Direction validates a direction value.
// It must be asc or desc.
func Direction(direction armaria.Direction) error {
	if direction != armaria.DirectionAsc && direction != armaria.DirectionDesc {
		return armaria.ErrInvalidDirection
	}

	return nil

}

// Query validates a query value.
// It's optional must be at least 3 chars long.
func Query(query null.NullString) error {
	if !query.Valid {
		return nil
	}

	if len(query.String) < MinQueryLength {
		return armaria.ErrQueryTooShort
	}

	return nil
}

// BookID validates a bookmark ID value.
// The target bookmark must exist.
func BookID(tx db.Transaction, ID string) error {
	exists, err := db.BookFolderExists(tx, ID, false)
	if err != nil {
		return fmt.Errorf("error checking if bookmark exists while validating bookmark ID")
	}

	if !exists {
		return armaria.ErrBookNotFound
	}

	return nil
}

// ParentID validates a parent ID value.
// It's optional but if it is provided the target parent folder must exist.
func ParentID(tx db.Transaction, parentID null.NullString) error {
	if !parentID.Valid {
		return nil
	}

	exists, err := db.BookFolderExists(tx, parentID.String, true)
	if err != nil {
		return fmt.Errorf("error checking if folder exists while validating parent ID")
	}

	if !exists {
		return armaria.ErrFolderNotFound
	}

	return nil
}

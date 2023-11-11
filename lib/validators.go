package lib

import (
	"fmt"
	"regexp"

	"github.com/samber/lo"
)

// This file contains validators that should be run on user inputs before any DB calls are made.

const MaxNameLength = 2048
const MaxURLLength = 2048
const MaxDescriptionLength = 4096
const MaxTagLength = 128
const MaxTags = 24
const MinQueryLength = 3

// validateURL validates a URL value.
// URL is only used for bookmarks and must have a length >= 1 and <= 2048.
func validateURL(URL NullString) error {
	if !URL.Valid || URL.String == "" {
		return ErrURLTooShort
	}

	if len(URL.String) > MaxURLLength {
		return ErrURLTooLong
	}

	return nil
}

// validateName validates a name value.
// Description is required and must have a length >= 1 and <= 1024.
func validateName(name NullString) error {
	if !name.Valid || name.String == "" {
		return ErrNameTooShort
	}

	if len(name.String) > MaxNameLength {
		return ErrNameTooLong
	}

	return nil
}

// validateDescription validates a description value.
// Description is optional and must have a length <= 5096.
func validateDescription(description NullString) error {
	if !description.Valid {
		return nil
	}

	if description.String == "" {
		return ErrDescriptionTooShort
	}

	if len(description.String) > MaxDescriptionLength {
		return ErrDescriptionTooLong
	}

	return nil
}

// validateTag validates a tags value.
// Tags must be unique.
// Each must have a length >= 1 and <= 128.
// Tags can have the chars A-Z a-z 0-9 - _
func validateTags(tags []string, existingTags []string) error {
	if len(tags)+len(existingTags) > 24 {
		return ErrTooManyTags
	}

	if len(tags) != len(lo.Uniq(tags)) {
		return ErrDuplicateTag
	}

	r, e := regexp.Compile(`^[a-zA-Z0-9\-_]*$`)
	if e != nil {
		return e
	}

	for _, tag := range tags {
		if tag == "" {
			return ErrTagTooShort
		}

		if len(tag) > MaxTagLength {
			return ErrTagTooLong
		}

		matched := r.MatchString(tag)
		if !matched {
			return ErrTagInvalidChar
		}

		if lo.Contains(existingTags, tag) {
			return ErrDuplicateTag
		}
	}

	return nil
}

// validateFirst validates a limit value.
// Limit is optional, but if it is provided it must be > 0.
func validateFirst(limit NullInt64) error {
	if !limit.Valid {
		return nil
	}

	if limit.Int64 <= 0 {
		return ErrFirstTooSmall
	}

	return nil
}

// validateOrder validates an order value.
// Order must be modified or name.
func validateOrder(order Order) error {
	if order != OrderModified && order != OrderName {
		return ErrInvalidOrder
	}

	return nil
}

// validateDirection validates a direction value.
// Direction must be asc or desc.
func validateDirection(direction Direction) error {
	if direction != DirectionAsc && direction != DirectionDesc {
		return ErrInvalidDirection
	}

	return nil

}

// validateQuery validates a query value.
// Query is optional must be at least 3 chars long.
func validateQuery(query NullString) error {
	if !query.Valid {
		return nil
	}

	if len(query.String) < MinQueryLength {
		return ErrQueryTooShort
	}

	return nil
}

// validateBookID validates a bookmark ID value.
// The target bookmark must exist.
func validateBookID(tx transaction, ID string) error {
	exists, err := bookFolderExistsDB(tx, ID, false)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	if !exists {
		return ErrBookNotFound
	}

	return nil
}

// validateParentID validates a parent ID value.
// Parent ID is optional but if it is provided the target parent folder must exist.
func validateParentID(tx transaction, parentID NullString) error {
	if !parentID.Valid {
		return nil
	}

	exists, err := bookFolderExistsDB(tx, parentID.String, true)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	if !exists {
		return ErrFolderNotFound
	}

	return nil
}

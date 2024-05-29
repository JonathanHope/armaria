package armaria

import (
	"fmt"
	"regexp"

	"github.com/jonathanhope/armaria/internal/db"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/order"
	"github.com/samber/lo"
)

const maxNameLength = 2048
const maxURLLength = 2048
const maxDescriptionLength = 4096
const maxTagLength = 128
const maxTags = 24
const minQueryLength = 3

// validateURL validates a validateURL value.
// It's only used for bookmarks and must have a length >= 1 and <= 2048.
func validateURL(URL null.NullString) error {
	if !URL.Valid || URL.String == "" {
		return ErrURLTooShort
	}

	if len(URL.String) > maxURLLength {
		return ErrURLTooLong
	}

	return nil
}

// validateName validates a name value.
// It's required and must have a length >= 1 and <= 1024.
func validateName(name null.NullString) error {
	if !name.Valid || name.String == "" {
		return ErrNameTooShort
	}

	if len(name.String) > maxNameLength {
		return ErrNameTooLong
	}

	return nil
}

// validateDescription validates a description value.
// IT's optional and must have a length <= 5096.
func validateDescription(description null.NullString) error {
	if !description.Valid {
		return nil
	}

	if description.String == "" {
		return ErrDescriptionTooShort
	}

	if len(description.String) > maxDescriptionLength {
		return ErrDescriptionTooLong
	}

	return nil
}

// validateTags validates a tags value.
// They must be unique.
// Each must have a length >= 1 and <= 128.
// Each must be comprised of the chars A-Z a-z 0-9 - _
func validateTags(tags []string, existingTags []string) error {
	if len(tags)+len(existingTags) > 24 {
		return ErrTooManyTags
	}

	if len(tags) != len(lo.Uniq(tags)) {
		return ErrDuplicateTag
	}

	r, err := regexp.Compile(`^[a-zA-Z0-9\-_]*$`)
	if err != nil {
		return fmt.Errorf("error compiling regex while validating tags: %w", err)
	}

	for _, tag := range tags {
		if tag == "" {
			return ErrTagTooShort
		}

		if len(tag) > maxTagLength {
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
// It's optional, but if it is provided it must be > 0.
func validateFirst(limit null.NullInt64) error {
	if !limit.Valid {
		return nil
	}

	if limit.Int64 <= 0 {
		return ErrFirstTooSmall
	}

	return nil
}

// validateOrder validates an order value.
// It must be modified or name.
func validateOrder(order Order) error {
	if order != OrderModified && order != OrderName && order != OrderManual {
		return ErrInvalidOrder
	}

	return nil
}

// validateDirection validates a direction value.
// It must be asc or desc.
func validateDirection(direction Direction) error {
	if direction != DirectionAsc && direction != DirectionDesc {
		return ErrInvalidDirection
	}

	return nil

}

// validateQuery validates a query value.
// It's optional must be at least 3 chars long.
func validateQuery(query null.NullString) error {
	if !query.Valid {
		return nil
	}

	if len(query.String) < minQueryLength {
		return ErrQueryTooShort
	}

	return nil
}

// validateBookID validates a bookmark ID value.
// The target bookmark must exist.
func validateBookID(tx db.Transaction, ID string) error {
	exists, err := db.BookFolderExists(tx, ID, false)
	if err != nil {
		return fmt.Errorf("error checking if bookmark exists while validating bookmark ID")
	}

	if !exists {
		return ErrBookNotFound
	}

	return nil
}

// validateParentID validates a parent ID value.
// It's optional but if it is provided the target parent folder must exist.
func validateParentID(tx db.Transaction, parentID null.NullString) error {
	if !parentID.Valid {
		return nil
	}

	exists, err := db.BookFolderExists(tx, parentID.String, true)
	if err != nil {
		return fmt.Errorf("error checking if folder exists while validating parent ID")
	}

	if !exists {
		return ErrFolderNotFound
	}

	return nil
}

// validateOrdering validates the values used for manual ordering.
// The previousBook value must be < the nextBook value.
func validateOrdering(tx db.Transaction, previousID null.NullString, nextID null.NullString) (string, error) {
	var previousOrder string
	var nextOrder string

	if previousID.Dirty {
		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:       previousID.String,
			IncludeBooks:   true,
			IncludeFolders: true,
		})
		if err != nil {
			return "", fmt.Errorf("error checking if book exists while validating ordering")
		}
		if len(books) == 0 {
			return "", ErrNotFound
		}
		previousOrder = books[0].Order
	}

	if nextID.Dirty {
		books, err := db.GetBooks(tx, db.GetBooksArgs{
			IDFilter:       nextID.String,
			IncludeBooks:   true,
			IncludeFolders: true,
		})
		if err != nil {
			return "", fmt.Errorf("error checking if book exists while validating ordering")
		}
		if len(books) == 0 {
			return "", ErrNotFound
		}
		nextOrder = books[0].Order
	}

	if previousOrder != "" && nextOrder != "" && previousOrder >= nextOrder {
		return "", ErrInvalidOrdering
	}

	if previousOrder != "" && nextOrder != "" {
		return order.Between(previousOrder, nextOrder)
	} else if previousOrder != "" {
		return order.End(previousOrder)
	} else if nextOrder != "" {
		return order.Start(nextOrder)
	}

	return "", nil
}

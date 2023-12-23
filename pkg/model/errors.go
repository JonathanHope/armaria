package armaria

import (
	"errors"
)

var (
	// ErrNoUpdate is returned when an update is requested with no updates.
	ErrNoUpdate = errors.New("no update")
	// ErrBookNotFound is returned when a target bookmark was not found.
	ErrBookNotFound = errors.New("bookmark not found")
	// ErrFolderNotFound is returned when a target folder was not found.
	ErrFolderNotFound = errors.New("folder not found")
	// ErrTagNotFound is returned when a target tag was not found.
	ErrTagNotFound = errors.New("tag not found")
	// ErrNotFound is returned when a target bookmark or folder was not found.
	ErrNotFound = errors.New("bookmark or folder not found")
	// ErrURLTooShort is returned when a provided URL is too short.
	ErrURLTooShort = errors.New("URL too short")
	// ErrURLTooLong is too long when a provided URL is too long.
	ErrURLTooLong = errors.New("URL too long")
	// ErrNameTooShort is returned when a provided name is too short.
	ErrNameTooShort = errors.New("name too short")
	// ErrNameTooLong is returned when a provided name is too long.
	ErrNameTooLong = errors.New("name too long")
	// ErrDescriptionTooShort is returned when a provided description is too short.
	ErrDescriptionTooShort = errors.New("description too short")
	// ErrDescriptionTooLong is returned when a provided description is too long.
	ErrDescriptionTooLong = errors.New("description too long")
	// ErrTagTooShort is returned when a provided tag is too short.
	ErrTagTooShort = errors.New("tag too short")
	// ErrTagTooLong is returned when a provided tag is too long.
	ErrTagTooLong = errors.New("tag too long")
	// ErrDuplicateTag is returned when a tag is applied twice to a bookmark.
	ErrDuplicateTag = errors.New("tags must be unique")
	// ErrTooManyTags is returned when too many tags have been applied to bookmark.
	ErrTooManyTags = errors.New("too many tags")
	// ErrTagInvalidChar is returned when a provided tag has an invalid character.
	ErrTagInvalidChar = errors.New("tag had invalid chars")
	// ErrFirstTooSmall is returned when a provided first is too small.
	ErrFirstTooSmall = errors.New("first too small")
	// ErrInvalidOrder is returned when a provided order is invalid.
	ErrInvalidOrder = errors.New("invalid order")
	// ErrInvalidDirection is returned when a provided direction is invalid.
	ErrInvalidDirection = errors.New("invalid direction")
	// ErrQueryTooShort is returned when a provided query is too short.
	ErrQueryTooShort = errors.New("query too short")
	// ErrConfigMissing is returned when the config file is missing.
	ErrConfigMissing = errors.New("config missing")
	// ErrInvalidOrdering is returned previous book >= next book for manual ordering.
	ErrInvalidOrdering = errors.New("invalid ordering")
)

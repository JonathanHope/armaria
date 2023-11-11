package cmd

import (
	"errors"
)

var (
	// ErrFolderNoFolderMutuallyExclusive is returned if folder and no-folder are both provided.
	ErrFolderNoFolderMutuallyExclusive = errors.New("folder/no-folder mutually exclusive")
	// ErrDescriptionNoDescriptionMutuallyExclusive is returned if no-description and description are both provoded.
	ErrDescriptionNoDescriptionMutuallyExclusive = errors.New("description/no-description mutually exclusive")
)

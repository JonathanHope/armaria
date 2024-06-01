package messaging

import (
	"encoding/json"

	"github.com/jonathanhope/armaria/internal/null"
)

// Each message kind has its own payload.
// These are all of the different payloads.
// There are also helpers to marshall/unmarshall the payloads.

// Payload are the payloads that a message can have.
type Payload interface {
	AddBookPayload | AddFolderPayload | AddTagsPayload | ListBooksPayload | ListTagsPayload | RemoveBookPayload | RemoveFolderPayload | RemoveTagsPayload | UpdateBookPayload | UpdateFolderPayload | ErrorPayload | BooksPayload | BookPayload | TagsPayload | VoidPayload | ConfigValuePayload | ParentNamesPayload
}

// AddBookPayload is a payload for a request to add a bookmark.
type AddBookPayload struct {
	DB          null.NullString `json:"db"`
	URL         string          `json:"url"`
	Name        null.NullString `json:"name"`
	Description null.NullString `json:"description"`
	ParentID    null.NullString `json:"parentId"`
	Tags        []string        `json:"tags"`
}

// AddFolderPayload is a payload for a request to add a folder.
type AddFolderPayload struct {
	DB       null.NullString `json:"db"`
	Name     string          `json:"name"`
	ParentID null.NullString `json:"parentId"`
}

// AddTagsPayload is a payload for a request to add tags.
type AddTagsPayload struct {
	DB   null.NullString `json:"db"`
	ID   string          `json:"id"`
	Tags []string        `json:"tags"`
}

// ListBooksPayload is a payload for a request to list bookmarks.
type ListBooksPayload struct {
	DB               null.NullString `json:"db"`
	IncludeBookmarks bool            `json:"includeBookmarks"`
	IncludeFolders   bool            `json:"includeFolders"`
	ParentID         null.NullString `json:"parentID"`
	WithoutParentID  bool            `json:"withoutParentID"`
	Query            null.NullString `json:"query"`
	Tags             []string        `json:"tags"`
	After            null.NullString `json:"after"`
	Order            string          `json:"order"`
	Direction        string          `json:"direction"`
	First            null.NullInt64  `json:"first"`
}

// ListTagsPayload is a payload for a request to list tags.
type ListTagsPayload struct {
	DB        null.NullString `json:"db"`
	Query     null.NullString `json:"query"`
	After     null.NullString `json:"after"`
	Direction string          `json:"direction"`
	First     null.NullInt64  `json:"first"`
}

// RemoveBookPayload is a payload for a request to delete a bookmark.
type RemoveBookPayload struct {
	DB null.NullString `json:"db"`
	ID string          `json:"id"`
}

// RemoveFolderPayload is a payload for a request to delete a folder.
type RemoveFolderPayload struct {
	DB null.NullString `json:"db"`
	ID string          `json:"id"`
}

// RemoveTagsPayload is a payload for a request to remove tags.
type RemoveTagsPayload struct {
	DB   null.NullString `json:"db"`
	ID   string          `json:"id"`
	Tags []string        `json:"tags"`
}

// UpdateBookPayload is a payload for a request to update a bookmark.
type UpdateBookPayload struct {
	DB                null.NullString `json:"db"`
	ID                string          `json:"id"`
	Name              null.NullString `json:"name"`
	URL               null.NullString `json:"url"`
	Description       null.NullString `json:"description"`
	ParentID          null.NullString `json:"parentId"`
	RemoveDescription bool            `json:"removeDescription"`
	RemoveParentID    bool            `json:"removeParentID"`
	PreviousBook      null.NullString `json:"previousBook"`
	NextBook          null.NullString `json:"nextBook"`
}

// UpdateBookPayload is a payload for a request to update a folder.
type UpdateFolderPayload struct {
	DB             null.NullString `json:"db"`
	ID             string          `json:"id"`
	Name           null.NullString `json:"name"`
	ParentID       null.NullString `json:"parentId"`
	RemoveParentID bool            `json:"removeParentID"`
	PreviousBook   null.NullString `json:"previousBook"`
	NextBook       null.NullString `json:"nextBook"`
}

// ErrorPayload is a payload for a response with an error in it.
type ErrorPayload struct {
	Error string `json:"error"`
}

// BooksPayload is a payload for a response with bookmarks/folders in it.
type BooksPayload struct {
	Books []BookDTO `json:"books"`
}

// BookPayload is a payload for a response with a bookmark/folder in it.
type BookPayload struct {
	Book BookDTO `json:"book"`
}

// TagsPayload is a payload for a response with tags in it.
type TagsPayload struct {
	Tags []string `json:"tags"`
}

// ConfigValuePayload is a payload for a response with a config value in it.
type ConfigValuePayload struct {
	Value string `json:"value"`
}

// ParentNamesPayload is a payload for a response with a books parents names in it.
type ParentNamesPayload struct {
	ParentNames []string
}

// VoidPayload is a payload for a response with nothing in it.
type VoidPayload struct{}

// GetPayload gets the underlying payload in a message.
func GetPayload[T Payload](msg NativeMessage) (T, error) {
	var payload T
	err := json.Unmarshal([]byte(msg.Payload), &payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

// PayloadToMessage converts a payload to a NativeMessage.
func PayloadToMessage[T Payload](kind MessageKind, payload T) (NativeMessage, error) {
	json, err := json.Marshal(payload)
	if err != nil {
		return NativeMessage{}, err
	}

	return NativeMessage{
		Kind:    kind,
		Payload: string(json),
	}, nil
}

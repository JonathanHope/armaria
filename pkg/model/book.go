package armaria

// Book is a bookmark or folder.
type Book struct {
	ID          string   // unique identifier of a bookmark/folder
	URL         *string  // address of a bookmark; not used for folders
	Name        string   // name of a bookmark/folder
	Description *string  // description of a bookmark/folder
	ParentID    *string  // optional ID of the parent folder for a bookmark/folder
	IsFolder    bool     // true if folder, and false otherwise
	ParentName  *string  // name of parent folder if bookmark/folder has one
	Tags        []string // tags applied to the bookmark
}

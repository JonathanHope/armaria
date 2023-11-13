package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/jonathanhope/armaria/lib"
)

// RootCmd is the top level CLI command for Armaria.
type RootCmd struct {
	DB        *string   `help:"Location of the bookmarks database."`
	Formatter Formatter `help:"How to format output: pretty/json." enum:"json,pretty" default:"pretty"`

	Add    AddCmd    `cmd:"" help:"Add a folder, bookmark, or tag."`
	Remove RemoveCmd `cmd:"" help:"Remove a folder, bookmark, or tag."`
	Update UpdateCmd `cmd:"" help:"Update a folder or bookmark."`
	List   ListCmd   `cmd:"" help:"List folders, bookmarks, or tags."`

	Config ConfigCmd `cmd:"" help:"Manage the configuration."`
}

// AddCmd is a CLI command to add a bookmark or folder.
type AddCmd struct {
	Book   AddBookCmd   `cmd:"" help:"Add a bookmark."`
	Folder AddFolderCmd `cmd:"" help:"Add a folder."`
	Tag    AddTagsCmd   `cmd:"" help:"Add tags to a bookmark."`
}

// ListCmd is a CLI command to list bookmarks/folders/tags.
type ListCmd struct {
	All     ListAllCmd     `cmd:"" help:"List bookmarks and folders."`
	Books   ListBooksCmd   `cmd:"" help:"List bookmarks."`
	Folders ListFoldersCmd `cmd:"" help:"List folders."`
	Tags    ListTagsCmd    `cmd:"" help:"List tags."`
}

// UpdateCmd is a CLI command to update a bookmark or folder.
type UpdateCmd struct {
	Book   UpdateBookCmd   `cmd:"" help:"Update a bookmark."`
	Folder UpdateFolderCmd `cmd:"" help:"Update a folder."`
}

// RemoveCmd is a CLI command to remove a folder or bookmark.
type RemoveCmd struct {
	Book   RemoveBookCmd   `cmd:"" help:"Remove a bookmark."`
	Folder RemoveFolderCmd `cmd:"" help:"Remove a folder."`
	Tag    RemoveTagsCmd   `cmd:"" help:"Remove tags from a bookmark."`
}

// AddBookCmd is a CLI command to add a bookmark.
type AddBookCmd struct {
	Folder      *string  `help:"Folder to add this bookmark to."`
	Name        *string  `help:"Name for the bookmark."`
	Description *string  `help:"Description of the bookmark."`
	Tag         []string `help:"Tag to apply to the bookmark."`

	URL string `arg:"" name:"url" help:"URL of the bookmark."`
}

// ConfigCmd is a CLI command to manage config.
type ConfigCmd struct {
	DB DBConfigCmd `cmd:"" help:"Manage the bookmarks database location configuration."`
}

// DBConfigCmd is a CLI command to manage the bookmarks database location config.
type DBConfigCmd struct {
	Get GetDBConfigCmd `cmd:"" help:"Get the location of the bookmarks database from the configuration."`
	Set SetDBConfigCmd `cmd:"" help:"Set the location of the bookmarks database in the configuration."`
}

// Run add a bookmark.
func (r *AddBookCmd) Run(ctx *Context) error {
	start := time.Now()

	options := lib.DefaultAddBookOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}
	if r.Name != nil {
		options.WithName(*r.Name)
	}
	if r.Description != nil {
		options.WithDescription(*r.Description)
	}
	if r.Tag != nil {
		options.WithTags(r.Tag)
	}

	book, err := lib.AddBook(r.URL, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []lib.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Added in %s", elapsed))

	return nil
}

// AddFolderCmd is a CLI command to add a folder.
type AddFolderCmd struct {
	Folder *string `help:"Folder to add this folder to."`

	Name string `arg:"" name:"name" help:"Name for the folder."`
}

// Run add a folder.
func (r *AddFolderCmd) Run(ctx *Context) error {
	start := time.Now()

	options := lib.DefaultAddFolderOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}

	book, err := lib.AddFolder(r.Name, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []lib.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Added in %s", elapsed))

	return nil
}

// AddTagsCmd is a CLI command to add tags to an existing bookmark.
type AddTagsCmd struct {
	Tag []string `help:"Tag to apply to the bookmark."`

	ID string `arg:"" name:"id" help:"ID of the bookmark to add tags to."`
}

// Run add tags.
func (r *AddTagsCmd) Run(ctx *Context) error {
	start := time.Now()

	options := lib.DefaultAddTagsOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	book, err := lib.AddTags(r.ID, r.Tag, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []lib.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Tagged in %s", elapsed))

	return nil
}

// ListAllCmd is a CLI command to list bookmarks and folders.
type ListAllCmd struct {
	Folder   *string       `help:"Folder to list bookmarks/folders in."`
	NoFolder bool          `help:"List top level bookmarks/folders."`
	After    *string       `help:"ID of bookmark/folder to return results after."`
	Query    *string       `help:"Query to search bookmarks/folders by."`
	Tag      []string      `help:"Tag to filter bookmarks/folders by."`
	Order    lib.Order     `help:"Field results are ordered on: modified/name." enum:"modified,name" default:"modified"`
	Dir      lib.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First    *int64        `help:"The max number of bookmarks/folders to return."`
}

// Run list bookmarks and folders.
func (r *ListAllCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := lib.DefaultListBooksOptions()
	options.WithFolders(true)
	options.WithBooks(true)
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}
	if r.NoFolder {
		options.WithoutParentID()
	}
	if r.After != nil {
		options.WithAfter(*r.After)
	}
	if r.Query != nil {
		options.WithQuery(*r.Query)
	}
	if r.Tag != nil {
		options.WithTags(r.Tag)
	}
	if r.Order != "" {
		options.WithOrder(r.Order)
	}
	if r.Dir != "" {
		options.WithDirection(r.Dir)
	}
	if r.First != nil {
		options.WithFirst(*r.First)
	}

	books, err := lib.ListBooks(options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, books)
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Listed in %s", elapsed))

	return nil
}

// ListBooksCmd is a CLI command to list bookmarks.
type ListBooksCmd struct {
	Folder   *string       `help:"Folder to list bookmarks in."`
	NoFolder bool          `help:"List top level bookmarks."`
	After    *string       `help:"ID of bookmark to return results after."`
	Query    *string       `help:"Query to search bookmarks by."`
	Tag      []string      `help:"Tag to filter bookmarks by."`
	Order    lib.Order     `help:"Field results are ordered on: modified/name." enum:"modified,name" default:"modified"`
	Dir      lib.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First    *int64        `help:"The max number of bookmarks to return."`
}

// Run list bookmarks.
func (r *ListBooksCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := lib.DefaultListBooksOptions()
	options.WithFolders(false)
	options.WithBooks(true)
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}
	if r.NoFolder {
		options.WithoutParentID()
	}
	if r.After != nil {
		options.WithAfter(*r.After)
	}
	if r.Query != nil {
		options.WithQuery(*r.Query)
	}
	if r.Tag != nil {
		options.WithTags(r.Tag)
	}
	if r.Order != "" {
		options.WithOrder(r.Order)
	}
	if r.Dir != "" {
		options.WithDirection(r.Dir)
	}
	if r.First != nil {
		options.WithFirst(*r.First)
	}

	books, err := lib.ListBooks(options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, books)
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Listed in %s", elapsed))

	return nil
}

// ListFoldersCmd is a CLI command to list folders.
type ListFoldersCmd struct {
	Folder   *string       `help:"Folder to list folders in."`
	NoFolder bool          `help:"List top level folders."`
	After    *string       `help:"ID of folder to return results after."`
	Query    *string       `help:"Query to search folders by."`
	Tag      []string      `help:"Tag to filter folders by."`
	Order    lib.Order     `help:"Field results are ordered on: modified/name." enum:"modified,name" default:"modified"`
	Dir      lib.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First    *int64        `help:"The max number of folders to return."`
}

// Run list folders.
func (r *ListFoldersCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := lib.DefaultListBooksOptions()
	options.WithFolders(true)
	options.WithBooks(false)
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}
	if r.NoFolder {
		options.WithoutParentID()
	}
	if r.After != nil {
		options.WithAfter(*r.After)
	}
	if r.Query != nil {
		options.WithQuery(*r.Query)
	}
	if r.Tag != nil {
		options.WithTags(r.Tag)
	}
	if r.Order != "" {
		options.WithOrder(r.Order)
	}
	if r.Dir != "" {
		options.WithDirection(r.Dir)
	}
	if r.First != nil {
		options.WithFirst(*r.First)
	}

	books, err := lib.ListBooks(options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, books)
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Listed in %s", elapsed))

	return nil
}

// ListTagsCmd is a CLI command to list tags.
type ListTagsCmd struct {
	Query *string       `help:"Query to search tags by."`
	After *string       `help:"ID of tags to return results after."`
	Dir   lib.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First *int64        `help:"The max number of tags to return."`
}

// Run list tags.
func (r *ListTagsCmd) Run(ctx *Context) error {
	start := time.Now()

	options := lib.DefaultListTagsOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.Query != nil {
		options.WithQuery(*r.Query)
	}
	if r.After != nil {
		options.WithAfter(*r.After)
	}
	if r.Dir != "" {
		options.WithDirection(r.Dir)
	}
	if r.First != nil {
		options.WithFirst(*r.First)
	}

	tags, err := lib.ListTags(options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatTagResults(ctx.Writer, ctx.Formatter, tags)
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Listed in %s", elapsed))

	return nil
}

// UpdateBookCmd is a CLI command to update a bookmark.
type UpdateBookCmd struct {
	Folder        *string `help:"Folder to move this bookmark to."`
	NoFolder      bool    `help:"Remove the parent folder."`
	Name          *string `help:"New name for this bookmark."`
	Description   *string `help:"New description for this bookmark."`
	NoDescription bool    `help:"Remove the description."`
	URL           *string `help:"New URL for this bookmark."`

	ID string `arg:"" name:"id" help:"ID of the bookmark to update."`
}

// Run update a bookmark.
func (r *UpdateBookCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoDescription && r.Description != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrDescriptionNoDescriptionMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := lib.DefaultUpdateBookOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.NoDescription {
		options.WithoutDescription()
	}
	if r.Description != nil {
		options.WithDescription(*r.Description)
	}
	if r.NoFolder {
		options.WithoutParentID()
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}
	if r.Name != nil {
		options.WithName(*r.Name)
	}
	if r.URL != nil {
		options.WithURL(*r.URL)
	}

	book, err := lib.UpdateBook(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []lib.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Updated in %s", elapsed))

	return nil
}

// UpdateFolderCmd is a CLI command to update a folder.
type UpdateFolderCmd struct {
	Name     *string `help:"New name for this folder."`
	Folder   *string `help:"Folder to move this folder to."`
	NoFolder bool    `help:"Remove the parent folder."`

	ID string `arg:"" name:"id" help:"ID of the folder to update."`
}

// Run update a folder.
func (r *UpdateFolderCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := lib.DefaultUpdateFolderOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.NoFolder {
		options.WithoutParentID()
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}
	if r.Name != nil {
		options.WithName(*r.Name)
	}

	book, err := lib.UpdateFolder(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []lib.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Updated in %s", elapsed))

	return nil
}

// RemoveBookCmd is a CLI command to remove a bookmark.
type RemoveBookCmd struct {
	ID string `arg:"" name:"id" help:"ID of the bookmark to remove."`
}

// Run remove a bookmark.
func (r *RemoveBookCmd) Run(ctx *Context) error {
	start := time.Now()

	options := lib.DefaultRemoveBookOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	err := lib.RemoveBook(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Removed in %s", elapsed))

	return nil
}

// RemoveFolderCmd is a CLI command to remove a folder.
type RemoveFolderCmd struct {
	ID string `arg:"" name:"id" help:"ID of the folder to remove."`
}

// Run remove a folder.
func (r *RemoveFolderCmd) Run(ctx *Context) error {
	start := time.Now()

	options := lib.DefaultRemoveFolderOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	err := lib.RemoveFolder(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Removed in %s", elapsed))

	return nil
}

// RemoveTagsCmd is a CLI command to add remove tags from an existing bookmark.
type RemoveTagsCmd struct {
	Tag []string `help:"Tag to remove from the bookmark."`

	ID string `arg:"" name:"id" help:"ID of the bookmark to remove tags from."`
}

// Run remove tags.
func (r *RemoveTagsCmd) Run(ctx *Context) error {
	start := time.Now()

	options := lib.DefaultRemoveTagsOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	book, err := lib.RemoveTags(r.ID, r.Tag, options)

	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []lib.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Untagged in %s", elapsed))

	return nil
}

// GetDBConfigCmd is a CLI command to get the location of the bookmarks database from the config.
type GetDBConfigCmd struct {
}

// Run get the location of the bookmarks database from the config.
func (r *GetDBConfigCmd) Run(ctx *Context) error {
	start := time.Now()

	config, err := lib.GetConfig()

	if err != nil && !errors.Is(err, lib.ErrConfigMissing) {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatConfigResult(ctx.Writer, ctx.Formatter, config.DB)
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Retreived in %s", elapsed))

	return nil
}

// SetDBConfigCmd is a CLI command to set the location of the bookmarks database in the config.
type SetDBConfigCmd struct {
	DB string `arg:"" name:"db" help:"Location of the bookmarks database."`
}

// Run set the location of the bookmarks database in the config.
func (r *SetDBConfigCmd) Run(ctx *Context) error {
	start := time.Now()

	err := lib.UpdateConfig(func(config *lib.Config) {
		config.DB = r.DB
	})

	if err != nil && !errors.Is(err, lib.ErrConfigMissing) {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Set in %s", elapsed))

	return nil
}

// RootCmdFactory creates a new RootCmd.
func RootCmdFactory() RootCmd {
	var cmd RootCmd
	return cmd
}

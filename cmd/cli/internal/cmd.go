package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/jonathanhope/armaria/cmd/cli/tui"
	"github.com/jonathanhope/armaria/pkg"
)

// RootCmd is the top level CLI command for Armaria.
type RootCmd struct {
	DB        *string   `help:"Location of the bookmarks database."`
	Formatter Formatter `help:"How to format output: pretty/json." enum:"json,pretty" default:"pretty"`

	Add    AddCmd    `cmd:"" help:"Add a folder, bookmark, or tag."`
	Remove RemoveCmd `cmd:"" help:"Remove a folder, bookmark, or tag."`
	Update UpdateCmd `cmd:"" help:"Update a folder or bookmark."`
	List   ListCmd   `cmd:"" help:"List folders, bookmarks, or tags."`
	Get    GetCmd    `cmd:"" help:"Get a folder or bookmark."`
	Query  QueryCmd  `cmd:"" help:"Query folders and bookmarks."`

	Config   ConfigCmd   `cmd:"" help:"Manage the configuration."`
	Manifest ManifestCmd `cmd:"" help:"Manage the app manifest."`

	TUI TUICommand `cmd:"" help:"Start the TUI."`

	Version VersionCmd `cmd:"" help:"Print the current version."`
}

// RootCmdFactory creates a new RootCmd.
func RootCmdFactory() RootCmd {
	var cmd RootCmd
	return cmd
}

// AddCmd is a CLI command to add a bookmark or folder.
type AddCmd struct {
	Book   AddBookCmd   `cmd:"" help:"Add a bookmark."`
	Folder AddFolderCmd `cmd:"" help:"Add a folder."`
	Tag    AddTagsCmd   `cmd:"" help:"Add tags to a bookmark."`
}

// ListCmd is a CLI command to list bookmarks/folders/tags.
type ListCmd struct {
	All         ListAllCmd         `cmd:"" help:"List bookmarks and folders."`
	Books       ListBooksCmd       `cmd:"" help:"List bookmarks."`
	Folders     ListFoldersCmd     `cmd:"" help:"List folders."`
	Tags        ListTagsCmd        `cmd:"" help:"List tags."`
	ParentNames ListParentNamesCmd `cmd:"" help:"List parent names."`
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

// GetCmd is a CLI command to get bookmarks/folders.
type GetCmd struct {
	All GetAllCmd `cmd:"" help:"Get a bookmark or folder."`
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

// ManifestCmd is a CLI command to manage the app manifest.
type ManifestCmd struct {
	Install InstallManifestCmd `cmd:"" help:"Install the app manifest"`
}

// InstallManifestCmd is a CLI command to the app manifest.
type InstallManifestCmd struct {
	Firefox  InstallFirefoxManifestCmd  `cmd:"" help:"Install the app manifest for Firefox."`
	Chrome   InstallChromeManifestCmd   `cmd:"" help:"Install the app manifest for Chrome."`
	Chromium InstallChromiumManifestCmd `cmd:"" help:"Install the app manifest for Chromium."`
}

// AddBookCmd is a CLI command to add a bookmark.
type AddBookCmd struct {
	Folder      *string  `help:"Folder to add this bookmark to."`
	Name        *string  `help:"Name for the bookmark."`
	Description *string  `help:"Description of the bookmark."`
	Tag         []string `help:"Tag to apply to the bookmark."`

	URL string `arg:"" name:"url" help:"URL of the bookmark."`
}

// Run add a bookmark.
func (r *AddBookCmd) Run(ctx *Context) error {
	start := time.Now()

	options := armaria.DefaultAddBookOptions()
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

	book, err := armaria.AddBook(r.URL, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []armaria.Book{book})
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

	options := armaria.DefaultAddFolderOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	if r.Folder != nil {
		options.WithParentID(*r.Folder)
	}

	book, err := armaria.AddFolder(r.Name, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []armaria.Book{book})
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

	options := armaria.DefaultAddTagsOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	book, err := armaria.AddTags(r.ID, r.Tag, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []armaria.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Tagged in %s", elapsed))

	return nil
}

// ListAllCmd is a CLI command to list bookmarks and folders.
type ListAllCmd struct {
	Folder   *string           `help:"Folder to list bookmarks/folders in."`
	NoFolder bool              `help:"List top level bookmarks/folders."`
	After    *string           `help:"ID of bookmark/folder to return results after."`
	Query    *string           `help:"Query to search bookmarks/folders by."`
	Tag      []string          `help:"Tag to filter bookmarks/folders by."`
	Order    armaria.Order     `help:"Field results are ordered on: modified/name/manual." enum:"modified,name,manual" default:"manual"`
	Dir      armaria.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First    *int64            `help:"The max number of bookmarks/folders to return."`
}

// Run list bookmarks and folders.
func (r *ListAllCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := armaria.DefaultListBooksOptions()
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

	books, err := armaria.ListBooks(options)
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
	Folder   *string           `help:"Folder to list bookmarks in."`
	NoFolder bool              `help:"List top level bookmarks."`
	After    *string           `help:"ID of bookmark to return results after."`
	Query    *string           `help:"Query to search bookmarks by."`
	Tag      []string          `help:"Tag to filter bookmarks by."`
	Order    armaria.Order     `help:"Field results are ordered on: modified/name/manual." enum:"modified,name,manual" default:"manual"`
	Dir      armaria.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First    *int64            `help:"The max number of bookmarks to return."`
}

// Run list bookmarks.
func (r *ListBooksCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := armaria.DefaultListBooksOptions()
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

	books, err := armaria.ListBooks(options)
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
	Folder   *string           `help:"Folder to list folders in."`
	NoFolder bool              `help:"List top level folders."`
	After    *string           `help:"ID of folder to return results after."`
	Query    *string           `help:"Query to search folders by."`
	Tag      []string          `help:"Tag to filter folders by."`
	Order    armaria.Order     `help:"Field results are ordered on: modified/name/manual." enum:"modified,name,manual" default:"manual"`
	Dir      armaria.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First    *int64            `help:"The max number of folders to return."`
}

// Run list folders.
func (r *ListFoldersCmd) Run(ctx *Context) error {
	start := time.Now()

	if r.NoFolder && r.Folder != nil {
		formatError(ctx.Writer, ctx.Formatter, ErrFolderNoFolderMutuallyExclusive)
		ctx.ReturnCode(1)
		return nil
	}

	options := armaria.DefaultListBooksOptions()
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

	books, err := armaria.ListBooks(options)
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
	Query *string           `help:"Query to search tags by."`
	After *string           `help:"ID of tags to return results after."`
	Dir   armaria.Direction `help:"Direction results are ordered by: asc/desc." enum:"asc,desc" default:"asc"`
	First *int64            `help:"The max number of tags to return."`
}

// Run list tags.
func (r *ListTagsCmd) Run(ctx *Context) error {
	start := time.Now()

	options := armaria.DefaultListTagsOptions()
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

	tags, err := armaria.ListTags(options)
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

// ListParentNamesCmd is a CLI command to get the parent names of a bookmark/folder.
type ListParentNamesCmd struct {
	ID string `arg:"" name:"id" help:"ID of the bookmark/folder to the parent names of."`
}

// Run get the parent names of a bookmark
func (r *ListParentNamesCmd) Run(ctx *Context) error {
	start := time.Now()

	options := armaria.DefaultGetParentNameOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	names, err := armaria.GetParentNames(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatParentNames(ctx.Writer, ctx.Formatter, names)
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
	Before        *string `help:"Book to order this bookmark before."`
	After         *string `help:"Book to order this bookmark after."`

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

	options := armaria.DefaultUpdateBookOptions()
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
	if r.Before != nil {
		options.WithOrderBefore(*r.Before)
	}
	if r.After != nil {
		options.WithOrderAfter(*r.After)
	}

	book, err := armaria.UpdateBook(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []armaria.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Updated in %s", elapsed))

	return nil
}

// UpdateFolderCmd is a CLI command to update a folder.
type UpdateFolderCmd struct {
	Name     *string `help:"New name for this folder."`
	Folder   *string `help:"Folder to move this folder to."`
	NoFolder bool    `help:"Remove the parent folder."`
	Before   *string `help:"Book to order this bookmark before."`
	After    *string `help:"Book to order this bookmark after."`

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

	options := armaria.DefaultUpdateFolderOptions()
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
	if r.Before != nil {
		options.WithOrderBefore(*r.Before)
	}
	if r.After != nil {
		options.WithOrderAfter(*r.After)
	}

	book, err := armaria.UpdateFolder(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []armaria.Book{book})
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

	options := armaria.DefaultRemoveBookOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	err := armaria.RemoveBook(r.ID, options)
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

	options := armaria.DefaultRemoveFolderOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	err := armaria.RemoveFolder(r.ID, options)
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

	options := armaria.DefaultRemoveTagsOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	book, err := armaria.RemoveTags(r.ID, r.Tag, options)

	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []armaria.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Untagged in %s", elapsed))

	return nil
}

// GetAllCmd is a CLI command to get a bookmark or folder.
type GetAllCmd struct {
	ID string `arg:"" name:"id" help:"ID of the bookmark or folder to get."`
}

// Run get bookmark or folder.
func (r *GetAllCmd) Run(ctx *Context) error {
	start := time.Now()

	options := armaria.DefaultGetBookOptions()
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}

	book, err := armaria.GetBook(r.ID, options)
	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatBookResults(ctx.Writer, ctx.Formatter, []armaria.Book{book})
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Got in %s", elapsed))

	return nil
}

// GetDBConfigCmd is a CLI command to get the location of the bookmarks database from the config.
type GetDBConfigCmd struct {
}

// Run get the location of the bookmarks database from the config.
func (r *GetDBConfigCmd) Run(ctx *Context) error {
	start := time.Now()

	config, err := armaria.GetConfig()

	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
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

	err := armaria.UpdateConfig(func(config *armaria.Config) {
		config.DB = r.DB
	})

	if err != nil && !errors.Is(err, armaria.ErrConfigMissing) {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Set in %s", elapsed))

	return nil
}

// InstallFirefoxManifestCmd is a CLI command to install the app manifest for Firefox.
type InstallFirefoxManifestCmd struct {
}

// Run install app manifest for Firefox.
func (r *InstallFirefoxManifestCmd) Run(ctx *Context) error {
	start := time.Now()

	err := armaria.InstallManifestFirefox()

	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Installed in %s", elapsed))

	return nil
}

// InstallChromeManifestCmd is a CLI command to install the app manifest for Chrome.
type InstallChromeManifestCmd struct {
}

// Run install app manifest for Chrome.
func (r *InstallChromeManifestCmd) Run(ctx *Context) error {
	start := time.Now()

	err := armaria.InstallManifestChrome()

	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Installed in %s", elapsed))

	return nil
}

// InstallChromiumManifestCmd is a CLI command to install the app manifest for Chromium.
type InstallChromiumManifestCmd struct {
}

// Run install app manifest for Chromium.
func (r *InstallChromiumManifestCmd) Run(ctx *Context) error {
	start := time.Now()

	err := armaria.InstallManifestChromium()

	if err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	elapsed := time.Since(start)

	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf("Installed in %s", elapsed))

	return nil
}

// VersionCmd is a CLI command to print the current version.
type VersionCmd struct {
}

// Run print the current version.
func (r *VersionCmd) Run(ctx *Context) error {
	formatSuccess(ctx.Writer, ctx.Formatter, fmt.Sprintf(ctx.Version))

	return nil
}

// QueryCmd is a CLI command to query bookmarks.
type QueryCmd struct {
	First int64 `help:"The max number of bookmarks/folders to return." default:"5"`

	Query string `arg:"" name:"query" help:"Query to search by."`
}

// Run query bookmarks.
func (r *QueryCmd) Run(ctx *Context) error {
	start := time.Now()

	options := armaria.DefaultListBooksOptions()
	options.WithFolders(true)
	options.WithBooks(true)
	if ctx.DB != nil {
		options.WithDB(*ctx.DB)
	}
	options.WithFirst(r.First)
	options.WithQuery(r.Query)

	books, err := armaria.ListBooks(options)
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

// TUICommand is a CLI command to start the TUI.
type TUICommand struct {
}

// Run start the TUI.
func (r *TUICommand) Run(ctx *Context) error {
	p := tui.Program()
	if _, err := p.Run(); err != nil {
		formatError(ctx.Writer, ctx.Formatter, err)
		ctx.ReturnCode(1)
		return nil
	}

	return nil
}

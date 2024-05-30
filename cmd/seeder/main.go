package main

import (
	"github.com/alecthomas/kong"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jonathanhope/armaria/pkg"
	"github.com/samber/lo"
)

// Used to seed data for testing.

const numTags = 3

type Context struct {
}

type SmallCmd struct {
}

func (r *SmallCmd) Run(_ *Context) error {
	// Add 10 folders.

	folders := make([]string, 0)
	for i := 0; i < 10; i++ {
		fo := armaria.DefaultAddFolderOptions()
		fr, err := armaria.AddFolder(gofakeit.ProductName(), fo)
		if err != nil {
			return err
		}

		folders = append(folders, fr.ID)
	}

	// Add 50 top level bookmarks.

	for i := 0; i < 50; i++ {
		bo := armaria.
			DefaultAddBookOptions().
			WithDescription(gofakeit.ProductDescription()).
			WithName(gofakeit.ProductName()).
			WithTags(tagsFactory())
		_, err := armaria.AddBook(gofakeit.URL(), bo)
		if err != nil {
			return err
		}
	}

	// Add 50 bookmarks to each folder.

	for _, f := range folders {
		for i := 0; i < 50; i++ {
			bo := armaria.
				DefaultAddBookOptions().
				WithDescription(gofakeit.ProductDescription()).
				WithName(gofakeit.ProductName()).
				WithTags(tagsFactory()).
				WithParentID(f)
			_, err := armaria.AddBook(gofakeit.URL(), bo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type MediumCmd struct {
}

func (r *MediumCmd) Run(_ *Context) error {
	// Add 50 folders.

	folders := make([]string, 0)
	for i := 0; i < 50; i++ {
		fo := armaria.DefaultAddFolderOptions()
		fr, err := armaria.AddFolder(gofakeit.ProductName(), fo)
		if err != nil {
			return err
		}

		folders = append(folders, fr.ID)
	}

	// Add 1000 top level bookmarks.

	for i := 0; i < 1000; i++ {
		bo := armaria.
			DefaultAddBookOptions().
			WithDescription(gofakeit.ProductDescription()).
			WithName(gofakeit.ProductName()).
			WithTags(tagsFactory())
		_, err := armaria.AddBook(gofakeit.URL(), bo)
		if err != nil {
			return err
		}
	}

	// Add 1000 bookmarks to each folder.

	for _, f := range folders {
		for i := 0; i < 1000; i++ {
			bo := armaria.
				DefaultAddBookOptions().
				WithDescription(gofakeit.ProductDescription()).
				WithName(gofakeit.ProductName()).
				WithTags(tagsFactory()).
				WithParentID(f)
			_, err := armaria.AddBook(gofakeit.URL(), bo)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func tagsFactory() []string {
	tags := []string{}

	for len(tags) < numTags {
		tag := gofakeit.NounCommon()
		count := lo.Count(tags, tag)
		if tag != "" && count == 0 {
			tags = append(tags, tag)
		}
	}

	return tags
}

type cli struct {
	Small  SmallCmd  `cmd:"" help:"Seed a small amount of data."`
	Medium MediumCmd `cmd:"" help:"Seed a medium amount of data."`
}

func main() {
	ctx := kong.Parse(&cli{})
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}

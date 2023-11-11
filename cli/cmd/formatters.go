package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jonathanhope/armaria/lib"
	"github.com/samber/lo"
	"golang.org/x/term"
)

// Formatter is the way output should be formatted.
type Formatter string

const (
	FormatterJSON   Formatter = "json"
	FormatterPretty Formatter = "pretty"
)

// BookDTO is a bookmark or folder that can be marshalled into JSON.
type BookDTO struct {
	ID          string         `json:"id"`
	URL         lib.NullString `json:"url"`
	Name        string         `json:"name"`
	Description lib.NullString `json:"description"`
	ParentID    lib.NullString `json:"parentId"`
	IsFolder    bool           `json:"idFolder"`
	ParentName  lib.NullString `json:"parentName"`
	Tags        []string       `json:"tags"`
}

// formatSuccess formats a success message.
// Success messages are not written in json mode.
func formatSuccess(writer io.Writer, formatter Formatter, message string) {
	switch formatter {

	case FormatterJSON:

	case FormatterPretty:
		style := lipgloss.
			NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.RoundedBorder())

		fmt.Fprintln(writer, style.Render(message))
	}
}

// formatError formats an error message.
func formatError(writer io.Writer, formatter Formatter, err error) {
	var errorString string
	if errors.Is(err, lib.ErrURLTooShort) {
		errorString = "URL too short"
	} else if errors.Is(err, lib.ErrURLTooLong) {
		errorString = "URL too long"
	} else if errors.Is(err, lib.ErrBookNotFound) {
		errorString = "Bookmark not found"
	} else if errors.Is(err, lib.ErrFolderNotFound) {
		errorString = "Folder not found"
	} else if errors.Is(err, lib.ErrNameTooShort) {
		errorString = "Name too short"
	} else if errors.Is(err, lib.ErrNameTooLong) {
		errorString = "Name too long"
	} else if errors.Is(err, lib.ErrDescriptionTooShort) {
		errorString = "Description too short"
	} else if errors.Is(err, lib.ErrDescriptionTooLong) {
		errorString = "Description too long"
	} else if errors.Is(err, lib.ErrTagTooShort) {
		errorString = "Tag too short"
	} else if errors.Is(err, lib.ErrTagTooLong) {
		errorString = "Tag too long"
	} else if errors.Is(err, lib.ErrDuplicateTag) {
		errorString = "Tags must be unique"
	} else if errors.Is(err, lib.ErrTooManyTags) {
		errorString = "Too many tags applied to bookmark"
	} else if errors.Is(err, lib.ErrTagInvalidChar) {
		errorString = "Tag has invalid chars"
	} else if errors.Is(err, lib.ErrNoUpdate) {
		errorString = "At least one update is required"
	} else if errors.Is(err, ErrFolderNoFolderMutuallyExclusive) {
		errorString = "Arguments folder and no-folder are mutually exclusive"
	} else if errors.Is(err, ErrDescriptionNoDescriptionMutuallyExclusive) {
		errorString = "Arguments description and no-description are mutually exclusive"
	} else if errors.Is(err, lib.ErrTagNotFound) {
		errorString = "Tag not found"
	} else if errors.Is(err, lib.ErrFirstTooSmall) {
		errorString = "First too small"
	} else if errors.Is(err, lib.ErrQueryTooShort) {
		errorString = "Query too short"
	} else {
		errorString = err.Error()
	}

	switch formatter {

	case FormatterJSON:
		fmt.Fprintf(writer, "\"%s\"\n", errorString)

	case FormatterPretty:
		style := lipgloss.
			NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.RoundedBorder())

		fmt.Fprintln(writer, style.Render(errorString))
	}
}

// formatBookResults formats a collection of bookmarks/folders.
func formatBookResults(writer io.Writer, formatter Formatter, books []lib.Book) {

	switch formatter {
	case FormatterJSON:
		dtos := lo.Map(books, func(x lib.Book, index int) BookDTO {
			return BookDTO{
				ID:          x.ID,
				URL:         lib.NullStringFromPtr(x.URL),
				Name:        x.Name,
				Description: lib.NullStringFromPtr(x.Description),
				ParentID:    lib.NullStringFromPtr(x.ParentID),
				IsFolder:    x.IsFolder,
				ParentName:  lib.NullStringFromPtr(x.ParentName),
				Tags:        x.Tags,
			}
		})

		json, err := json.Marshal(&dtos)
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(writer, string(json))

	case FormatterPretty:
		width, _, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}

		headerStyle := lipgloss.
			NewStyle().
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			Width(16)

		rowStyle := lipgloss.
			NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Width(width - 16)

		for _, book := range books {
			rows := [][]string{
				{formatIsFolder(book.IsFolder), book.ID},
				{"Name", book.Name},
				{"URL", formatNullableString(book.URL)},
				{"Description", formatNullableString(book.Description)},
				{"Folder", formatNullableString(book.ParentName)},
				{"Tags", formatTags(book.Tags)},
			}

			table := table.New().
				Border(lipgloss.RoundedBorder()).
				BorderRow(true).
				BorderColumn(true).
				Width(width).
				StyleFunc(func(row, col int) lipgloss.Style {
					switch {
					case col == 0:
						return headerStyle
					default:
						return rowStyle
					}
				}).
				Rows(rows...)

			fmt.Fprintln(writer, table)
		}
	}
}

// formatTagResults formats a collection of tags.
func formatTagResults(writer io.Writer, formatter Formatter, tags []string) {
	switch formatter {

	case FormatterJSON:
		json, err := json.Marshal(&tags)

		if err != nil {
			panic(err)
		}

		fmt.Fprintln(writer, string(json))

	case FormatterPretty:
		width, _, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}

		style := lipgloss.
			NewStyle().
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.RoundedBorder()).
			MaxWidth(width - 2)

		for _, tag := range tags {
			fmt.Fprintln(writer, style.Render(fmt.Sprintf("🏷  %s", tag)))
		}
	}
}

// formatIsFolder formats an is folder value.
func formatIsFolder(isFolder bool) string {
	if isFolder {
		return "🗁"
	}

	return "🕮"
}

// formatNullableString formats a nullable string.
func formatNullableString(str *string) string {
	if str != nil {
		return *str
	}

	return "NULL"
}

// formatTags formats a tags value.
func formatTags(tags []string) string {
	return strings.Join(tags, ", ")
}

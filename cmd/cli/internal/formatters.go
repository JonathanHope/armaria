package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jonathanhope/armaria/cmd/cli/internal/messaging"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/samber/lo"
)

// The Armaria CLI can format its output in multiple ways.
// This file contains the logic for all of the output formatting.

// Formatter is the way output should be formatted.
type Formatter string

const (
	FormatterJSON   Formatter = "json"   // format results as JSON.
	FormatterPretty Formatter = "pretty" // format results in human readable way.
)

// marshalMsg will marshal a native message as JSON.
func marshalMsg[T messaging.Payload](kind messaging.MessageKind, payload T) string {
	msg, err := messaging.PayloadToMessage(kind, payload)
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return string(json)
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
	if errors.Is(err, armaria.ErrURLTooShort) {
		errorString = "URL too short"
	} else if errors.Is(err, armaria.ErrURLTooLong) {
		errorString = "URL too long"
	} else if errors.Is(err, armaria.ErrBookNotFound) {
		errorString = "Bookmark not found"
	} else if errors.Is(err, armaria.ErrFolderNotFound) {
		errorString = "Folder not found"
	} else if errors.Is(err, armaria.ErrNotFound) {
		errorString = "Bookmark or folder not found"
	} else if errors.Is(err, armaria.ErrNameTooShort) {
		errorString = "Name too short"
	} else if errors.Is(err, armaria.ErrNameTooLong) {
		errorString = "Name too long"
	} else if errors.Is(err, armaria.ErrDescriptionTooShort) {
		errorString = "Description too short"
	} else if errors.Is(err, armaria.ErrDescriptionTooLong) {
		errorString = "Description too long"
	} else if errors.Is(err, armaria.ErrTagTooShort) {
		errorString = "Tag too short"
	} else if errors.Is(err, armaria.ErrTagTooLong) {
		errorString = "Tag too long"
	} else if errors.Is(err, armaria.ErrDuplicateTag) {
		errorString = "Tags must be unique"
	} else if errors.Is(err, armaria.ErrTooManyTags) {
		errorString = "Too many tags applied to bookmark"
	} else if errors.Is(err, armaria.ErrTagInvalidChar) {
		errorString = "Tag has invalid chars"
	} else if errors.Is(err, armaria.ErrNoUpdate) {
		errorString = "At least one update is required"
	} else if errors.Is(err, ErrFolderNoFolderMutuallyExclusive) {
		errorString = "Arguments folder and no-folder are mutually exclusive"
	} else if errors.Is(err, ErrDescriptionNoDescriptionMutuallyExclusive) {
		errorString = "Arguments description and no-description are mutually exclusive"
	} else if errors.Is(err, armaria.ErrTagNotFound) {
		errorString = "Tag not found"
	} else if errors.Is(err, armaria.ErrFirstTooSmall) {
		errorString = "First too small"
	} else if errors.Is(err, armaria.ErrQueryTooShort) {
		errorString = "Query too short"
	} else {
		errorString = err.Error()
	}

	switch formatter {

	case FormatterJSON:
		json := marshalMsg(messaging.MessageKindError, messaging.ErrorPayload{
			Error: errorString,
		})
		fmt.Fprint(writer, json)

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
func formatBookResults(writer io.Writer, formatter Formatter, books []armaria.Book) {

	switch formatter {

	case FormatterJSON:
		dtos := lo.Map(books, func(x armaria.Book, index int) messaging.BookDTO {
			return messaging.BookDTO{
				ID:          x.ID,
				URL:         null.NullStringFromPtr(x.URL),
				Name:        x.Name,
				Description: null.NullStringFromPtr(x.Description),
				ParentID:    null.NullStringFromPtr(x.ParentID),
				IsFolder:    x.IsFolder,
				ParentName:  null.NullStringFromPtr(x.ParentName),
				Tags:        x.Tags,
			}
		})

		json := marshalMsg(messaging.MessageKindBooks, messaging.BooksPayload{
			Books: dtos,
		})

		fmt.Fprintln(writer, json)

	case FormatterPretty:
		width, _ := consolesize.GetConsoleSize()

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
		json := marshalMsg(messaging.MessageKindTags, messaging.TagsPayload{
			Tags: tags,
		})

		fmt.Fprintln(writer, json)

	case FormatterPretty:
		width, _ := consolesize.GetConsoleSize()

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

// formatConfigResult formats a config value.
func formatConfigResult(writer io.Writer, formatter Formatter, value string) {
	switch formatter {

	case FormatterJSON:
		json := marshalMsg(messaging.MessageKindConfigValue, messaging.ConfigValuePayload{
			Value: value,
		})

		fmt.Fprint(writer, json)

	case FormatterPretty:
		width, _ := consolesize.GetConsoleSize()

		style := lipgloss.
			NewStyle().
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.RoundedBorder()).
			MaxWidth(width - 2)

		if value == "" {
			fmt.Fprintln(writer, style.Render("⚙ N/A"))
		} else {
			fmt.Fprintln(writer, style.Render(fmt.Sprintf("⚙ %s", value)))
		}
	}
}

// formatParentNames formats a collection of parent names.
func formatParentNames(writer io.Writer, formatter Formatter, parentNames []string) {
	switch formatter {

	case FormatterJSON:
		json := marshalMsg(messaging.MessageKindParentNames, messaging.ParentNamesPayload{
			ParentNames: parentNames,
		})

		fmt.Fprintln(writer, json)

	case FormatterPretty:
		width, _ := consolesize.GetConsoleSize()

		style := lipgloss.
			NewStyle().
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.RoundedBorder()).
			MaxWidth(width - 2)

		fmt.Fprintln(writer, style.Render(strings.Join(parentNames, " > ")))
	}
}

// formatIsFolder formats an is folder value.
func formatIsFolder(isFolder bool) string {
	if isFolder {
		return "📁"
	}

	return "📖"
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

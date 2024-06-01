package test

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/cucumber/godog"
	"github.com/google/shlex"
	"github.com/google/uuid"
	"github.com/jonathanhope/armaria/cmd/cli/internal"
	"github.com/jonathanhope/armaria/cmd/cli/internal/messaging"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg"
)

// invokeCli runs the Armaria CLI with the provided args.
func invokeCli(args string) (string, error) {
	// All of this is to invoke a Kong CLI app directly in code.
	// A buffer is used to intercept output.

	rootCmd := cmd.RootCmdFactory()
	w := bytes.NewBuffer(nil)

	options := []kong.Option{
		kong.Name("test"),
		kong.Exit(func(int) {
			panic(true)
		}),
		kong.Writers(w, w),
	}

	parser, err := kong.New(&rootCmd, options...)
	if err != nil {
		return "", err
	}

	// We need to use this parser to take shell quoting into account.
	tokens, err := shlex.Split(args)
	if err != nil {
		return "", err
	}

	ctx, err := parser.Parse(tokens)
	if err != nil {
		return "", err
	}

	err = ctx.Run(&cmd.Context{
		DB:         rootCmd.DB,
		Formatter:  rootCmd.Formatter,
		Writer:     w,
		ReturnCode: noop})
	if err != nil {
		return "", err
	}

	return w.String(), nil
}

// insertArgs is params to insert a bookmark or folder.
type insertArgs struct {
	db          string
	vars        map[string]interface{}
	id          string
	parentID    string
	isFolder    string
	name        string
	url         string
	description string
	tags        string
}

// insert inserts a bookmark or folder directly into the bookmarks DB.
func insert(args insertArgs) error {
	_, storeId, idKey, err := handleString(args.vars, args.id)
	if err != nil {
		return err
	}

	parentID, storeParentID, parentIdKey, err := handleNullString(args.vars, args.parentID)
	if err != nil {
		return err
	}

	isFolder, storeIsFolder, isFolderKey, err := handleBool(args.vars, args.isFolder)
	if err != nil {
		return err
	}

	name, storeName, nameKey, err := handleString(args.vars, args.name)
	if err != nil {
		return err
	}

	URL, storeURL, URLKey, err := handleNullString(args.vars, args.url)
	if err != nil {
		return err
	}

	desc, storeDesc, descKey, err := handleNullString(args.vars, args.description)
	if err != nil {
		return err
	}

	tags, storeTags, tagsKey, err := handleCompoundString(args.vars, args.tags)
	if err != nil {
		return err
	}

	var book armaria.Book
	if !isFolder {
		options := armaria.DefaultAddBookOptions()
		options.WithDB(args.db)
		if parentID.Valid {
			options.WithParentID(parentID.String)
		}
		if name != "" {
			options.WithName(name)
		}
		if desc.Valid {
			options.WithDescription(desc.String)
		}
		if tags != nil {
			options.WithTags(tags)
		}

		book, err = armaria.AddBook(URL.String, options)
		if err != nil {
			return err
		}
	} else {
		options := armaria.DefaultAddFolderOptions()
		options.WithDB(args.db)
		if parentID.Valid {
			options.WithParentID(parentID.String)
		}

		book, err = armaria.AddFolder(name, options)
		if err != nil {
			return err
		}
	}

	// Store any variables that were denoted with {...} in the cucumber table.

	if storeId {
		args.vars[idKey] = book.ID
	}

	if storeParentID {
		args.vars[parentIdKey] = null.NullStringFromPtr(book.ParentID)
	}

	if storeIsFolder {
		args.vars[isFolderKey] = book.IsFolder
	}

	if storeName {
		args.vars[nameKey] = book.Name
	}

	if storeURL {
		args.vars[URLKey] = null.NullStringFromPtr(book.URL)
	}

	if storeDesc {
		args.vars[descKey] = null.NullStringFromPtr(book.Description)
	}

	if storeTags {
		args.vars[tagsKey] = book.Tags
	}

	return nil
}

// tableToBooks converts a cucumber table to a collection of bookmarks/folders.
func tableToBooks(vars map[string]interface{}, actual []messaging.BookDTO, table *godog.Table) ([]messaging.BookDTO, error) {
	books := make([]messaging.BookDTO, 0)
	for _, row := range table.Rows[1:] {
		id, storeId, _, err := handleString(vars, row.Cells[0].Value)
		if err != nil {
			return nil, err
		}

		parentId, storeParentId, _, err := handleNullString(vars, row.Cells[1].Value)
		if err != nil {
			return nil, err
		}

		isFolder, storeIsFolder, _, err := handleBool(vars, row.Cells[2].Value)
		if err != nil {
			return nil, err
		}

		name, storeName, _, err := handleString(vars, row.Cells[3].Value)
		if err != nil {
			return nil, err
		}

		URL, storeURL, _, err := handleNullString(vars, row.Cells[4].Value)
		if err != nil {
			return nil, err
		}

		desc, storeDesc, _, err := handleNullString(vars, row.Cells[5].Value)
		if err != nil {
			return nil, err
		}

		tags, storeTags, _, err := handleCompoundString(vars, row.Cells[6].Value)
		if err != nil {
			return nil, err
		}

		// The {...} is used as a wildcard in the result table.
		// Using the last inserted result was the best way I could think of to model this.

		last := actual[len(actual)-1]

		if storeId {
			id = last.ID
		}

		if storeParentId {
			parentId = last.ParentID
		}

		if storeIsFolder {
			isFolder = last.IsFolder
		}

		if storeName {
			name = last.Name
		}

		if storeURL {
			URL = last.URL
		}

		if storeDesc {
			desc = last.Description
		}

		if storeTags {
			tags = last.Tags
		}

		book := messaging.BookDTO{
			ID:          id,
			ParentID:    parentId,
			IsFolder:    isFolder,
			Name:        name,
			URL:         URL,
			Description: desc,
			Tags:        tags,
		}

		books = append(books, book)
	}

	return books, nil
}

// tableToTags converts a cucumber table to a collection of tags.
func tableToTags(table *godog.Table) []string {
	tags := make([]string, 0)
	for _, row := range table.Rows[1:] {
		tags = append(tags, row.Cells[0].Value)
	}

	return tags
}

// handleString handles a string value in a cucumber table.
// Returns a flag to denote if result should be stored as a variable.
// Will retreive a variable if the [...] denotation is used.
func handleString(vars map[string]interface{}, val string) (res string, store bool, key string, err error) {
	if isVariableSet(val) {
		store = true
		key = stripVariable(val)
		return
	} else if isVariableGet(val) {
		val, ok := vars[stripVariable(val)]
		if !ok {
			err = errors.New("Missing variable")
			return
		}

		if res, ok = val.(string); !ok {
			err = fmt.Errorf("Could not convert variable to string; type was %T", val)
			return
		}
	} else {
		res = val
	}

	return
}

// handleNullString handles a nullable string value in a cucumber table.
// Returns a flag to denote if result should be stored as a variable.
// Will retreive a variable if the [...] denotation is used.
// Handles NULL literal correctly.
func handleNullString(vars map[string]interface{}, val string) (res null.NullString, store bool, key string, err error) {
	if isVariableSet(val) {
		store = true
		key = stripVariable(val)
		return
	} else if isVariableGet(val) {
		val, ok := vars[stripVariable(val)]
		if !ok {
			err = errors.New("Missing variable")
			return
		}

		if converted, ok := val.(null.NullString); ok {
			res = converted
			return
		} else if converted, ok := val.(string); ok {
			res = null.NullStringFrom(converted)
		} else {
			err = fmt.Errorf("Could not convert variable to string; type was %T", val)
		}
	} else if val != "NULL" {
		res = null.NullStringFrom(val)
	}

	return
}

// handleBool handles a bool value in a cucumber table.
// Returns a flag to denote if result should be stored as a variable.
// Will retreive a variable if the [...] denotation is used.
func handleBool(vars map[string]interface{}, val string) (res bool, store bool, key string, err error) {
	if isVariableSet(val) {
		store = true
		key = stripVariable(val)
		return
	} else if isVariableGet(val) {
		val, ok := vars[stripVariable(val)]
		if !ok {
			err = errors.New("Missing variable")
			return
		}

		if res, ok = val.(bool); !ok {
			err = fmt.Errorf("Could not convert variable to bool; type was %T", val)
			return
		}
	} else {
		parsed, parseErr := strconv.ParseBool(val)
		if parseErr != nil {
			err = parseErr
			return
		}

		res = parsed
	}

	return
}

// handleCompoundString handles a comma separated string value.
// Returns a flag to denote if result should be stored as a variable.
// Will retreive a variable if the [...] denotation is used.
func handleCompoundString(vars map[string]interface{}, val string) (res []string, store bool, key string, err error) {
	if val == "" {
		res = make([]string, 0)
		return
	}

	if isVariableSet(val) {
		store = true
		key = stripVariable(val)
		return
	} else if isVariableGet(val) {
		val, ok := vars[stripVariable(val)]
		if !ok {
			err = errors.New("Missing variable")
			return
		}

		if res, ok = val.([]string); !ok {
			err = fmt.Errorf("Could not convert variable to []string; type was %T", val)
			return
		}
	} else {
		res = strings.Split(val, ", ")
	}

	return
}

// isVariableGet returns true if the cell is a request to get a variable.
func isVariableGet(value string) bool {
	return strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]")
}

// isVariableSet returns true if the cell is a request to store a variable.
func isVariableSet(value string) bool {
	return strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")
}

// stripVariable returns the variable name if it is wrapped in {...} or [...].
func stripVariable(value string) string {
	return value[1 : len(value)-1]
}

// variableToString converts a stored variable to string.
func variableToString(value interface{}) (string, error) {
	if v, ok := value.(string); ok {
		return v, nil
	} else if v, ok := value.(bool); ok {
		return strconv.FormatBool(v), nil
	} else if v, ok := value.(null.NullString); ok {
		if v.Valid {
			return v.String, nil
		} else {
			return "NULL", nil
		}
	} else {
		return "", errors.New("Could not convert variable to string")
	}
}

// noop does nothing.
// We don't want to actually os.exit here since that would exit the test.
func noop(code int) {}

// markDirty marks every field as dirty.
// The dirty field is not relevant to tests.
func markDirty(expected []messaging.BookDTO, actual []messaging.BookDTO) {
	for i := range expected {
		expected[i].Description.Dirty = false
		expected[i].ParentID.Dirty = false
		expected[i].ParentName.Dirty = false
		expected[i].URL.Dirty = false
		expected[i].ParentName = null.NullString{}
	}

	for i := range actual {
		actual[i].Description.Dirty = false
		actual[i].ParentID.Dirty = false
		actual[i].ParentName.Dirty = false
		actual[i].URL.Dirty = false
		actual[i].ParentName = null.NullString{}
	}
}

// processCommand swaps out special tokens in the CLI command.
func processCommand(vars map[string]interface{}, cmd string) (string, error) {
	// Replace any variable templates in the CLI args.
	// The per scenario DB name is passed in.
	// The JSON formatter is used so the output is machine readable.
	for k := range vars {
		str, err := variableToString(vars[k])
		if err != nil {
			return "", err
		}

		cmd = strings.Replace(cmd, fmt.Sprintf("[%s]", k), str, -1)
	}

	// Syntax like %repeat:{str}:{num}% can be used to generate strings of arbitrary length.
	// This is useful for bounds checking while keeping the Gherkin clean.
	var compRegEx = regexp.MustCompile(`\%repeat:(.+?):(\d+?)\%`)
	match := compRegEx.FindStringSubmatch(cmd)
	if len(match) == 3 {
		substr := match[1]
		length, err := strconv.Atoi(match[2])
		if err != nil {
			return "", err
		}

		var str = ""
		for i := 0; i < length; i++ {
			str = str + substr
		}

		cmd = compRegEx.ReplaceAllString(cmd, str)
	}

	// Syntax like [uuid] can be used to generate GUIDs.
	for strings.Contains(cmd, "[uuid]") {
		cmd = strings.Replace(cmd, "[uuid]", uuid.New().String(), 1)
	}

	return cmd, nil
}

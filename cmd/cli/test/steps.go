package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jonathanhope/armaria/cmd/cli/internal"
	"github.com/jonathanhope/armaria/internal/null"
)

// Context keys.

type dbContextKey struct{}
type outputContextKey struct{}
type variablesContextKey struct{}

// InitializeTestSuite wires up events.
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.ScenarioContext().Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// Each scenario run gets its own DB.
		db := fmt.Sprintf("%s.db", uuid.New())
		ctx = context.WithValue(ctx, dbContextKey{}, db)

		ctx = context.WithValue(ctx, variablesContextKey{}, make(map[string]interface{}))

		return ctx, nil
	})

	ctx.ScenarioContext().After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			return nil, err
		}

		db, ok := ctx.Value(dbContextKey{}).(string)
		if !ok {
			return ctx, errors.New("Missing DB name")
		}

		// When the test scenario is over delete the per scenario DB.

		if _, err := os.Stat(db); err == nil {
			if err = os.Remove(db); err != nil {
				os.Remove(db)
			}
		}

		shm := fmt.Sprintf("%s-shm", db)
		if _, err := os.Stat(shm); err == nil {
			if err = os.Remove(shm); err != nil {
				os.Remove(db)
			}
		}

		wal := fmt.Sprintf("%s-wal", db)
		if _, err := os.Stat(wal); err == nil {
			if err = os.Remove(wal); err != nil {
				os.Remove(db)
			}
		}

		return ctx, nil
	})
}

// InitializeScenario wires up the steps.
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the DB already has the following entries:$`, theDBAlreadyHasTheFollowingEntries)
	ctx.Step(`^I run it with the following args:$`, iRunItWithTheFollowingArgs)
	ctx.Step(`^the following bookmarks\/folders exist:$`, theFollowingBookmarksFoldersExist)
	ctx.Step(`^the folllowing tags exist:$`, theFollowingTagsExist)
	ctx.Step(`^the following error is returned:$`, theFollowingErrorIsReturned)
	ctx.Step(`^the folllowing tags are returned:$`, theFolllowingTagsAreReturned)
	ctx.Step(`^the folllowing names are returned:$`, theFolllowingNamesAreReturned)
	ctx.Step(`^the folllowing books are returned:$`, theFolllowingBooksAreReturned)
}

// theDBAlreadyHasTheFollowingEntries inserts data from a cucumber table into the bookmarks database.
func theDBAlreadyHasTheFollowingEntries(ctx context.Context, table *godog.Table) (context.Context, error) {
	vars, ok := ctx.Value(variablesContextKey{}).(map[string]interface{})
	if !ok {
		return ctx, errors.New("Missing variables")
	}

	db, ok := ctx.Value(dbContextKey{}).(string)
	if !ok {
		return ctx, errors.New("Missing DB name")
	}

	for _, row := range table.Rows[1:] {
		err := insert(insertArgs{
			db:          db,
			vars:        vars,
			id:          row.Cells[0].Value,
			parentID:    row.Cells[1].Value,
			isFolder:    row.Cells[2].Value,
			name:        row.Cells[3].Value,
			url:         row.Cells[4].Value,
			description: row.Cells[5].Value,
			tags:        row.Cells[6].Value,
		})
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

// iRunItWithTheFollowingArgs runs the CLI with the provided args.
func iRunItWithTheFollowingArgs(ctx context.Context, args *godog.DocString) (context.Context, error) {
	db, ok := ctx.Value(dbContextKey{}).(string)
	if !ok {
		return ctx, errors.New("Missing DB name")
	}

	vars, ok := ctx.Value(variablesContextKey{}).(map[string]interface{})
	if !ok {
		return ctx, errors.New("Missing variables")
	}

	cmd, err := processCommand(vars, fmt.Sprintf("%s --db %s --formatter json", args.Content, db))
	if err != nil {
		return ctx, err
	}

	// Store the output for future use.
	output, err := invokeCli(fmt.Sprintf("%s --db %s --formatter json", cmd, db))
	return context.WithValue(ctx, outputContextKey{}, output), err
}

// theFollowingBookmarksFoldersExist compares the JSON output of the list all command with a cucumber results table.
func theFollowingBookmarksFoldersExist(ctx context.Context, table *godog.Table) error {
	vars, ok := ctx.Value(variablesContextKey{}).(map[string]interface{})
	if !ok {
		return errors.New("Missing variables")
	}

	db, ok := ctx.Value(dbContextKey{}).(string)
	if !ok {
		return errors.New("Missing DB name")
	}

	output, ok := ctx.Value(outputContextKey{}).(string)
	if !ok {
		return errors.New("Missing DB name")
	}

	if strings.HasPrefix(output, `"`) {
		fmt.Printf("Unexpected Error: %s\n", output)
	}

	output, err := invokeCli(fmt.Sprintf("list all --db %s --formatter json", db))
	if err != nil {
		return err
	}

	var actual []cmd.BookDTO
	if err := json.Unmarshal([]byte(output), &actual); err != nil {
		return err
	}

	expected, err := tableToBooks(vars, actual, table)
	if err != nil {
		return err
	}

	// The cucumber tables don't have parent name so it gets nulled out.
	for i := range actual {
		actual[i].ParentName = null.NullStringFromPtr(nil)
	}

	markDirty(expected, actual)

	diff := cmp.Diff(expected, actual)
	if diff != "" {
		return fmt.Errorf("Expected and actual books different:\n%s", diff)
	}

	return nil
}

// theFollowingTagsExist compares the JSON output of the list all command with a cucumber results table.
func theFollowingTagsExist(ctx context.Context, table *godog.Table) error {
	db, ok := ctx.Value(dbContextKey{}).(string)
	if !ok {
		return errors.New("Missing DB name")
	}

	output, err := invokeCli(fmt.Sprintf("list tags --db %s --formatter json", db))
	if err != nil {
		return err
	}

	var actual []string
	if err := json.Unmarshal([]byte(output), &actual); err != nil {
		return err
	}

	expected := tableToTags(table)
	if err != nil {
		return err
	}

	if len(expected) == 0 && len(actual) == 0 {
		return nil
	}

	diff := cmp.Diff(expected, actual)
	if diff != "" {
		return fmt.Errorf("Expected and actual tags different:\n%s", diff)
	}

	return nil
}

// theFollowingErrorIsReturned compares the JSON output of the CLI run with a cucumber error string.
func theFollowingErrorIsReturned(ctx context.Context, message *godog.DocString) error {
	output, ok := ctx.Value(outputContextKey{}).(string)
	if !ok {
		return errors.New("Missing output")
	}

	expected := strings.TrimSpace(fmt.Sprintf("\"%s\"", message.Content))

	actual := strings.TrimSpace(output)

	diff := cmp.Diff(expected, actual)
	if diff != "" {
		return fmt.Errorf("Expected and actual errors different:\n%s", diff)
	}

	return nil
}

// theFolllowingTagsAreReturned compares the JSON output of the CLI with a set of tags.
func theFolllowingTagsAreReturned(ctx context.Context, table *godog.Table) error {
	output, ok := ctx.Value(outputContextKey{}).(string)
	if !ok {
		return errors.New("Missing output")
	}

	var expected []string
	for _, row := range table.Rows[1:] {
		expected = append(expected, row.Cells[0].Value)
	}

	var actual []string
	err := json.Unmarshal([]byte(output), &actual)
	if err != nil {
		return err
	}

	diff := cmp.Diff(expected, actual)
	if diff != "" {
		return fmt.Errorf("Expected and actual tags different:\n%s", diff)
	}

	return nil
}

// theFolllowingNamesAreReturned compares the JSON output of the CLI with a set of names.
func theFolllowingNamesAreReturned(ctx context.Context, table *godog.Table) error {
	output, ok := ctx.Value(outputContextKey{}).(string)
	if !ok {
		return errors.New("Missing output")
	}

	var expected []string
	for _, row := range table.Rows {
		expected = append(expected, row.Cells[0].Value)
	}

	var actual []string
	err := json.Unmarshal([]byte(output), &actual)
	if err != nil {
		return err
	}

	diff := cmp.Diff(expected, actual)
	if diff != "" {
		return fmt.Errorf("Expected and actual names different:\n%s", diff)
	}

	return nil
}

// theFolllowingBooksAreReturned compares the JSON output of the CLI with a set of bookmarks/folders.
func theFolllowingBooksAreReturned(ctx context.Context, table *godog.Table) error {
	output, ok := ctx.Value(outputContextKey{}).(string)
	if !ok {
		return errors.New("Missing output")
	}

	vars, ok := ctx.Value(variablesContextKey{}).(map[string]interface{})
	if !ok {
		return errors.New("Missing variables")
	}

	var actual []cmd.BookDTO
	err := json.Unmarshal([]byte(output), &actual)
	if err != nil {
		return err
	}

	expected, err := tableToBooks(vars, actual, table)
	if err != nil {
		return err
	}

	markDirty(expected, actual)

	diff := cmp.Diff(expected, actual)
	if diff != "" {
		return fmt.Errorf("Expected and actual books different:\n%s", diff)
	}

	return nil
}

package lib

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nullism/bqb"
	"github.com/samber/lo"
)

// This file contains the low level logic to access the bookmarks DB.

// types

// bookDto is a DTO to stuff DB results into.
type bookDto struct {
	ID          string     `db:"id"`
	URL         NullString `db:"url"`
	Name        string     `db:"name"`
	Description NullString `db:"description"`
	ParentID    NullString `db:"parent_id"`
	IsFolder    bool       `db:"is_folder"`
	ParentName  NullString `db:"parent_name"`
	Tags        string     `db:"tags"`
}

// create

// addBookDB inserts a book into the bookmarks database.
func addBookDB(tx transaction, url string, name string, description NullString, parentID NullString) (string, error) {
	id := uuid.New().String()

	insert := bqb.New(`INSERT INTO "bookmarks"("id", "url", "is_folder", "name", "description", "parent_id")`)
	insert.Space("VALUES(?, ?, ?, ?, ?, ?)", id, url, false, name, description, parentID)

	err := exec(tx, insert)
	return id, err
}

// addFolderDB inserts a folder into the bookmarks database.
func addFolderDB(tx transaction, name string, parentID NullString) (string, error) {
	id := uuid.New().String()

	insert := bqb.New(`INSERT INTO "bookmarks"("id", "is_folder", "name", "parent_id")`)
	insert.Space("VALUES(?, ?, ?, ?)", id, true, name, parentID)

	err := exec(tx, insert)
	return id, err
}

// addTagsDB inserts tags into the bookmarks database.
func addTagsDB(tx transaction, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	insert := bqb.New(`INSERT INTO "tags"("tag")`)
	insert.Space(`VALUES (?)`, tags[0])

	for _, tag := range tags[1:] {
		insert.Comma(`(?)`, tag)
	}

	return exec(tx, insert)
}

// linkTagsDB adds tags to bookmark.
func linkTagsDB(tx transaction, bookmarkID string, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	insert := bqb.New(`INSERT INTO "bookmarks_tags"("bookmark_id", "tag")`)
	insert.Space(`VALUES (?, ?)`, bookmarkID, tags[0])

	for _, tag := range tags[1:] {
		insert.Comma(`(?, ?)`, bookmarkID, tag)
	}

	return exec(tx, insert)
}

// read

// getBooksDBArgs are the args for getBooksDB.
type getBooksDBArgs struct {
	idFilter       string
	includeBooks   bool
	includeFolders bool
	parentID       NullString
	query          NullString
	tags           []string
	after          NullString
	order          Order
	direction      Direction
	first          NullInt64
}

// getBooksDB lists bookmarks/folders in the bookmarks DB.
func getBooksDB(tx transaction, args getBooksDBArgs) ([]Book, error) {
	tags := bqb.New(`SELECT GROUP_CONCAT("tag")`)
	tags.Space(`FROM "bookmarks_tags"`)
	tags.Space(`WHERE "bookmark_id" = "child"."id"`)

	books := bqb.New(`SELECT "child"."id"`)
	books.Comma(`"child"."url"`)
	books.Comma(`"child"."name"`)
	books.Comma(`"child"."description"`)
	books.Comma(`"child"."parent_id"`)
	books.Comma(`"child"."is_folder"`)
	books.Comma(`"parent"."name" AS "parent_name"`)
	books.Comma(`IFNULL((?), '') AS "tags"`, tags)
	books.Space(`FROM "bookmarks" AS "child"`)
	books.Space(`LEFT JOIN "bookmarks" AS "parent" ON "parent"."id" = "child"."parent_id"`)

	where := bqb.Optional("WHERE")

	if args.idFilter != "" {
		where.And(`"child"."id" = ?`, args.idFilter)
	}

	if args.includeBooks && !args.includeFolders {
		where.And(`"child"."is_folder" = ?`, false)
	}

	if args.includeFolders && !args.includeBooks {
		where.And(`"child"."is_folder" = ?`, true)
	}

	if args.parentID.Dirty && args.parentID.Valid {
		where.And(`"child"."parent_id" = ?`, args.parentID.String)
	} else if args.parentID.Dirty && !args.parentID.Valid {
		where.And(`"child"."parent_id" IS NULL`)
	}

	if args.query.Dirty && args.query.Valid {
		searchFilter := bqb.New(`SELECT "id"`)
		searchFilter.Space(`FROM "bookmarks_fts"`)
		searchFilter.Space(`WHERE "name" LIKE ?`, fmt.Sprintf("%%%s%%", args.query.String))
		searchFilter.Space(`OR "description" LIKE ?`, fmt.Sprintf("%%%s%%", args.query.String))
		searchFilter.Space(`OR "url" LIKE ?`, fmt.Sprintf("%%%s%%", args.query.String))

		where.And(`"child"."id" IN (?)`, searchFilter)
	}

	if len(args.tags) > 0 {
		tagsFilter := bqb.New(`SELECT "bookmark_id"`)
		tagsFilter.Space(`FROM "bookmarks_tags"`)
		tagsFilter.Space(`WHERE "tag" IN (?)`, args.tags)

		where.And(`"child"."id" IN (?)`, tagsFilter)
	}

	if args.after.Dirty && args.after.Valid {
		if args.order == OrderName && args.direction == DirectionAsc {
			where.And(`("child"."name" > (SELECT "name" FROM "bookmarks" WHERE "id" = ?)`, args.after.String)
		} else if args.order == OrderName && args.direction == DirectionDesc {
			where.And(`("child"."name" < (SELECT "name" FROM "bookmarks" WHERE "id" = ?)`, args.after.String)
		} else if args.order == OrderModified && args.direction == DirectionAsc {
			where.And(`("child"."modified" > (SELECT "modified" FROM "bookmarks" WHERE "id" = ?)`, args.after.String)
		} else if args.order == OrderModified && args.direction == DirectionDesc {
			where.And(`("child"."modified" < (SELECT "modified" FROM "bookmarks" WHERE "id" = ?)`, args.after.String)
		}

		if args.order == OrderName {
			where.Or(`("child"."name" = (SELECT "name" from "bookmarks" WHERE "id" = ?) AND "child"."id" > ?))`, args.after.String, args.after.String)
		} else if args.order == OrderModified {
			where.Or(`("child"."modified" = (SELECT "modified" from "bookmarks" WHERE "id" = ?) AND "child"."id" > ?))`, args.after.String, args.after.String)
		}
	}

	books.Space("?", where)

	if args.direction == DirectionAsc && args.order == OrderName {
		books.Space(`ORDER BY "child"."name" ASC`)
	} else if args.direction == DirectionDesc && args.order == OrderName {
		books.Space(`ORDER BY "child"."name" DESC`)
	} else if args.direction == DirectionAsc && args.order == OrderModified {
		books.Space(`ORDER BY "child"."modified" ASC`)
	} else if args.direction == DirectionDesc && args.order == OrderModified {
		books.Space(`ORDER BY "child"."modified" DESC`)
	}

	if args.first.Dirty && args.first.Valid {
		books.Space(`LIMIT ?`, args.first.Int64)
	}

	results, err := query[bookDto](tx, books)
	return lo.Map(results, func(x bookDto, index int) Book {
		return Book{
			ID:          x.ID,
			URL:         PtrFromNullString(x.URL),
			Name:        x.Name,
			Description: PtrFromNullString(x.Description),
			ParentID:    PtrFromNullString(x.ParentID),
			IsFolder:    x.IsFolder,
			ParentName:  PtrFromNullString(x.ParentName),
			Tags:        parseTags(x.Tags),
		}
	}), err
}

// getTagsDBArgs are the args for getTagsDB.
type getTagsDBArgs struct {
	idFilter   NullString
	tagsFilter []string
	query      NullString
	after      NullString
	direction  Direction
	first      NullInt64
}

// getTagsDB lists tags in the bookmarks DB.
func getTagsDB(tx transaction, args getTagsDBArgs) ([]string, error) {
	tags := bqb.New(`SELECT "tag"`)
	tags.Space(`FROM "tags"`)

	where := bqb.Optional(`WHERE`)

	if len(args.tagsFilter) > 0 {
		where.And(`"tag" IN (?)`, args.tagsFilter)
	}

	if args.query.Dirty && args.query.Valid {
		searchFilter := bqb.New(`SELECT "tag"`)
		searchFilter.Space(`FROM "tags_fts"`)
		searchFilter.Space(`WHERE "tag" LIKE ?`, fmt.Sprintf("%%%s%%", args.query.String))

		where.And(`"tag" IN (?)`, searchFilter)
	}

	if args.after.Dirty && args.after.Valid {
		if args.direction == DirectionAsc {
			where.And(`"tag" > ?`, args.after.String)
		} else {
			where.And(`"tag" < ?`, args.after.String)
		}
	}

	tags.Space(`?`, where)

	if args.direction == DirectionAsc {
		tags.Space(`ORDER BY "tag" ASC, "id" ASC`)
	} else {
		tags.Space(`ORDER BY "tag" DESC, "id" ASC`)
	}

	if args.first.Dirty && args.first.Valid {
		tags.Space(`LIMIT ?`, args.first.Int64)
	}

	return query[string](tx, tags)
}

// bookFolderExistsDB returns true if the target book or folder exists.
func bookFolderExistsDB(tx transaction, ID string, isFolder bool) (bool, error) {
	books := bqb.New(`SELECT COUNT(1) AS "num"`)
	books.Space(`FROM "bookmarks"`)
	books.Space(`WHERE "bookmarks"."id" = ?`, ID)
	books.Space(`AND "bookmarks"."is_folder" = ?`, isFolder)

	count, err := count(tx, books)
	return count == 1, err
}

// getParentAndChildren gets a parent and all of its children.
func getParentAndChildren(tx transaction, ID string) ([]Book, error) {
	tags := bqb.New(`SELECT GROUP_CONCAT("tag")`)
	tags.Space(`FROM "bookmarks_tags"`)
	tags.Space(`WHERE "bookmark_id" = "child"."id"`)

	first := bqb.New(`SELECT "child"."id"`)
	first.Comma(`"child"."url"`)
	first.Comma(`"child"."name"`)
	first.Comma(`"child"."description"`)
	first.Comma(`"child"."parent_id"`)
	first.Comma(`"child"."is_folder"`)
	first.Comma(`"parent"."name" AS "parent"`)
	first.Comma(`IFNULL((?), '') AS "tags"`, tags)
	first.Space(`FROM "bookmarks" AS "child"`)
	first.Space(`LEFT JOIN "bookmarks" AS "parent" ON "parent"."id" = "child"."parent_id"`)
	first.Space(`WHERE "child"."id" = ?`, ID)

	rest := bqb.New(`SELECT "child"."id"`)
	rest.Comma(`"child"."url"`)
	rest.Comma(`"child"."name"`)
	rest.Comma(`"child"."description"`)
	rest.Comma(`"child"."parent_id"`)
	rest.Comma(`"child"."is_folder"`)
	rest.Comma(`"parent"."name" AS "parent"`)
	rest.Comma(`IFNULL((?), '') AS "tags"`, tags)
	rest.Space(`FROM "bookmarks" AS "child"`)
	rest.Space(`LEFT JOIN "bookmarks" AS "parent" ON "parent"."id" = "child"."parent_id"`)
	rest.Space(`INNER JOIN BOOK ON BOOK.id = "child"."parent_id"`)

	books := bqb.New(`WITH RECURSIVE BOOK AS (? UNION ALL ?)`, first, rest)
	books.Space(`SELECT "id"`)
	books.Comma(`"url"`)
	books.Comma(`"name"`)
	books.Comma(`"description"`)
	books.Comma(`"parent_id"`)
	books.Comma(`"is_folder"`)
	books.Comma(`"parent"`)
	books.Comma(`"tags"`)
	books.Space(`FROM BOOK`)

	results, err := query[bookDto](tx, books)
	return lo.Map(results, func(x bookDto, index int) Book {
		return Book{
			ID:          x.ID,
			URL:         PtrFromNullString(x.URL),
			Name:        x.Name,
			Description: PtrFromNullString(x.Description),
			ParentID:    PtrFromNullString(x.ParentID),
			IsFolder:    x.IsFolder,
			ParentName:  PtrFromNullString(x.ParentName),
			Tags:        parseTags(x.Tags),
		}
	}), err

}

// update

// updateBookDBArgs are the args for updateBookDB.
type updateBookDBArgs struct {
	name        NullString
	url         NullString
	description NullString
	parentID    NullString
}

// updateBookDB updates a book in the bookmarks database.
func updateBookDB(tx transaction, ID string, args updateBookDBArgs) error {
	update := bqb.New(`UPDATE "bookmarks"`)
	set := bqb.Optional(`SET`)

	if args.name.Dirty {
		set.Comma(`"name" = ?`, args.name)
	}

	if args.url.Dirty {
		set.Comma(`"url" = ?`, args.url)
	}

	if args.description.Dirty {
		set.Comma(`"description" = ?`, args.description)
	}

	if args.parentID.Dirty {
		set.Comma(`"parent_id" = ?`, args.parentID)
	}

	update.Space(`?`, set)
	update.Space(`WHERE "id" = ?`, ID)
	update.Space(`AND "is_folder" = ?`, false)

	return exec(tx, update)
}

// updateFolderDBArgs are the args for updateFolderDB.
type updateFolderDBArgs struct {
	name     NullString
	parentID NullString
}

// updateFolderDB updates a folder in the bookmarks database.
func updateFolderDB(tx transaction, ID string, args updateFolderDBArgs) error {
	update := bqb.New(`UPDATE "bookmarks"`)
	set := bqb.Optional(`SET`)

	if args.name.Dirty {
		set.Comma(`"name" = ?`, args.name)
	}

	if args.parentID.Dirty {
		set.Comma(`"parent_id" = ?`, args.parentID)
	}

	update.Space(`?`, set)
	update.Space(`WHERE "id" = ?`, ID)
	update.Space(`AND "is_folder" = ?`, true)

	return exec(tx, update)
}

// delete

// unlinkTagsDB removes tags from a bookmark.
func unlinkTagsDB(tx transaction, ID string, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	remove := bqb.New(`DELETE FROM "bookmarks_tags"`)
	remove.Space(`WHERE "bookmark_id" = ?`, ID)
	remove.Space(`AND "tag" IN (?)`, tags)

	return exec(tx, remove)
}

// removeBookDB deletes a bookmark from the bookmarks DB.
func removeBookDB(tx transaction, ID string) error {
	remove := bqb.New(`DELETE FROM "bookmarks"`)
	remove.Space(`WHERE id = ?`, ID)
	remove.Space(`AND is_folder = ?`, false)

	return exec(tx, remove)
}

// removeFolderDB deletes a folder from the bookmarks DB.
func removeFolderDB(tx transaction, ID string) error {
	remove := bqb.New(`DELETE FROM "bookmarks"`)
	remove.Space(`WHERE id = ?`, ID)
	remove.Space(`AND is_folder = ?`, true)

	return exec(tx, remove)
}

// cleanOrphanedTagsDB removes any tags that aren't applied to a bookmark.
func cleanOrphanedTagsDB(tx transaction, tags []string) error {
	existing := bqb.New(`SELECT 1`)
	existing.Space(`FROM "bookmarks_tags"`)
	existing.Space(`WHERE "bookmarks_tags"."tag" = "tags"."tag"`)

	remove := bqb.New(`DELETE FROM "tags"`)
	remove.Space(`WHERE "tag" IN (?)`, tags)
	remove.Space(`AND NOT EXISTS (?)`, existing)

	return exec(tx, remove)
}

// helpers

// parseTags parses the tags coming back from the database.
func parseTags(tags string) []string {
	if tags == "" {
		return make([]string, 0)
	}

	return strings.Split(tags, ",")
}

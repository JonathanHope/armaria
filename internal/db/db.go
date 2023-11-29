package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/pkg/model"
	"github.com/nullism/bqb"
	"github.com/samber/lo"
)

// This file contains the low level logic to access the bookmarks DB.

// create

// AddBook inserts a book into the bookmarks database.
func AddBook(tx Transaction, url string, name string, description null.NullString, parentID null.NullString) (string, error) {
	id := uuid.New().String()

	insert := bqb.New(`INSERT INTO "bookmarks"("id", "url", "is_folder", "name", "description", "parent_id")`)
	insert.Space("VALUES(?, ?, ?, ?, ?, ?)", id, url, false, name, description, parentID)

	err := exec(tx, insert)
	return id, err
}

// AddFolder inserts a folder into the bookmarks database.
func AddFolder(tx Transaction, name string, parentID null.NullString) (string, error) {
	id := uuid.New().String()

	insert := bqb.New(`INSERT INTO "bookmarks"("id", "is_folder", "name", "parent_id")`)
	insert.Space("VALUES(?, ?, ?, ?)", id, true, name, parentID)

	err := exec(tx, insert)
	return id, err
}

// AddTags inserts tags into the bookmarks database.
func AddTags(tx Transaction, tags []string) error {
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

// LinkTags adds tags to bookmark.
func LinkTags(tx Transaction, bookmarkID string, tags []string) error {
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

// GetBooksArgs are the args for getBooksDB.
type GetBooksArgs struct {
	IDFilter       string
	IncludeBooks   bool
	IncludeFolders bool
	ParentID       null.NullString
	Query          null.NullString
	Tags           []string
	After          null.NullString
	Order          armaria.Order
	Direction      armaria.Direction
	First          null.NullInt64
}

// GetBooks lists bookmarks/folders in the bookmarks DB.
func GetBooks(tx Transaction, args GetBooksArgs) ([]armaria.Book, error) {
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

	if args.IDFilter != "" {
		where.And(`"child"."id" = ?`, args.IDFilter)
	}

	if args.IncludeBooks && !args.IncludeFolders {
		where.And(`"child"."is_folder" = ?`, false)
	}

	if args.IncludeFolders && !args.IncludeBooks {
		where.And(`"child"."is_folder" = ?`, true)
	}

	if args.ParentID.Dirty && args.ParentID.Valid {
		where.And(`"child"."parent_id" = ?`, args.ParentID.String)
	} else if args.ParentID.Dirty && !args.ParentID.Valid {
		where.And(`"child"."parent_id" IS NULL`)
	}

	if args.Query.Dirty && args.Query.Valid {
		searchFilter := bqb.New(`SELECT "id"`)
		searchFilter.Space(`FROM "bookmarks_fts"`)
		searchFilter.Space(`WHERE "name" LIKE ?`, fmt.Sprintf("%%%s%%", args.Query.String))
		searchFilter.Space(`OR "description" LIKE ?`, fmt.Sprintf("%%%s%%", args.Query.String))
		searchFilter.Space(`OR "url" LIKE ?`, fmt.Sprintf("%%%s%%", args.Query.String))

		where.And(`"child"."id" IN (?)`, searchFilter)
	}

	if len(args.Tags) > 0 {
		tagsFilter := bqb.New(`SELECT "bookmark_id"`)
		tagsFilter.Space(`FROM "bookmarks_tags"`)
		tagsFilter.Space(`WHERE "tag" IN (?)`, args.Tags)

		where.And(`"child"."id" IN (?)`, tagsFilter)
	}

	if args.After.Dirty && args.After.Valid {
		if args.Order == armaria.OrderName && args.Direction == armaria.DirectionAsc {
			where.And(`("child"."name" > (SELECT "name" FROM "bookmarks" WHERE "id" = ?)`, args.After.String)
		} else if args.Order == armaria.OrderName && args.Direction == armaria.DirectionDesc {
			where.And(`("child"."name" < (SELECT "name" FROM "bookmarks" WHERE "id" = ?)`, args.After.String)
		} else if args.Order == armaria.OrderModified && args.Direction == armaria.DirectionAsc {
			where.And(`("child"."modified" > (SELECT "modified" FROM "bookmarks" WHERE "id" = ?)`, args.After.String)
		} else if args.Order == armaria.OrderModified && args.Direction == armaria.DirectionDesc {
			where.And(`("child"."modified" < (SELECT "modified" FROM "bookmarks" WHERE "id" = ?)`, args.After.String)
		}

		if args.Order == armaria.OrderName {
			where.Or(`("child"."name" = (SELECT "name" from "bookmarks" WHERE "id" = ?) AND "child"."id" > ?))`, args.After.String, args.After.String)
		} else if args.Order == armaria.OrderModified {
			where.Or(`("child"."modified" = (SELECT "modified" from "bookmarks" WHERE "id" = ?) AND "child"."id" > ?))`, args.After.String, args.After.String)
		}
	}

	books.Space("?", where)

	if args.Direction == armaria.DirectionAsc && args.Order == armaria.OrderName {
		books.Space(`ORDER BY "child"."name" ASC`)
	} else if args.Direction == armaria.DirectionDesc && args.Order == armaria.OrderName {
		books.Space(`ORDER BY "child"."name" DESC`)
	} else if args.Direction == armaria.DirectionAsc && args.Order == armaria.OrderModified {
		books.Space(`ORDER BY "child"."modified" ASC`)
	} else if args.Direction == armaria.DirectionDesc && args.Order == armaria.OrderModified {
		books.Space(`ORDER BY "child"."modified" DESC`)
	}

	if args.First.Dirty && args.First.Valid {
		books.Space(`LIMIT ?`, args.First.Int64)
	}

	results, err := query[bookDTO](tx, books)
	return lo.Map(results, func(x bookDTO, index int) armaria.Book {
		return x.toBook()
	}), err
}

// GetTagsArgs are the args for getTagsDB.
type GetTagsArgs struct {
	IDFilter   null.NullString
	TagsFilter []string
	Query      null.NullString
	After      null.NullString
	Direction  armaria.Direction
	First      null.NullInt64
}

// GetTags lists tags in the bookmarks DB.
func GetTags(tx Transaction, args GetTagsArgs) ([]string, error) {
	tags := bqb.New(`SELECT "tag"`)
	tags.Space(`FROM "tags"`)

	where := bqb.Optional(`WHERE`)

	if len(args.TagsFilter) > 0 {
		where.And(`"tag" IN (?)`, args.TagsFilter)
	}

	if args.Query.Dirty && args.Query.Valid {
		searchFilter := bqb.New(`SELECT "tag"`)
		searchFilter.Space(`FROM "tags_fts"`)
		searchFilter.Space(`WHERE "tag" LIKE ?`, fmt.Sprintf("%%%s%%", args.Query.String))

		where.And(`"tag" IN (?)`, searchFilter)
	}

	if args.After.Dirty && args.After.Valid {
		if args.Direction == armaria.DirectionAsc {
			where.And(`"tag" > ?`, args.After.String)
		} else {
			where.And(`"tag" < ?`, args.After.String)
		}
	}

	tags.Space(`?`, where)

	if args.Direction == armaria.DirectionAsc {
		tags.Space(`ORDER BY "tag" ASC, "id" ASC`)
	} else {
		tags.Space(`ORDER BY "tag" DESC, "id" ASC`)
	}

	if args.First.Dirty && args.First.Valid {
		tags.Space(`LIMIT ?`, args.First.Int64)
	}

	return query[string](tx, tags)
}

// BookFolderExists returns true if the target book or folder exists.
func BookFolderExists(tx Transaction, ID string, isFolder bool) (bool, error) {
	books := bqb.New(`SELECT COUNT(1) AS "num"`)
	books.Space(`FROM "bookmarks"`)
	books.Space(`WHERE "bookmarks"."id" = ?`, ID)
	books.Space(`AND "bookmarks"."is_folder" = ?`, isFolder)

	count, err := count(tx, books)
	return count == 1, err
}

// GetParentAndChildren gets a parent and all of its children.
func GetParentAndChildren(tx Transaction, ID string) ([]armaria.Book, error) {
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

	results, err := query[bookDTO](tx, books)
	return lo.Map(results, func(x bookDTO, index int) armaria.Book {
		return x.toBook()
	}), err

}

// update

// UpdateBookArgs are the args for updateBookDB.
type UpdateBookArgs struct {
	Name        null.NullString
	URL         null.NullString
	Description null.NullString
	ParentID    null.NullString
}

// UpdateBook updates a book in the bookmarks database.
func UpdateBook(tx Transaction, ID string, args UpdateBookArgs) error {
	update := bqb.New(`UPDATE "bookmarks"`)
	set := bqb.Optional(`SET`)

	if args.Name.Dirty {
		set.Comma(`"name" = ?`, args.Name)
	}

	if args.URL.Dirty {
		set.Comma(`"url" = ?`, args.URL)
	}

	if args.Description.Dirty {
		set.Comma(`"description" = ?`, args.Description)
	}

	if args.ParentID.Dirty {
		set.Comma(`"parent_id" = ?`, args.ParentID)
	}

	update.Space(`?`, set)
	update.Space(`WHERE "id" = ?`, ID)
	update.Space(`AND "is_folder" = ?`, false)

	return exec(tx, update)
}

// UpdateFolderArgs are the args for updateFolderDB.
type UpdateFolderArgs struct {
	Name     null.NullString
	ParentID null.NullString
}

// UpdateFolder updates a folder in the bookmarks database.
func UpdateFolder(tx Transaction, ID string, args UpdateFolderArgs) error {
	update := bqb.New(`UPDATE "bookmarks"`)
	set := bqb.Optional(`SET`)

	if args.Name.Dirty {
		set.Comma(`"name" = ?`, args.Name)
	}

	if args.ParentID.Dirty {
		set.Comma(`"parent_id" = ?`, args.ParentID)
	}

	update.Space(`?`, set)
	update.Space(`WHERE "id" = ?`, ID)
	update.Space(`AND "is_folder" = ?`, true)

	return exec(tx, update)
}

// delete

// UnlinkTags removes tags from a bookmark.
func UnlinkTags(tx Transaction, ID string, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	remove := bqb.New(`DELETE FROM "bookmarks_tags"`)
	remove.Space(`WHERE "bookmark_id" = ?`, ID)
	remove.Space(`AND "tag" IN (?)`, tags)

	return exec(tx, remove)
}

// RemoveBook deletes a bookmark from the bookmarks DB.
func RemoveBook(tx Transaction, ID string) error {
	remove := bqb.New(`DELETE FROM "bookmarks"`)
	remove.Space(`WHERE id = ?`, ID)
	remove.Space(`AND is_folder = ?`, false)

	return exec(tx, remove)
}

// RemoveFolder deletes a folder from the bookmarks DB.
func RemoveFolder(tx Transaction, ID string) error {
	remove := bqb.New(`DELETE FROM "bookmarks"`)
	remove.Space(`WHERE id = ?`, ID)
	remove.Space(`AND is_folder = ?`, true)

	return exec(tx, remove)
}

// CleanOrphanedTags removes any tags that aren't applied to a bookmark.
func CleanOrphanedTags(tx Transaction, tags []string) error {
	existing := bqb.New(`SELECT 1`)
	existing.Space(`FROM "bookmarks_tags"`)
	existing.Space(`WHERE "bookmarks_tags"."tag" = "tags"."tag"`)

	remove := bqb.New(`DELETE FROM "tags"`)
	remove.Space(`WHERE "tag" IN (?)`, tags)
	remove.Space(`AND NOT EXISTS (?)`, existing)

	return exec(tx, remove)
}

-- +goose Up
-- +goose StatementBegin
CREATE TABLE "bookmarks" (
  "id" TEXT PRIMARY KEY NOT NULL,
  "parent_id" TEXT NULL,
  "is_folder" INTEGER NOT NULL,
  "name" TEXT NOT NULL CHECK(LENGTH("name") <= 2048),
  "url" TEXT CHECK("url" IS NULL OR LENGTH("url") <= 2048),
  "description" TEXT CHECK("description" IS NULL OR LENGTH("description") <= 4096),
  "modified" TEXT NOT NULL DEFAULT(datetime()),
  FOREIGN KEY("parent_id") REFERENCES "bookmarks"("id"),
  CHECK ("is_folder" IN (0, 1)),
  CHECK (("is_folder" = 0 AND "url" IS NOT NULL) OR ("is_folder" = 1 AND "url" IS NULL))
) STRICT;

CREATE TABLE "tags" (
  "tag" TEXT PRIMARY KEY NOT NULL CHECK(LENGTH("tag") <= 128),
  "modified" TEXT NOT NULL DEFAULT(datetime())
) STRICT;

CREATE TABLE "bookmarks_tags" (
  "bookmark_id" TEXT NOT NULL,
  "tag" TEXT NOT NULL,
  "modified" TEXT NOT NULL DEFAULT(datetime()),
  PRIMARY KEY ("bookmark_id", "tag"),
  FOREIGN KEY("bookmark_id") REFERENCES "bookmarks"("id"),
  FOREIGN KEY("tag") REFERENCES "tags"("tag")
) STRICT;

CREATE VIRTUAL TABLE "bookmarks_fts"
USING fts5("id" UNINDEXED, "name", "url", "description", tokenize="trigram");

CREATE INDEX "ix_bookmarks_parent_id"
ON "bookmarks"("parent_id");

CREATE INDEX "ix_bookmarks_is_folder"
ON "bookmarks"("is_folder");

CREATE INDEX "ix_bookmarks_tags_bookmark_id"
ON "bookmarks_tags"("bookmark_id");

CREATE INDEX "ix_bookmarks_tags_tag"
ON "bookmarks_tags"("tag");

CREATE INDEX "ix_bookmarks_name"
ON "bookmarks"("name");

CREATE INDEX "ix_bookmarks_modified"
ON "bookmarks"("modifed");

CREATE TRIGGER "after_bookmarks_insert" AFTER INSERT ON "bookmarks" BEGIN
  INSERT INTO bookmarks_fts (
    "id",
    "name",
    "url",
    "description"
  )
  VALUES(
    new."id",
    new."name",
    new."url",
    new."description"
  );
END;

CREATE TRIGGER "after_bookmarks_update" UPDATE ON "bookmarks" BEGIN
  UPDATE "bookmarks_fts"
   SET "name" = new."name",
   "url" = new."url",
   "description" = new."description"
  WHERE "id" = old."id";
END;

CREATE TRIGGER "after_bookmarks_delete" AFTER DELETE ON "bookmarks" BEGIN
  DELETE FROM "bookmarks_fts"
  WHERE "id" = old."id";
END;

CREATE VIRTUAL TABLE "tags_fts"
USING fts5("tag");

CREATE TRIGGER "after_tags_insert" AFTER INSERT ON "tags" BEGIN
  INSERT INTO tags_fts (
    "tag"
  )
  VALUES(
    new."tag"
  );
END;

CREATE TRIGGER "after_tags_delete" AFTER DELETE ON "tags" BEGIN
  DELETE FROM "tags_fts"
  WHERE "tag" = old."tag";
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "ix_bookmarks_parent_id";

DROP INDEX "ix_bookmarks_is_folder";

DROP INDEX ix_bookmarks_tags_bookmark_id;

DROP INDEX "ix_bookmarks_tags_tag";

DROP INDEX "ix_bookmarks_name";

DROP INDEX "ix_bookmarks_modified";

DROP TRIGGER "after_bookmarks_insert";

DROP TRIGGER "after_bookmarks_update";

DROP TRIGGER "after_bookmarks_delete";

DROP TABLE "bookmarks_tags";

DROP TABLE "bookmarks_fts";

DROP TABLE "tags_fts";

DROP TABLE "bookmarks";

DROP TABLE "tags";
-- +goose StatementEnd

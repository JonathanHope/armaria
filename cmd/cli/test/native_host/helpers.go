package test

import (
	"bytes"
	"database/sql"

	"github.com/blockloop/scan/v2"
	"github.com/jonathanhope/armaria/cmd/cli/internal/messaging"
	"github.com/nullism/bqb"
)

// nativeMessageLoop performs a full loop through native messaging.
// It ends a native message to the host and then returns the response.
func nativeMessageLoop[T messaging.Payload](kind messaging.MessageKind, payload T) (messaging.NativeMessage, error) {
	msg, err := messaging.PayloadToMessage(kind, payload)
	if err != nil {
		return messaging.NativeMessage{}, err
	}

	in := bytes.NewBuffer(nil)
	err = msg.SendMessage(in)
	if err != nil {
		return messaging.NativeMessage{}, err
	}

	out := bytes.NewBuffer(nil)
	if err := messaging.Dispatch(in, out); err != nil {
		return messaging.NativeMessage{}, err
	}

	res, err := messaging.ReceiveMessage(out)
	if err != nil {
		return messaging.NativeMessage{}, err
	}

	return res, err
}

// getLastInsertedID gets the ID of the last inserted bookmark.
func getLastInsertedID(dbLocation string, ignoreIds []string) (string, error) {
	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := bqb.New(`SELECT "id" FROM "bookmarks"`)
	where := bqb.Optional("WHERE")
	if len(ignoreIds) > 0 {
		where.And(`"id" NOT IN (?)`, ignoreIds)
	}
	query.Space(`?`, where)
	query.Space(`ORDER BY "modified" DESC LIMIT 1`)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	rows, err := db.Query(sql, args...)
	if err != nil {
		return "", err
	}

	ids := make([]string, 0)
	err = scan.RowsStrict(&ids, rows)
	if err != nil {
		return "", err
	}

	return ids[0], err
}

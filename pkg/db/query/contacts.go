package query

import (
	"context"

	"github.com/connorvoisey/shgrid_api/.gen/todo/public/table"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	. "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func ListContacts(log zerolog.Logger, context context.Context, req ContactReq, res ContactRes) {
	log.Debug().Any("inputBody", input.Body).Int64("limit", int64(input.Body.Limit)).Msg("/contacts")
	resp := &ContactRes{}
	stmt := SELECT(
		table.Contacts.AllColumns,
	).FROM(
		table.Contacts,
	)

	contactsAddWhere(stmt, input.Body.Filters)
	stmt.LIMIT(int64(input.Body.Limit)).OFFSET(int64(input.Body.Offset))

	sql, params := stmt.Sql()
	log.Debug().Str("sql", sql).Any("params", params).Msg("/contacts query")

	err := stmt.Query(
		db,
		&resp.Body.Data,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch contacts from database")
	}
	if resp.Body.Data == nil {
		resp.Body.Data = make([]Contacts, 0)
	}

	var dest struct {
		Count int
	}
	stmtQuery := SELECT(COUNT(STAR)).FROM(table.Contacts)
	contactsAddWhere(stmtQuery, input.Body.Filters)
	countErr := stmtQuery.Query(
		db,
		&dest,
	)
	sql, params = stmtQuery.Sql()
	log.Debug().Str("sql", sql).Any("params", params).Any("dest", dest).Msg("/contacts query")
	if err != nil {
		log.Error().Err(countErr).Msg("Failed to fetch contacts count from database")
	}
	resp.Body.Count = dest.Count
	return resp, nil
}

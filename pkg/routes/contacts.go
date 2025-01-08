package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/connorvoisey/shgrid_api/.gen/todo/public/table"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type ContactRes struct {
	Body struct {
		Rows  []ContactsWithOrganisation `json:"rows"`
		Count int                        `json:"count"`
	}
}

type ContactFilter struct {
	ID        *uuid.UUID `sql:"primary_key"  json:"id,omitempty" doc:"Id of record"`
	FirstName *string    `json:"first_name,omitempty" maxLength:"255" doc:"First Name of contact"`
	LastName  *string    `json:"last_name,omitempty" maxLength:"255" doc:"Last Name of contact"`
	Email     *string    `json:"email,omitempty" maxLength:"255" doc:"Email of contact"`
}

type ContactSorter struct {
	Key string `json:"key" enum:"id,first_name,last_name,email"`
	Dir string `json:"dir" enum:"asc,desc"`
}

type ContactReq struct {
	Body struct {
		Limit   int             `minimum:"0" maximum:"100" json:"limit" example:"5"`
		Offset  int             `minimum:"0" json:"offset"`
		Filters ContactFilter   `json:"filters"`
		Sorters []ContactSorter `json:"sorters"`
	}
}
type Contacts struct {
	ID        uuid.UUID `sql:"primary_key" json:"id" doc:"Id of record"`
	FirstName *string   `json:"first_name" maxLength:"255" doc:"First Name of contact"`
	LastName  string    `json:"last_name" maxLength:"255" doc:"Last Name of contact"`
	Email     string    `json:"email" maxLength:"255" doc:"Email of contact"`
}

type ContactsWithOrganisation struct {
	Contacts
	Organisation *Organisations `json:"organisation"`
	Address      *Addresses     `json:"address"`
}

func contactsAddWhere(query SelectStatement, filter ContactFilter) {
	condition := Bool(true)
	if filter.ID != nil {
		condition = condition.AND(table.Contacts.ID.EQ(String(filter.ID.String())))
	}
	if filter.FirstName != nil {
		condition = condition.AND(table.Contacts.FirstName.LIKE(String(fmt.Sprintf("%%%s%%", *filter.FirstName))))
	}
	if filter.LastName != nil {
		condition = condition.AND(table.Contacts.LastName.LIKE(String(fmt.Sprintf("%%%s%%", *filter.LastName))))
	}
	if filter.Email != nil {
		condition = condition.AND(table.Contacts.Email.LIKE(String(fmt.Sprintf("%%%s%%", *filter.Email))))
	}
	query.WHERE(condition)
}

func getContactsSortableColumn(key string) *ColumnString {
	switch key {
	case "id":
		return &table.Contacts.ID
	case "first_name":
		return &table.Contacts.FirstName
	case "last_name":
		return &table.Contacts.LastName
	case "email":
		return &table.Contacts.Email
	}
	return nil
}

func contactsAddSorting(query SelectStatement, sorters []ContactSorter) {
	var orderBys []OrderByClause
	for _, sorter := range sorters {
		col := getContactsSortableColumn(sorter.Key)
		orderBy := getOrderBy(sorter.Dir, *col)
		orderBys = append(orderBys, orderBy)
	}
	query.ORDER_BY(orderBys...)
}

func listContacts(ctx context.Context, input *ContactReq, log *zerolog.Logger, db *sql.DB) (*ContactRes, error) {
    // time.Sleep(time.Duration(time.Second * 2))
	log.Debug().Any("inputBody", input.Body).Int64("limit", int64(input.Body.Limit)).Msg("/contacts")

	// orgAddress := table.Addresses.AS("OrgAddress")
	resp := &ContactRes{}
	stmt := SELECT(
		table.Contacts.AllColumns,
		// table.Organisations.AllColumns,
		// table.Addresses.AllColumns,
		// orgAddress.AllColumns,
	).FROM(
		table.Contacts,
		// LEFT_JOIN(table.Organisations, table.Contacts.OrganisationID.EQ(table.Organisations.ID)).
		// LEFT_JOIN(table.Addresses, table.Contacts.AddressID.EQ(table.Addresses.ID)).
		// LEFT_JOIN(orgAddress, table.Organisations.AddressID.EQ(orgAddress.ID)),
	)

	contactsAddWhere(stmt, input.Body.Filters)
	contactsAddSorting(stmt, input.Body.Sorters)
	stmt.LIMIT(int64(input.Body.Limit)).OFFSET(int64(input.Body.Offset))

	sql, params := stmt.Sql()
	log.Debug().Str("sql", sql).Any("params", params).Msg("/contacts query")

	err := stmt.Query(
		db,
		&resp.Body.Rows,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch contacts from database")
		return nil, huma.Error500InternalServerError("Failed to fetch contacts from database")
	}
	if resp.Body.Rows == nil {
		resp.Body.Rows = make([]ContactsWithOrganisation, 0)
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
	countSql, countParams := stmtQuery.Sql()
	log.Debug().Str("sql", countSql).Any("params", countParams).Any("dest", dest).Msg("/contacts query")
	if countErr != nil {
		log.Error().Err(countErr).Msg("Failed to fetch contacts count from database")
		return nil, huma.Error500InternalServerError("Failed to fetch contacts from database")
	}
	resp.Body.Count = dest.Count
	return resp, nil
}

func AddContactRoutes(api *huma.API, log *zerolog.Logger, db *sql.DB) {
	huma.Register(*api, huma.Operation{
		OperationID:   "list-contacts",
		Method:        http.MethodPost,
		Path:          "/contacts",
		Summary:       "List contacts",
		Tags:          []string{"Contacts"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *ContactReq) (*ContactRes, error) {
		return listContacts(ctx, input, log, db)
	})
}

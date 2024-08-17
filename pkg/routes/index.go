package routes

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/connorvoisey/shgrid_api/.gen/todo/public/table"
	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"net/http"
)

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

// ReviewInput represents the review operation request.
type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Author of the review"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message"`
	}
}
type Review struct {
	Body struct {
		Id      int    `json:"id" minimum:"1" doc:"Id of the review"`
		Author  string `json:"author" maxLength:"10" doc:"Author of the review"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message"`
	}
}

type ContactRes struct {
	Body struct {
		Data  []Contacts `json:"data"`
		Count int        `json:"count"`
	}
}

type ContactFilter struct {
	ID        *uuid.UUID `sql:"primary_key"  json:"id,omitempty" doc:"Id of record"`
	FirstName *string    `json:"first_name,omitempty" maxLength:"255" doc:"First Name of contact"`
	LastName  *string    `json:"last_name,omitempty" maxLength:"255" doc:"Last Name of contact"`
	Email     *string    `json:"email,omitempty" maxLength:"255" doc:"Email of contact"`
}

type ContactSorter struct {
	Field string `json:"field" enum:"first_name,last_name,email"`
	Dir   string `json:"dir" enum:"asc,desc"`
}

type ContactReq struct {
	Body struct {
		Limit   int             `minimum:"0" maximum:"100" json:"limit" example:"100"`
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

type Addresses struct {
	ID           uuid.UUID `sql:"primary_key" json:"id" doc:"Id of record"`
	AddressLine1 string    `json:"address_line_one" maxLength:"255"`
	AddressLine2 string    `json:"address_line_two" maxLength:"255"`
	TownOrCity   string    `json:"town_or_city" maxLength:"255"`
	Country      string    `json:"country" maxLength:"255"`
	Postcode     string    `json:"postcode" maxLength:"255"`
}

type Organisations struct {
	ID   uuid.UUID `sql:"primary_key" json:"id" doc:"Id of record"`
	Name string    `json:"name" maxLength:"255"`
}

type ContactsWithOrganisation struct {
	Contacts
	Organisation Organisations `json:"organisation"`
}
type CountDest struct {
	Count int `sql:"primary_key"`
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
	case "first_name":
		return &table.Contacts.FirstName
	case "last_name":
		return &table.Contacts.LastName
	case "email":
		return &table.Contacts.Email
	}
	return nil
}

func getOrderBy(dirString string, col ColumnString) OrderByClause {
	switch dirString {
	case "asc":
		return col.ASC()
	case "desc":
		return col.DESC()
	}
	return nil
}

func contactsAddSorting(query SelectStatement, sorters []ContactSorter) {
	var orderBys []OrderByClause
	for _, sorter := range sorters {
		col := getContactsSortableColumn(sorter.Field)
		orderBys = append(orderBys, getOrderBy(sorter.Dir, *col))
	}
	query.ORDER_BY(orderBys...)
}

func AddRoutes(api huma.API, log *zerolog.Logger, db *sql.DB) {
	// Register GET /greeting/{name} handler.
	huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	// Register POST /reviews
	huma.Register(api, huma.Operation{
		OperationID:   "list-contacts",
		Method:        http.MethodPost,
		Path:          "/contacts",
		Summary:       "List contacts",
		Tags:          []string{"Contacts"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *ContactReq) (*ContactRes, error) {
		log.Debug().Any("inputBody", input.Body).Int64("limit", int64(input.Body.Limit)).Msg("/contacts")
		resp := &ContactRes{}
		stmt := SELECT(
			table.Contacts.AllColumns,
		).FROM(
			table.Contacts,
		)

		contactsAddWhere(stmt, input.Body.Filters)
		contactsAddSorting(stmt, input.Body.Sorters)
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
	})

	// Register POST /reviews
	huma.Register(api, huma.Operation{
		OperationID:   "post-review",
		Method:        http.MethodPost,
		Path:          "/reviews",
		Summary:       "Post a review",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *ReviewInput) (*Review, error) {
		// TODO: save review in data store.
		resp := &Review{}
		resp.Body.Id = 1
		resp.Body.Author = i.Body.Author
		resp.Body.Rating = i.Body.Rating
		resp.Body.Message = i.Body.Message
		return resp, nil
	})

}

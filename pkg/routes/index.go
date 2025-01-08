package routes

import (
	"context"
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Addresses struct {
	ID           uuid.UUID `sql:"primary_key" json:"id" doc:"Id of record"`
	AddressLine1 string    `json:"address_line_one" maxLength:"255"`
	AddressLine2 string    `json:"address_line_two" maxLength:"255"`
	TownOrCity   string    `json:"town_or_city" maxLength:"255"`
	Country      string    `json:"country" maxLength:"255"`
	Postcode     string    `json:"postcode" maxLength:"255"`
}

type Organisations struct {
	ID      uuid.UUID `sql:"primary_key" json:"id" doc:"Id of record"`
	Name    string    `json:"name" maxLength:"255"`
	Address Addresses `json:"address" alias:"OrgAddress"`
}
type CountDest struct {
	Count int `sql:"primary_key"`
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

type HtmlRes struct {
	Body []byte
}

func AddRoutes(api *huma.API, log *zerolog.Logger, db *sql.DB) {
	huma.Get(*api, "/docs", func(ctx context.Context, _ *struct{}) (*HtmlRes, error) {
		return &HtmlRes{Body: []byte(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`)}, nil
	})
	AddContactRoutes(api, log, db)
}

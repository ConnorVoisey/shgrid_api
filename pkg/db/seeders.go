package db

import (
	"database/sql"
	"fmt"
	"github.com/connorvoisey/shgrid_api/.gen/todo/public/table"
	_ "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
	_ "github.com/lib/pq"
	"math/rand"
	"sync"
)

type (
	Contacts struct {
		ID             uuid.UUID  `sql:"primary_key" json:"id"`
		FirstName      *string    `json:"first_name"`
		LastName       string     `json:"last_name"`
		Email          string     `json:"email"`
		OrganisationID *uuid.UUID `json:"organisation_id"`
		AddressID      *uuid.UUID `json:"address_id"`
	}

	Organisations struct {
		ID        uuid.UUID  `sql:"primary_key" json:"id"`
		Name      string     `json:"name"`
		AddressID *uuid.UUID `json:"address_id"`
	}
	Addresses struct {
		ID           uuid.UUID `sql:"primary_key" json:"id"`
		AddressLine1 *string   `json:"address_line_1"`
		AddressLine2 *string   `json:"address_line_2"`
		TownOrCity   string    `json:"town_or_city"`
		Country      string    `json:"country"`
		Postcode     string    `json:"postcode"`
	}
)

func Seed(db *sql.DB, addCount int, orgCount int, conCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	addresses, err := createAddresses(addCount, db)
	panicOnError(err, "Failed to create addresses")
	organisations, err := createOrganisations(orgCount, db, addresses)
	panicOnError(err, "Failed to create organisations")
	contacts, err := createContacts(conCount, db, addresses, organisations)
	panicOnError(err, "Failed to create contacts")
	fmt.Printf("Created %v addresses, %v organisations, %v contacts\n", len(addresses), len(organisations), len(contacts))
}

func createContacts(count int, db *sql.DB, addresses []Addresses, organisations []Organisations) ([]Contacts, error) {
	fake := faker.New()
	var contactFactories []Contacts
	for i := 0; i < count; i++ {
		var randAddId *uuid.UUID = nil
		if rand.Intn(10) < 8 {
			randAddId = &addresses[rand.Intn(len(addresses))].ID
		}

		var randOrgId *uuid.UUID = nil
		if rand.Intn(10) < 8 {
			randOrgId = &organisations[rand.Intn(len(organisations))].ID
		}
		contactInput, err := ContactFactory(fake, randAddId, randOrgId)
		if err != nil {
			return nil, err
		}
		contactFactories = append(contactFactories, *contactInput)
	}
	query := table.Contacts.INSERT(
		table.Contacts.ID,
		table.Contacts.FirstName,
		table.Contacts.LastName,
		table.Contacts.Email,
		table.Contacts.OrganisationID,
		table.Contacts.AddressID,
	).MODELS(contactFactories).RETURNING(table.Contacts.AllColumns)

	// sql, params := query.Sql()
	// fmt.Println(fmt.Sprintf("SQL: %s, params: %v", sql, params))

	var dest []Contacts

	err := query.Query(db, &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}
func createAddresses(count int, db *sql.DB) ([]Addresses, error) {
	fake := faker.New()
	var addressFactories []Addresses
	for i := 0; i < count; i++ {
		addInput, err := AddressFactory(fake)
		if err != nil {
			return nil, err
		}
		addressFactories = append(addressFactories, *addInput)
	}
	query := table.Addresses.INSERT(
		table.Addresses.ID,
		table.Addresses.AddressLine1,
		table.Addresses.AddressLine2,
		table.Addresses.TownOrCity,
		table.Addresses.Country,
		table.Addresses.Postcode,
	).MODELS(addressFactories).RETURNING(table.Addresses.AllColumns)

	// sql, params := query.Sql()
	// fmt.Println(fmt.Sprintf("SQL: %s, params: %v", sql, params))

	var dest []Addresses

	err := query.Query(db, &dest)
	return dest, err
}

func createOrganisations(count int, db *sql.DB, addresses []Addresses) ([]Organisations, error) {
	fake := faker.New()
	var orgFactories []Organisations
	for i := 0; i < count; i++ {
		var addId *uuid.UUID = nil
		if rand.Intn(10) < 9 {
			randAddressI := rand.Intn(len(addresses) - 1)
			addId = &addresses[randAddressI].ID
		}
		newOrg, err := OrganisationFactory(fake, addId)
		if err != nil {
			return nil, err
		}
		orgFactories = append(orgFactories, *newOrg)
	}
	query := table.Organisations.INSERT(
		table.Organisations.ID,
		table.Organisations.Name,
		table.Organisations.AddressID,
	).MODELS(orgFactories).RETURNING(table.Organisations.AllColumns)

	var dest []Organisations

	err := query.Query(db, &dest)

	return dest, err
}
func panicOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

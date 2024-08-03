package db

import (
	_ "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
	_ "github.com/lib/pq"
)

func ContactFactory(fake faker.Faker, addressId *uuid.UUID, organisationId *uuid.UUID) (*Contacts, error) {
	person := fake.Person()
	contact := person.Contact()
	firstName := person.FirstName()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Contacts{
		ID:             id,
		FirstName:      &firstName,
		LastName:       person.LastName(),
		Email:          contact.Email,
		OrganisationID: organisationId,
		AddressID:      addressId,
	}, nil
}

func AddressFactory(fake faker.Faker) (*Addresses, error) {
	address := fake.Address()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	lineOne := address.StreetName()
	lineTwo := address.State()
	return &Addresses{
		ID:           id,
		AddressLine1: &lineOne,
		AddressLine2: &lineTwo,
		TownOrCity:   address.City(),
		Country:      address.Country(),
		Postcode:     address.PostCode(),
	}, nil
}
func OrganisationFactory(fake faker.Faker, addressId *uuid.UUID) (*Organisations, error) {
	company := fake.Company()
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &Organisations{
		ID:        id,
		Name:      company.Name(),
		AddressID: addressId,
	}, err
}

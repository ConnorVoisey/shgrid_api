//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
)

type Contacts struct {
	ID             uuid.UUID `sql:"primary_key"`
	FirstName      *string
	LastName       string
	Email          string
	OrganisationID *uuid.UUID
	AddressID      *uuid.UUID
}

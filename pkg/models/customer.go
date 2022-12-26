package models

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"

// Customer stream demo
// @def primary                  ID
//
//go:generate toolkit gen model Customer --database DB
type Customer struct {
	ID        string `json:"id"        db:"id"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName"  db:"last_name"`
	Age       int    `json:"age"       db:"age"`
	TaxNumber string `json:"taxNumber" db:"tax_number,default=''"`
	City      string `json:"city"   db:"city"`
	datatypes.OperationTimesWithDeleted
}

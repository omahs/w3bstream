package models

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"

// SourceCustomer stream demo
// @def primary                  ID
//
//go:generate toolkit gen model SourceCustomer --database DB
type SourceCustomer struct {
	ID        string `json:"id"        db:"id"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName"  db:"last_name"`
	Age       int    `json:"age"       db:"age"`
	City      string `json:"city"   db:"city"`
	datatypes.OperationTimesWithDeleted
}

package model

//go:generate easyjson -all source_customer.go
//easyjson:json
type SourceCustomer struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	City      string `json:"city"`
}

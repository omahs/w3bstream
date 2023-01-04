package model

//go:generate easyjson -all customer.go
//easyjson:json
type Customer struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	TaxNumber string `json:"taxNumber"`
	City      string `json:"city"`
}

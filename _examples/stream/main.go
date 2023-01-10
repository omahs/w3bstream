package main

import (
	"fmt"
	common "github.com/machinefi/w3bstream/_examples/wasm_common_go"
	"github.com/machinefi/w3bstream/stream/model"
	"github.com/mailru/easyjson"
	"github.com/tidwall/gjson"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//export filterAge
func _filterAge(rid uint32) int32 {
	common.Log(fmt.Sprintf("start rid: %d", rid))

	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log(err.Error())
		return -1
	}

	common.Log(fmt.Sprintf("get start resource all %v: %s", rid, string(message)))

	sourceCustomer := model.SourceCustomer{}
	easyjson.Unmarshal(message, &sourceCustomer)
	age := sourceCustomer.Age
	common.Log(fmt.Sprintf("get start resource age %v: %d", rid, age))

	if age >= 18 {
		common.Log(fmt.Sprintf("filter the Customer's age more than 18 %v: `%s`", rid, string(message)))
		common.SetDataByRID(rid, "true")
	} else if age < 18 {
		common.Log(fmt.Sprintf("filter the Customer's age less than 18 %v: `%s`", rid, string(message)))
		common.SetDataByRID(rid, "false")
	}

	return int32(rid)
}

//export mapTax
func _mapTax(rid uint32) int32 {
	common.Log(fmt.Sprintf("mapTax rid: %d", rid))

	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log(err.Error())
		return -1
	}

	common.Log(fmt.Sprintf("get mapTax resource all %v: %s", rid, string(message)))

	sourceCustomer := model.SourceCustomer{}
	easyjson.Unmarshal(message, &sourceCustomer)

	//TODO generate an error
	//common.Log(fmt.Sprintf("get mapTax sourceCustomer %d", sourceCustomer.Age))

	customer := model.Customer{}
	customer.ID = sourceCustomer.ID
	customer.FirstName = sourceCustomer.FirstName
	customer.LastName = sourceCustomer.LastName
	customer.Age = sourceCustomer.Age
	customer.City = sourceCustomer.City

	if customer.Age >= 30 {
		common.Log(fmt.Sprintf("the Customer's age is more than 30 %v: %s", rid, string(message)))
		customer.TaxNumber = "19832106687"
	}

	if b, err := easyjson.Marshal(customer); err != nil {
		common.Log(fmt.Sprintf("%v marshal error", sourceCustomer))
		return -1
	} else {
		common.SetDataByRID(rid, string(b))
	}

	return int32(rid)
}

//export groupByAge
func _groupByAge(rid uint32) int32 {
	common.Log(fmt.Sprintf("groupByAge rid: %d", rid))

	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log(err.Error())
		return -1
	}

	common.Log(fmt.Sprintf("get groupByAge resource all %v: `%s`", rid, string(message)))

	customer := string(message)

	city := gjson.Get(customer, "city").String()
	common.Log(fmt.Sprintf("get groupByAge resource city %v: %s", rid, city))

	common.SetDataByRID(rid, city)

	return int32(rid)
}

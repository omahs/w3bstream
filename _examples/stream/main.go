package main

import (
	"fmt"
	"github.com/google/uuid"
	common "github.com/machinefi/w3bstream/_examples/wasm_common_go"
	"github.com/machinefi/w3bstream/stream/model"
	"github.com/mailru/easyjson"
	"github.com/tidwall/gjson"
)

// main is required for TinyGo to compile to Wasm.
func main() {}

//export start
func _start(rid uint32) int32 {
	common.Log(fmt.Sprintf("start rid: %d", rid))

	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log(err.Error())
		return -1
	}

	defer func() {
		if common.FreeResource(rid) {
			common.Log(fmt.Sprintf("resource %v released", rid))
		}
	}()

	common.Log(fmt.Sprintf("get start resource all %v: %s", rid, string(message)))

	customer := model.Customer{}
	easyjson.Unmarshal(message, &customer)
	age := customer.Age
	//customer := string(message)
	//
	//age := gjson.Get(customer, "age").Int()
	common.Log(fmt.Sprintf("get start resource age %v: %d", rid, age))

	id := uuid.New().ID() % (^uint32(0) >> 1)
	if age >= 18 {
		common.Log(fmt.Sprintf("filter the Customer's age more than 18 %v: `%s`", rid, string(message)))
		common.SetDataByRID(id, "true")
	} else if age < 18 {
		common.Log(fmt.Sprintf("filter the Customer's age less than 18 %v: `%s`", rid, string(message)))
		common.SetDataByRID(id, "false")
	}

	return int32(id)
}

//export mapTax
func _mapTax(rid uint32) int32 {
	common.Log(fmt.Sprintf("mapTax rid: %d", rid))

	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log(err.Error())
		return -1
	}

	defer func() {
		if common.FreeResource(rid) {
			common.Log(fmt.Sprintf("resource %v released", rid))
		}
	}()

	common.Log(fmt.Sprintf("get mapTax resource all %v: %s", rid, string(message)))

	customer := model.Customer{}
	//e := customer.UnmarshalJSON(message)
	easyjson.Unmarshal(message, &customer)
	//common.Log(e.Error())

	//TODO generate an error
	//common.Log(fmt.Sprintf("get mapTax customer %d", customer.Age))

	id := uuid.New().ID() % (^uint32(0) >> 1)
	if customer.Age >= 30 {
		common.Log(fmt.Sprintf("the Customer's age is more than 30 %v: %s", rid, string(message)))
		customer.TaxNumber = "19832106687"
	}

	if b, err := easyjson.Marshal(customer); err != nil {
		common.Log(fmt.Sprintf("%v marshal error", customer))
		return -1
	} else {
		common.SetDataByRID(id, string(b))
	}

	return int32(id)
}

//export groupByAge
func _groupByAge(rid uint32) int32 {
	common.Log(fmt.Sprintf("groupByAge rid: %d", rid))

	message, err := common.GetDataByRID(rid)
	if err != nil {
		common.Log(err.Error())
		return -1
	}

	defer func() {
		if common.FreeResource(rid) {
			common.Log(fmt.Sprintf("resource %v released", rid))
		}
	}()

	common.Log(fmt.Sprintf("get groupByAge resource all %v: `%s`", rid, string(message)))

	customer := string(message)

	city := gjson.Get(customer, "city").String()
	common.Log(fmt.Sprintf("get groupByAge resource city %v: %s", rid, city))

	id := uuid.New().ID() % (^uint32(0) >> 1)
	common.SetDataByRID(id, city)

	return int32(id)
}

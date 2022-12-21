package main

import (
	"fmt"
	"github.com/google/uuid"
	common "github.com/machinefi/w3bstream/_examples/wasm_common_go"
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

	common.Log(fmt.Sprintf("get resource all %v: `%s`", rid, string(message)))

	customer := string(message)

	age := gjson.Get(customer, "age").Int()
	common.Log(fmt.Sprintf("get resource age %v: %d", rid, age))

	id := uuid.New().ID() % (^uint32(0) >> 1)
	if age >= 18 {
		common.Log(fmt.Sprintf("###### 大于等于 18 %v: `%s`", rid, string(message)))
		common.SetDataByRID(id, "true")
	} else if age < 18 {
		common.Log(fmt.Sprintf("&&&&&& 小于 18 %v: `%s`", rid, string(message)))
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

	common.Log(fmt.Sprintf("get resource all %v: `%s`", rid, string(message)))

	customer := string(message)

	age := gjson.Get(customer, "age").Int()
	common.Log(fmt.Sprintf("get resource age %v: %d", rid, age))

	id := uuid.New().ID() % (^uint32(0) >> 1)
	if age >= 30 {
		common.Log(fmt.Sprintf("###### 大于等于 30 %v: `%s`", rid, string(message)))
		common.SetDataByRID(id, "true")
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

	common.Log(fmt.Sprintf("get resource all %v: `%s`", rid, string(message)))

	customer := string(message)

	city := gjson.Get(customer, "city").String()
	common.Log(fmt.Sprintf("get resource age %v: %s", rid, city))

	id := uuid.New().ID() % (^uint32(0) >> 1)
	common.SetDataByRID(id, city)

	return int32(id)
}

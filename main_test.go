package main

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"testing"
)

func TestGjson(t *testing.T) {
	value := gjson.Get(" {\n    \"type\":\"hello\",\n    \"mac\":{\n\t\"hello\":\"world\"\n}\n}", "mac")
	log.Println(value)
}

type test struct {
	Str string `json:"str"`
}

func TestJson(t *testing.T) {
	m := make(map[string]interface{})
	x := make([]*test, 3)
	x[0] = &test{Str: "123"}
	m["type"] = x
	marshal, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(marshal))
}

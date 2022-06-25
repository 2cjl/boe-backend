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

func TestJson(t *testing.T) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte("{\n    \"type\":\"hello\",\n    \"mac\":{\n\t\"hello\":\"world\"\n}\n}"), &m)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(m["mac"].(map[string]interface{})["hello"])
}

package db

import (
	"boe-backend/internal/util/config"
	"flag"
	"log"
	"testing"
)

func TestGetOrganizationByUser(t *testing.T) {
	flag.Parse()
	config.InitViper()
	o := GetOrganizationByUser(1)
	log.Println(o)
}

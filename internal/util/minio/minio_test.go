package minio

import (
	"boe-backend/internal/util/config"
	"testing"
)

func TestMinio(t *testing.T) {
	config.InitViper()
	Init()
	PreSignObject()
}

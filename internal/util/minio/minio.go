package minio

import (
	"boe-backend/internal/util/config"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"time"
)

var (
	minioClient *minio.Client
)

func Init() {
	cfg := config.GetConfig().Minio
	var err error
	minioClient, err = minio.New(fmt.Sprintf("%s:%s", cfg.Endpoint, cfg.Port), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up
}

func PreSignObject() {
	expiry := time.Second * 24 * 60 * 60 // 1 day.
	presignedURL, err := minioClient.PresignedPutObject(context.Background(), "boe", "myobject", expiry)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)
}

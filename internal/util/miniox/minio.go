package miniox

import (
	"boe-backend/internal/util/config"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/url"
	"sync"
	"time"
)

const (
	expiry = time.Hour
)

var (
	minioClient *minio.Client
	once        sync.Once
	bucketName  string
)

func getClient() {
	if minioClient == nil {
		once.Do(func() {
			cfg := config.GetConfig().Minio
			var err error
			minioClient, err = minio.New(fmt.Sprintf("%s:%s", cfg.Endpoint, cfg.Port), &minio.Options{
				Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
				Secure: false,
			})
			if err != nil {
				panic(err)
			}
			bucketName = cfg.Bucket
		})
	}
}

func PreSignObject(objectPath string) *url.URL {
	getClient()
	preSignedURL, err := minioClient.PresignedPutObject(context.Background(), bucketName, objectPath, expiry)
	if err != nil {
		log.Println(err)
		return nil
	}
	return preSignedURL
}

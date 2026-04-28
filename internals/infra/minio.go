package infra

import (
	"log"
	"sync"

	"github.com/deltrexgg/ai-code-editor-server/internals/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	client *minio.Client
	once   sync.Once
)

func InitMinio(cfg config.MinioConfig) {
	once.Do(func() {
		var err error
		client, err = minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.UseSSL,
		})
		if err != nil {
			log.Fatalf("MinIO init failed: %v", err)
		}

		log.Println("MinIO connected:", cfg.Endpoint)
	})
}

func GetMinio() *minio.Client {
	if client == nil {
		log.Fatal("MinIO not initialized")
	}
	return client
}
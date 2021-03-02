package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"photo-gallery/callback"
	"photo-gallery/listener"
	"photo-gallery/storage"
	"photo-gallery/thumb"
)

const (
	minioEndpointEnv     = "MINIO_ENDPOINT"
	minioRootUserEnv     = "MINIO_ROOT_USER"
	minioRootPasswordEnv = "MINIO_ROOT_PASSWORD"
	photosBucketEnv      = "PHOTOS_BUCKET"
	thumbsBucketEnv      = "THUMBS_BUCKET"

	createdEvent = "s3:ObjectCreated:*"
)

func main() {
	_ = godotenv.Load()

	client, err := minio.New(os.Getenv(minioEndpointEnv), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv(minioRootUserEnv), os.Getenv(minioRootPasswordEnv), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	photosStorage := storage.NewStorage(client, os.Getenv(photosBucketEnv))
	thumbsStorage := storage.NewStorage(client, os.Getenv(thumbsBucketEnv))
	service := thumb.NewThumbsService(thumbsStorage)
	handler := callback.NewHandler(photosStorage, service)

	photosListener := listener.NewListener(os.Getenv(photosBucketEnv), client)
	photosListener.Listen(context.Background(), handler.Handle, createdEvent)
}

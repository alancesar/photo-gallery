package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"worker/callback"
	"worker/exif"
	"worker/listener"
	"worker/storage"
	"worker/thumb"
)

const (
	minioEndpointEnv     = "MINIO_ENDPOINT"
	minioRootUserEnv     = "MINIO_ROOT_USER"
	minioRootPasswordEnv = "MINIO_ROOT_PASSWORD"
	photosBucketEnv      = "PHOTOS_BUCKET"
	thumbsBucketEnv      = "THUMBS_BUCKET"
	exifBucketEnv        = "EXIF_BUCKET"

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
	thumbsService := thumb.NewService(thumbsStorage)
	exifService := exif.NewService(os.Getenv(photosBucketEnv), os.Getenv(exifBucketEnv), client)
	handler := callback.NewHandler(photosStorage, thumbsService, exifService)

	photosListener := listener.NewListener(os.Getenv(photosBucketEnv), client)
	photosListener.Listen(context.Background(), handler.Handle, createdEvent)
}

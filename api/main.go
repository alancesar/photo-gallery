package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"os"
	"photo-gallery/callback"
	"photo-gallery/database"
	"photo-gallery/handler"
	"photo-gallery/listener"
	"photo-gallery/metadata"
	"photo-gallery/storage"
)

const (
	dbHostEnv            = "DB_HOST"
	dbUserEnv            = "DB_USER"
	dbPasswordEnv        = "DB_PASSWORD"
	dbNameEnv            = "DB_NAME"
	dbPortEnv            = "DB_PORT"
	minioEndpointEnv     = "MINIO_ENDPOINT"
	minioRootUserEnv     = "MINIO_ROOT_USER"
	minioRootPasswordEnv = "MINIO_ROOT_PASSWORD"
	photosBucketEnv      = "PHOTOS_BUCKET"
	thumbsBucketEnv      = "THUMBS_BUCKET"
	metadataApiEnv       = "METADATA_API"

	createdEvent = "s3:ObjectCreated:*"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv(dbHostEnv),
		os.Getenv(dbUserEnv), os.Getenv(dbPasswordEnv), os.Getenv(dbNameEnv), os.Getenv(dbPortEnv))
	conn, err := database.NewConnection(dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db := database.NewDatabase(conn)

	client, err := minio.New(os.Getenv(minioEndpointEnv), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv(minioRootUserEnv), os.Getenv(minioRootPasswordEnv), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	photosStorage := storage.NewStorage(client, os.Getenv(photosBucketEnv))
	thumbsStorage := storage.NewStorage(client, os.Getenv(thumbsBucketEnv))
	metadataService := metadata.NewService(os.Getenv(metadataApiEnv), http.DefaultClient)
	callbackHandler := callback.NewCallbackHandler(db, photosStorage, thumbsStorage, metadataService)
	photosListener := listener.NewListener(os.Getenv(photosBucketEnv), client)

	go photosListener.Listen(context.Background(), callbackHandler.Handle, createdEvent)

	engine := gin.Default()
	engine.Use(cors.Default())
	engine.Handle(http.MethodGet, "/api/photos", handler.PhotoHandler(db))
	if err := engine.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}

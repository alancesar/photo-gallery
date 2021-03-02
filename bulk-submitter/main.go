package main

import (
	"bulk-submitter/handler"
	"bulk-submitter/storage"
	"flag"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

const (
	minioEndpointEnv  = "MINIO_ENDPOINT"
	minioAccessKeyEnv = "MINIO_ACCESS_KEY"
	minioSecretKeyEnv = "MINIO_SECRET_KEY"
	bucketNameEnv     = "PHOTOS_BUCKET"
)

var rootPath *string

func init() {
	rootPath = flag.String("root-path", "", "root path for bulk submitter")
	flag.Parse()

	if *rootPath == "" {
		log.Fatalln("root path must be defined for bulk submitter")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	client, err := minio.New(os.Getenv(minioEndpointEnv), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv(minioAccessKeyEnv), os.Getenv(minioSecretKeyEnv), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	s := storage.NewStorage(client, os.Getenv(bucketNameEnv))
	bs := handler.NewBulkSubmitter(s)
	bs.Submit(*rootPath)
}

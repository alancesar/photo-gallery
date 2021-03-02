package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"os"
)

type Storage struct {
	client *minio.Client
	bucket string
}

func NewStorage(client *minio.Client, bucket string) *Storage {
	return &Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *Storage) Put(ctx context.Context, filename, contentType string, file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(ctx, s.bucket, filename, file, stat.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

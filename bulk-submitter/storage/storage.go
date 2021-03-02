package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
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

func (s *Storage) Put(ctx context.Context, filename, filepath string) error {
	_, err := s.client.FPutObject(ctx, s.bucket, filename, filepath, minio.PutObjectOptions{})

	return err
}

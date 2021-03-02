package thumb

import (
	"context"
	"io"
	"os"
	"photo-gallery/photo"
	"photo-gallery/storage"
)

const (
	thumbWidth  = 800
	thumbHeight = 600
)

type Service struct {
	storage *storage.Storage
}

func NewThumbsService(s *storage.Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) CreateThumbs(ctx context.Context, input io.Reader, filename, contentType string) error {
	resized, err := photo.Fit(input, thumbWidth, thumbHeight)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(resized.Name())
	}()

	return s.storage.Put(ctx, filename, contentType, resized)
}

package thumb

import (
	"context"
	"github.com/disintegration/imaging"
	"io"
	"io/ioutil"
	"os"
	"worker/storage"
)

const (
	thumbWidth  = 800
	thumbHeight = 600
)

type Service struct {
	storage *storage.Storage
}

func NewService(s *storage.Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) CreateThumbs(ctx context.Context, input io.Reader, filename, contentType string) error {
	resized, err := fit(input, thumbWidth, thumbHeight)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(resized.Name())
	}()

	return s.storage.Put(ctx, filename, contentType, resized)
}

func fit(reader io.Reader, width, height int) (*os.File, error) {
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	resized := imaging.Fit(img, width, height, imaging.Lanczos)
	output, err := ioutil.TempFile("", "thumb.*.jpg")
	if err != nil {
		return nil, err
	}

	err = imaging.Save(resized, output.Name())
	return output, err
}

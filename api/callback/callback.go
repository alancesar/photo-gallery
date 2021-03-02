package callback

import (
	"context"
	"github.com/google/uuid"
	"log"
	"photo-gallery/database"
	"photo-gallery/metadata"
	"photo-gallery/photo"
	"photo-gallery/storage"
	"photo-gallery/thumb"
)

type PhotoHandler struct {
	photosStorage *storage.Storage
	service       *thumb.Service
}

func NewPhotoCallbackHandler(photosStorage *storage.Storage, service *thumb.Service) *PhotoHandler {
	return &PhotoHandler{
		photosStorage: photosStorage,
		service:       service,
	}
}

func (h *PhotoHandler) Handle(ctx context.Context, key, contentType string) {
	item, err := h.photosStorage.Get(ctx, key)
	if err != nil {
		log.Println(err)
		return
	}

	if err = h.service.CreateThumbs(ctx, item, key, contentType); err != nil {
		log.Println(err)
		return
	}
}

type ThumbHandler struct {
	db              *database.Database
	metadataService *metadata.Service
}

func NewThumbHandler(db *database.Database, metadataService *metadata.Service) *ThumbHandler {
	return &ThumbHandler{
		db:              db,
		metadataService: metadataService,
	}
}

func (h ThumbHandler) Handle(_ context.Context, key, _ string) {
	exif, err := h.metadataService.GetExif(key)
	if err != nil {
		log.Println(err)
		return
	}

	width, height, _ := exif.GetSizes()

	if err = h.db.Save(photo.Photo{
		ID:       uuid.New(),
		Exif:     exif,
		Filename: key,
		Width:    width,
		Height:   height,
	}); err != nil {
		log.Println(err)
		return
	}
}

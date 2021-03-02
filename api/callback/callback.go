package callback

import (
	"context"
	"github.com/google/uuid"
	"log"
	"os"
	"photo-gallery/database"
	"photo-gallery/metadata"
	"photo-gallery/photo"
	"photo-gallery/storage"
)

const (
	thumbWidth  = 800
	thumbHeight = 600
)

type Handler struct {
	db              *database.Database
	photosStorage   *storage.Storage
	thumbsStorage   *storage.Storage
	metadataService *metadata.Service
}

func NewCallbackHandler(
	db *database.Database,
	photosStorage *storage.Storage,
	thumbsStorage *storage.Storage,
	metadataService *metadata.Service,
) *Handler {
	return &Handler{
		db:              db,
		photosStorage:   photosStorage,
		thumbsStorage:   thumbsStorage,
		metadataService: metadataService,
	}
}

func (h *Handler) Handle(ctx context.Context, key, contentType string) {
	item, err := h.photosStorage.Get(ctx, key)
	if err != nil {
		log.Println(err)
		return
	}

	resized, err := photo.Fit(item, thumbWidth, thumbHeight)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = os.Remove(resized.Name())
	}()

	if err = h.thumbsStorage.Put(ctx, key, contentType, resized); err != nil {
		log.Println(err)
		return
	}

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

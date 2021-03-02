package callback

import (
	"context"
	"github.com/google/uuid"
	"log"
	"photo-gallery/database"
	"photo-gallery/metadata"
	"photo-gallery/photo"
)

type Handler struct {
	db              *database.Database
	metadataService *metadata.Service
}

func NewHandler(db *database.Database, metadataService *metadata.Service) *Handler {
	return &Handler{
		db:              db,
		metadataService: metadataService,
	}
}

func (h Handler) Handle(_ context.Context, key, _ string) {
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

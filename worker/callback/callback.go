package callback

import (
	"context"
	"log"
	"photo-gallery/storage"
	"photo-gallery/thumb"
)

type Handler struct {
	photosStorage *storage.Storage
	service       *thumb.Service
}

func NewHandler(photosStorage *storage.Storage, service *thumb.Service) *Handler {
	return &Handler{
		photosStorage: photosStorage,
		service:       service,
	}
}

func (h *Handler) Handle(ctx context.Context, key, contentType string) {
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

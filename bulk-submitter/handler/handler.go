package handler

import (
	"bulk-submitter/storage"
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
)

const (
	expectedMIME     = "image/jpeg"
	defaultExtension = ".jpg"
)

type BulkSubmitter struct {
	storage *storage.Storage
}

func NewBulkSubmitter(storage *storage.Storage) BulkSubmitter {
	return BulkSubmitter{
		storage: storage,
	}
}

func (bs *BulkSubmitter) Submit(rootPath string) {
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		mime, err := mimetype.DetectFile(path)
		if err != nil {
			log.Println(fmt.Sprintf("error on open %s", path))
			return err
		}

		if !mime.Is(expectedMIME) {
			log.Println(fmt.Sprintf("ignored %s", path))
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Println(fmt.Sprintf("error on open %s", path))
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		filename := fmt.Sprintf("%s%s", uuid.New().String(), defaultExtension)
		if err = bs.storage.Put(context.Background(), filename, mime.String(), file); err != nil {
			log.Println(fmt.Sprintf("error while uploading %s", path))
			return err
		}

		log.Println(fmt.Sprintf("upload %s successfuly", path))
		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}
}

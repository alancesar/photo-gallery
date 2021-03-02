package handler

import (
	"bulk-submitter/storage"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
)

const (
	expectedMimetype = "image/jpeg"
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
		if info.IsDir() || mime.TypeByExtension(filepath.Ext(path)) != expectedMimetype {
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

		hash := sha256.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			log.Println(fmt.Sprintf("error on calculate file hash %s", path))
			return err
		}

		filename := fmt.Sprintf("%s%s", hex.EncodeToString(hash.Sum(nil)), defaultExtension)
		if err = bs.storage.Put(context.Background(), filename, path); err != nil {
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

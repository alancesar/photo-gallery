package database

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"photo-gallery/photo"
	"time"
)

const (
	photosTableName = "photos"
)

type entity struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Exif      datatypes.JSON
	Filename  string
	Width     int
	Height    int
}

func (*entity) TableName() string {
	return photosTableName
}

type Database struct {
	db *gorm.DB
}

func NewConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func NewDatabase(db *gorm.DB) *Database {
	_ = db.AutoMigrate(&entity{})
	return &Database{
		db: db,
	}
}

func (d *Database) Get(id string) (photo.Photo, error) {
	e := entity{
		ID: uuid.MustParse(id),
	}

	if query := d.db.Take(&e); query.Error != nil {
		return photo.Photo{}, query.Error
	}

	ex := photo.Exif{}
	if err := json.Unmarshal(e.Exif, &ex); err != nil {
		return photo.Photo{}, err
	}

	return photo.Photo{
		ID:       e.ID,
		Exif:     ex,
		Filename: e.Filename,
		Width:    e.Width,
		Height:   e.Width,
	}, nil
}

func (d *Database) GetAll() ([]photo.Photo, error) {
	var entities []entity

	if query := d.db.Find(&entities); query.Error != nil {
		return nil, query.Error
	}

	photos := make([]photo.Photo, len(entities))
	for index := range entities {
		photos[index] = photo.Photo{
			ID:       entities[index].ID,
			Filename: entities[index].Filename,
			Width:    entities[index].Width,
			Height:   entities[index].Height,
		}
	}

	return photos, nil
}

func (d Database) Save(p photo.Photo) error {
	ex, err := json.Marshal(&p.Exif)
	if err != nil {
		return err
	}

	e := entity{
		ID:       p.ID,
		Exif:     ex,
		Filename: p.Filename,
		Width:    p.Width,
		Height:   p.Height,
	}

	return d.db.Create(&e).Error
}

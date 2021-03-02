package photo

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
	"strconv"
)

var (
	sizeRegex            = regexp.MustCompile(`(\d{1,8})`)
	errSizeTagNotPresent = errors.New("size tag not present")
)

const (
	imageWidthKey  = "Image Width"
	imageHeightKey = "Image Height"
)

type Exif map[string]string

type Photo struct {
	ID       uuid.UUID `json:"id"`
	Exif     Exif      `json:"exif,omitempty"`
	Filename string    `json:"filename"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
}

func (e Exif) GetSizes() (width int, height int, err error) {
	if stringWidth, ok := e[imageWidthKey]; !ok {
		return 0, 0, errSizeTagNotPresent
	} else {
		if width, err = strconv.Atoi(sizeRegex.FindString(stringWidth)); err != nil {
			return 0, 0, err
		}
	}

	if stringHeight, ok := e[imageHeightKey]; !ok {
		return 0, 0, errSizeTagNotPresent
	} else {
		if height, err = strconv.Atoi(sizeRegex.FindString(stringHeight)); err != nil {
			return 0, 0, err
		}
	}

	return width, height, nil
}

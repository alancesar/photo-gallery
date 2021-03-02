package photo

import (
	"github.com/disintegration/imaging"
	"io"
	"io/ioutil"
	"os"
)

func Fit(reader io.Reader, width, height int) (*os.File, error) {
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

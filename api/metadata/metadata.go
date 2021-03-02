package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"photo-gallery/photo"
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Service struct {
	host   string
	client HttpClient
}

func NewService(host string, client HttpClient) *Service {
	return &Service{
		host:   host,
		client: client,
	}
}

func (s *Service) GetExif(filename string) (photo.Exif, error) {
	url := fmt.Sprintf("%s/api/exif/%s", s.host, filename)
	response, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status from %s: %s", response.Request.URL.String(), response.Status)
	}

	var exif photo.Exif
	err = json.NewDecoder(response.Body).Decode(&exif)
	return exif, err
}

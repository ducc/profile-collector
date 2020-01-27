package scraper

import (
	"context"
	"github.com/ducc/profile-collector/protos"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
)

type Scraper interface {
	Scrape(ctx context.Context, address string) (*protos.Profile, error)
}

type httpScraper struct {
	client *http.Client
}

func NewHTTP() (Scraper, error) {
	return &httpScraper{
		client: http.DefaultClient,
	}, nil
}

func (s *httpScraper) Scrape(ctx context.Context, address string) (*protos.Profile, error) {
	req, err := s.buildRequest(address)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	// todo log res status

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	profile := &protos.Profile{}
	if err := proto.Unmarshal(data, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *httpScraper) buildRequest(address string) (*http.Request, error) {
	return http.NewRequest("GET", address, nil)
}

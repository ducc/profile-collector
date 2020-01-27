package service

import "time"

type Service struct {
	Name           string
	Address        string
	ScrapeInterval time.Duration
	LastScrape     time.Time
}

func NewService(name, address string, scrapeInterval time.Duration) *Service {
	return &Service{
		Name:           name,
		Address:        address,
		ScrapeInterval: scrapeInterval,
	}
}

func (s *Service) ScrapeNow() bool {
	return time.Now().After(s.LastScrape.Add(s.ScrapeInterval))
}

package scraper

import "github.com/antoniopantaleo/wwdc/internal/domain"

type StubScraper struct {
	Events []domain.WWDCEvent
	Err    error
}

func (s *StubScraper) Scrape() ([]domain.WWDCEvent, error) {
	return s.Events, s.Err
}

package usecases

import "github.com/antoniopantaleo/wwdc/internal/domain"

type ScrapeAndExportUseCase struct {
	scraper  domain.WWDCScraper
	exporter domain.WWDCExporter
}

func NewScrapeAndExportUseCase(scraper domain.WWDCScraper, exporter domain.WWDCExporter) *ScrapeAndExportUseCase {
	return &ScrapeAndExportUseCase{
		scraper:  scraper,
		exporter: exporter,
	}
}

func (s *ScrapeAndExportUseCase) Execute() error {
	events, err := s.scraper.Scrape()
	if err != nil {
		return err
	}
	return s.exporter.Export(events)
}

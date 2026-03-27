package usecases

import (
	"errors"
	"testing"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

type mockScraper struct {
	scrapeFunc func() ([]domain.WWDCEvent, error)
}

func (m *mockScraper) Scrape() ([]domain.WWDCEvent, error) {
	return m.scrapeFunc()
}

type mockExporter struct {
	exportFunc func(events []domain.WWDCEvent) error
}

func (m *mockExporter) Export(events []domain.WWDCEvent) error {
	return m.exportFunc(events)
}

func TestScrapeHappyPath(t *testing.T) {
	mockEvents := []domain.WWDCEvent{
		{
			Title:    "WWDC26",
			Year:     2026,
			CoverURL: "https://example.com/wwdc26.jpg",
			Videos: []domain.WWDCVideo{
				{
					Title:    "Session 1",
					VideoURL: "https://example.com/session1.mp4",
					Content:  "This is the content of session 1.",
				},
			},
		},
	}
	scraper := &mockScraper{
		scrapeFunc: func() ([]domain.WWDCEvent, error) {
			return mockEvents, nil
		},
	}
	var exportedEvents []domain.WWDCEvent
	exporter := &mockExporter{
		exportFunc: func(events []domain.WWDCEvent) error {
			exportedEvents = events
			return nil
		},
	}
	sut := NewScrapeAndExportUseCase(scraper, exporter)
	err := sut.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !areEventsEqual(exportedEvents, mockEvents) {
		t.Fatalf("expected exported events to be %v, got %v", mockEvents, exportedEvents)
	}
}

func TestScrapeScraperError(t *testing.T) {
	errScrape := errors.New("some scraping error")
	scraper := &mockScraper{
		scrapeFunc: func() ([]domain.WWDCEvent, error) {
			return nil, errScrape
		},
	}
	exportFunctionCalled := false
	exporter := &mockExporter{
		exportFunc: func(events []domain.WWDCEvent) error {
			exportFunctionCalled = true
			return nil
		},
	}

	sut := NewScrapeAndExportUseCase(scraper, exporter)
	err := sut.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if exportFunctionCalled {
		t.Fatal("expected export function not to be called, but it was")
	}
}

func TestScrapeExportError(t *testing.T) {
	mockEvents := []domain.WWDCEvent{
		{
			Title:    "WWDC26",
			Year:     2026,
			CoverURL: "https://example.com/wwdc26.jpg",
			Videos: []domain.WWDCVideo{
				{
					Title:    "Session 1",
					VideoURL: "https://example.com/session1.mp4",
					Content:  "This is the content of session 1.",
				},
			},
		},
	}
	scraper := &mockScraper{
		scrapeFunc: func() ([]domain.WWDCEvent, error) {
			return mockEvents, nil
		},
	}
	errExport := errors.New("some export error")
	exporter := &mockExporter{
		exportFunc: func(events []domain.WWDCEvent) error {
			return errExport
		},
	}

	sut := NewScrapeAndExportUseCase(scraper, exporter)
	err := sut.Execute()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func areEventsEqual(a, b []domain.WWDCEvent) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Title != b[i].Title ||
			a[i].Year != b[i].Year ||
			a[i].CoverURL != b[i].CoverURL ||
			!areVideosEqual(a[i].Videos, b[i].Videos) {
			return false
		}
	}
	return true
}

func areVideosEqual(a, b []domain.WWDCVideo) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Title != b[i].Title ||
			a[i].VideoURL != b[i].VideoURL ||
			a[i].Content != b[i].Content {
			return false
		}
	}
	return true
}

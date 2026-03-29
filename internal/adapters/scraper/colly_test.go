package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollyScraperEventsList(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()

	scraper := NewCollyScraper(ts.URL)
	events, err := scraper.Scrape()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(events) != 6 {
		t.Fatalf("expected 6 events, got %d", len(events))
	}
}

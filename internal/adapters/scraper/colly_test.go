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
	for _, event := range events {
		for _, video := range event.Videos {
			if video.Title != "Bring your Live Activity to Apple Watch" {
				t.Fatalf("expected Bring your Live Activity to Apple Watch title, got %s instead", video.Title)
			}
			if video.VideoURL != "https://devstreaming-cdn.apple.com/videos/wwdc/2024/10068/4/C621DA91-3F64-481C-8D10-25A5C5FCD587/downloads/wwdc2024-10068_hd.mp4?dl=1" {
			t.Fatalf("expected https://devstreaming-cdn.apple.com/videos/wwdc/2024/10068/4/C621DA91-3F64-481C-8D10-25A5C5FCD587/downloads/wwdc2024-10068_hd.mp4?dl=1 videoURL, got %s instead", video.VideoURL)
		}
		}
	}
}

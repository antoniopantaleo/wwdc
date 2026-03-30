package scraper

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

func TestCollyScraperEventsList(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()

	logger := log.New(os.Stdout, "[COLLY TEST] ", log.LstdFlags)
	scraper := NewCollyScraper(ts.URL, logger)
	events, err := scraper.Scrape()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(events) != 6 {
		t.Fatalf("expected 6 events, got %d", len(events))
	}
	wwdc2024Idx := slices.IndexFunc(events, func(e domain.WWDCEvent) bool {
		return e.Year == 2024
	})
	if wwdc2024Idx == -1 || len(events) <= wwdc2024Idx {
		t.Fatal("unable to find wwdc2024 event")
	}
	wwdc2024 := events[wwdc2024Idx]
	appleWatchVideoIdx := slices.IndexFunc(wwdc2024.Videos, func (v domain.WWDCVideo) bool {
		return v.Title == "Bring your Live Activity to Apple Watch"
	})
	if appleWatchVideoIdx == -1 || len(wwdc2024.Videos) <= appleWatchVideoIdx {
		t.Fatal("unable to find Apple Watch video")
	}
	appleWatchVideo := wwdc2024.Videos[appleWatchVideoIdx]
	if appleWatchVideo.VideoURL != "https://devstreaming-cdn.apple.com/videos/wwdc/2024/10068/4/C621DA91-3F64-481C-8D10-25A5C5FCD587/downloads/wwdc2024-10068_hd.mp4?dl=1" {
			t.Fatalf("expected https://devstreaming-cdn.apple.com/videos/wwdc/2024/10068/4/C621DA91-3F64-481C-8D10-25A5C5FCD587/downloads/wwdc2024-10068_hd.mp4?dl=1 videoURL, got %s instead", appleWatchVideo.VideoURL)
		}
	
}

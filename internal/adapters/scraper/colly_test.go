package scraper

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

type mockProgressReporter struct {
	InfoFunc    func(message string)
	WarningFunc func(message string)
}

func (m *mockProgressReporter) Info(message string) {
	if m.InfoFunc != nil {
		m.InfoFunc(message)
	}
}

func (m *mockProgressReporter) Warning(message string) {
	if m.WarningFunc != nil {
		m.WarningFunc(message)
	}
}

func noOpReporter() *mockProgressReporter {
	return &mockProgressReporter{}
}

func TestCollyScraperEventsList(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()

	scraper := NewCollyScraper(ts.URL, noOpReporter())
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
	appleWatchVideoIdx := slices.IndexFunc(wwdc2024.Videos, func(v domain.WWDCVideo) bool {
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

func TestCollyScraperReturnsErrorWhenEventsPageFails(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/videos", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	mux.HandleFunc("/videos/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	scraper := NewCollyScraper(ts.URL, noOpReporter())
	_, err := scraper.Scrape()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestCollyScraperLogsWarningOnFailedVideoPage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/videos/play/wwdc2024/10068/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	mux.Handle("/", http.FileServer(http.Dir("testdata")))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	var warnings []string
	reporter := &mockProgressReporter{
		WarningFunc: func(message string) {
			warnings = append(warnings, message)
		},
	}
	scraper := NewCollyScraper(ts.URL, reporter)
	events, err := scraper.Scrape()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) == 0 {
		t.Fatal("expected at least some events")
	}

	found := slices.ContainsFunc(warnings, func(w string) bool {
		return strings.Contains(w, "10068")
	})
	if !found {
		t.Fatal("expected a warning about video 10068, got none")
	}
}

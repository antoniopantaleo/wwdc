package exporter

import (
	"bytes"
	"testing"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

func TestJSONExportHappyPath(t *testing.T) {
	var buf bytes.Buffer
	sut := NewJSONExporter(&buf)
	events := []domain.WWDCEvent{
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
	err := sut.Export(events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expectedJSON := `{
	"events": [
		{
			"title": "WWDC26",
			"year": 2026,
			"coverUrl": "https://example.com/wwdc26.jpg",
			"videos": [
				{
					"title": "Session 1",
					"videoUrl": "https://example.com/session1.mp4",
					"content": "This is the content of session 1."
				}
			]
		}
	]
}
`
	if buf.String() != expectedJSON {
		t.Fatalf("expected JSON to be %s, got %s", expectedJSON, buf.String())
	}
}

package exporter

import (
	"os"
	"testing"

	"github.com/antoniopantaleo/wwdc/internal/adapters/filesystem"
	"github.com/antoniopantaleo/wwdc/internal/domain"
)

func TestMarkdownExporterOSFileSystemWritesFoldersAndFiles(t *testing.T) {
	base := t.TempDir()
	fs := filesystem.NewOSFileSystem(base)
	sut := NewMarkdownExporter(fs)
	events := []domain.WWDCEvent{
		{
			Title:    "WWDC24",
			Year:     2024,
			CoverURL: "https://example.com/wwdc24.jpg",
			Videos: []domain.WWDCVideo{
				{
					Title:    "Session 1",
					VideoURL: "https://example.com/session1.mp4",
					Content:  "This is the content of session 1.",
				},
				{
					Title:    "Session 2",
					VideoURL: "https://example.com/session2.mp4",
					Content:  "This is the content of session 2.",
				},
			},
		},

		{
			Title:    "WWDC23",
			Year:     2023,
			CoverURL: "https://example.com/wwdc23.jpg",
			Videos: []domain.WWDCVideo{
				{
					Title:    "Session 3",
					VideoURL: "https://example.com/session3.mp4",
					Content:  "This is the content of session 3.",
				},
				{
					Title:    "Session 4",
					VideoURL: "https://example.com/session4.mp4",
					Content:  "This is the content of session 4.",
				},
			},
		},
	}
	err := sut.Export(events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expectedDirectories := []string{
		base + "/WWDC24",
		base + "/WWDC23",
	}
	expectedFilePaths := []string{
		base + "/WWDC24/Session 1.md",
		base + "/WWDC24/Session 2.md",
		base + "/WWDC23/Session 3.md",
		base + "/WWDC23/Session 4.md",
	}
	for _, dir := range expectedDirectories {
		if !dirExists(dir) {
			t.Fatalf("expected directory %s to exist, but it does not", dir)
		}
	}
	for _, filePath := range expectedFilePaths {
		if !fileExists(filePath) {
			t.Fatalf("expected file %s to exist, but it does not", filePath)
		}
	}
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

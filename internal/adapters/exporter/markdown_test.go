package exporter

import (
	"testing"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

type mockFileSystem struct {
	makeDirFunc   func(path string) error
	writeFileFunc func(path string, data []byte) error
}

func (m *mockFileSystem) MakeDir(path string) error {
	return m.makeDirFunc(path)
}

func (m *mockFileSystem) WriteFile(path string, data []byte) error {
	return m.writeFileFunc(path, data)
}

func TestMarkdownExportHappyPath(t *testing.T) {
	var (
		writtenDirectory string
		writtenFilePath  string
		writtenData      []byte
	)
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			writtenDirectory = path
			return nil
		},
		writeFileFunc: func(path string, data []byte) error {
			writtenFilePath = path
			writtenData = data
			return nil
		},
	}
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
	sut := NewMarkdownExporter(&fs)
	err := sut.Export(events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if writtenDirectory != "WWDC26" {
		t.Fatalf("expected directory to be %s, got %s", "WWDC26", writtenDirectory)
	}
	if writtenFilePath != "WWDC26/Session 1.md" {
		t.Fatalf("expected file path to be %s, got %s", "WWDC26/Session 1.md", writtenFilePath)
	}
	expectedData := []byte("# Session 1\n\n[Video](https://example.com/session1.mp4)\n\nThis is the content of session 1.")

	if string(writtenData) != string(expectedData) {
		t.Fatalf("expected data to be %s, got %s", expectedData, writtenData)
	}
}

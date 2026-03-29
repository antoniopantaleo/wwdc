package exporter

import (
	"errors"
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

func TestMarkdownExportHappyPathOneEventOneVideo(t *testing.T) {
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
	expectedData := []byte("# Session 1\n\n![Video](https://example.com/session1.mp4)\n\nThis is the content of session 1.")

	if string(writtenData) != string(expectedData) {
		t.Fatalf("expected data to be %s, got %s", expectedData, writtenData)
	}
}

func TestMarkdownExportHappyPathOneEventMultipleVideos(t *testing.T) {
	var (
		writtenDirectory string
		writtenFilePaths  []string
		writtenData      [][]byte
	)
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			writtenDirectory = path
			return nil
		},
		writeFileFunc: func(path string, data []byte) error {
			writtenFilePaths = append(writtenFilePaths, path)
			writtenData = append(writtenData, data)
			return nil
		},
	}
	events := []domain.WWDCEvent{
		{
			Title:    "WWDC25",
			Year:     2025,
			CoverURL: "https://example.com/wwdc25.jpg",
			Videos: []domain.WWDCVideo{
				{
					Title:    "Session 1",
					VideoURL: "https://example.com/session1.mp4",
					Content:  "This is the content of session 1.",
				},
				{
					Title: "Session 2",
					VideoURL: "https://example.com/session2.mp4",
					Content:  "This is the content of session 2.",
				},
			},
		},
	}
	sut := NewMarkdownExporter(&fs)
	err := sut.Export(events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if writtenDirectory != "WWDC25" {
		t.Fatalf("expected directory to be %s, got %s", "WWDC25", writtenDirectory)
	}
	if writtenFilePaths[0] != "WWDC25/Session 1.md" {
		t.Fatalf("expected file path to be %s, got %s", "WWDC25/Session 1.md", writtenFilePaths[0])
	}
	expectedData := []byte("# Session 1\n\n![Video](https://example.com/session1.mp4)\n\nThis is the content of session 1.")

	if string(writtenData[0]) != string(expectedData) {
		t.Fatalf("expected data to be %s, got %s", expectedData, writtenData[0])
	}

	if writtenFilePaths[1] != "WWDC25/Session 2.md" {
		t.Fatalf("expected file path to be %s, got %s", "WWDC25/Session 2.md", writtenFilePaths[1])
	}
	expectedData2 := []byte("# Session 2\n\n![Video](https://example.com/session2.mp4)\n\nThis is the content of session 2.")

	if string(writtenData[1]) != string(expectedData2) {
		t.Fatalf("expected data to be %s, got %s", expectedData2, writtenData[1])
	}
}

func TestMarkdownExportPathHappyPathMultipleEvents(t *testing.T) {
	var (
		writtenDirectories []string
		writtenFilePaths  []string
		writtenData      [][]byte
	)
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			writtenDirectories = append(writtenDirectories, path)
			return nil
		},
		writeFileFunc: func(path string, data []byte) error {
			writtenFilePaths = append(writtenFilePaths, path)
			writtenData = append(writtenData, data)
			return nil
		},
	}
	events := []domain.WWDCEvent{
		{
			Title: "WWDC24",
			Year:  2024,
			CoverURL: "https://example.com/wwdc24.jpg",
			Videos: []domain.WWDCVideo{
				{
					Title:    "Session 1",
					VideoURL: "https://example.com/session1.mp4",
					Content:  "This is the content of session 1.",
				},	
				{
					Title: "Session 2",
					VideoURL: "https://example.com/session2.mp4",
					Content:  "This is the content of session 2.",
				},
			},
		},

		{
			Title: "WWDC23",
			Year:  2023,
			CoverURL: "https://example.com/wwdc23.jpg",
			Videos: []domain.WWDCVideo{
				{
					Title:    "Session 3",
					VideoURL: "https://example.com/session3.mp4",
					Content:  "This is the content of session 3.",
				},
				{
					Title: "Session 4",
					VideoURL: "https://example.com/session4.mp4",
					Content:  "This is the content of session 4.",
				},
			},
		},
	}
	sut := NewMarkdownExporter(&fs)
	err := sut.Export(events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if writtenDirectories[0] != "WWDC24" {
		t.Fatalf("expected directory to be %s, got %s", "WWDC24", writtenDirectories[0])
	}
	if writtenDirectories[1] != "WWDC23" {
		t.Fatalf("expected directory to be %s, got %s", "WWDC23", writtenDirectories[1])
	}
	if writtenFilePaths[0] != "WWDC24/Session 1.md" {
		t.Fatalf("expected file path to be %s, got %s", "WWDC24/Session 1.md", writtenFilePaths[0])
	}
	expectedData := []byte("# Session 1\n\n![Video](https://example.com/session1.mp4)\n\nThis is the content of session 1.")

	if string(writtenData[0]) != string(expectedData) {
		t.Fatalf("expected data to be %s, got %s", expectedData, writtenData[0])
	}

	if writtenFilePaths[1] != "WWDC24/Session 2.md" {
		t.Fatalf("expected file path to be %s, got %s", "WWDC24/Session 2.md", writtenFilePaths[1])
	}
	expectedData2 := []byte("# Session 2\n\n![Video](https://example.com/session2.mp4)\n\nThis is the content of session 2.")

	if string(writtenData[1]) != string(expectedData2) {
		t.Fatalf("expected data to be %s, got %s", expectedData2, writtenData[1])
	}

	if writtenFilePaths[2] != "WWDC23/Session 3.md" {
		t.Fatalf("expected file path to be %s, got %s", "WWDC23/Session 3.md", writtenFilePaths[2])
	}
	expectedData3 := []byte("# Session 3\n\n![Video](https://example.com/session3.mp4)\n\nThis is the content of session 3.")

	if string(writtenData[2]) != string(expectedData3) {
		t.Fatalf("expected data to be %s, got %s", expectedData3, writtenData[2])
	}

	if writtenFilePaths[3] != "WWDC23/Session 4.md" {
		t.Fatalf("expected file path to be %s, got %s", "WWDC23/Session 4.md", writtenFilePaths[3])
	}
	expectedData4 := []byte("# Session 4\n\n![Video](https://example.com/session4.mp4)\n\nThis is the content of session 4.")

	if string(writtenData[3]) != string(expectedData4) {
		t.Fatalf("expected data to be %s, got %s", expectedData4, writtenData[3])
	}
}

func TestMarkdownExportOneEventNoVideos(t *testing.T) {
	var (
		writtenDirectory string
		writtenFilePaths  []string
		writtenData      [][]byte
	)
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			writtenDirectory = path
			return nil
		},
		writeFileFunc: func(path string, data []byte) error {
			writtenFilePaths = append(writtenFilePaths, path)
			writtenData = append(writtenData, data)
			return nil
		},
	}
	events := []domain.WWDCEvent{
		{
			Title: "WWDC24",
			Year:  2024,
			CoverURL: "https://example.com/wwdc24.jpg",
			Videos: []domain.WWDCVideo{},
		},
	}
	sut := NewMarkdownExporter(&fs)
	err := sut.Export(events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if writtenDirectory != "WWDC24" {
		t.Fatalf("expected directory to be %s, got %s", "WWDC24", writtenDirectory)
	}
	if len(writtenFilePaths) != 0 {
		t.Fatalf("expected no file to be written, but got %d files", len(writtenFilePaths))
	}
	if len(writtenData) != 0 {
		t.Fatalf("expected no data to be written, but got %d data entries", len(writtenData))
	}
}

func TestMarkdownExportNoEvents(t *testing.T) {
	var (
		writtenDirectory string
		writtenFilePaths  []string
		writtenData      [][]byte
	)
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			writtenDirectory = path
			return nil
		},
		writeFileFunc: func(path string, data []byte) error {
			writtenFilePaths = append(writtenFilePaths, path)
			writtenData = append(writtenData, data)
			return nil
		},
	}
	events := []domain.WWDCEvent{}
	sut := NewMarkdownExporter(&fs)
	err := sut.Export(events)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if writtenDirectory != "" {
		t.Fatalf("expected no directory to be created, got %s", writtenDirectory)
	}
	if len(writtenFilePaths) != 0 {
		t.Fatalf("expected no file to be written, but got %d files", len(writtenFilePaths))
	}
	if len(writtenData) != 0 {
		t.Fatalf("expected no data to be written, but got %d data entries", len(writtenData))
	}
}

func TestMarkdownExportInvalidYear(t *testing.T) {
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			return nil
		},
		writeFileFunc: func(path string, data []byte) error {
			return nil
		},
	}
	events := []domain.WWDCEvent{
		{
			Title: "WWDC24",
			Year:  1999,
			CoverURL: "https://example.com/wwdc24.jpg",
			Videos: []domain.WWDCVideo{},
		},
	}
	sut := NewMarkdownExporter(&fs)
	err := sut.Export(events)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	expectedErrorMessage := "Unable to build year"
	if err.Error() != expectedErrorMessage {
		t.Fatalf("expected error message to be %s, got %s", expectedErrorMessage, err.Error())
	}
}

func TestMarkdownExportMakeDirError(t *testing.T) {
	errMakeDir := errors.New("some error creating directory")
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			return errMakeDir
		},
		writeFileFunc: func(path string, data []byte) error {
			return nil
		},
	}
	events := []domain.WWDCEvent{
		{
			Title: "WWDC24",
			Year:  2024,
			CoverURL: "https://example.com/wwdc24.jpg",
			Videos: []domain.WWDCVideo{},
		},
	}
	sut := NewMarkdownExporter(&fs)
	err := sut.Export(events)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, errMakeDir) {
		t.Fatalf("expected error to be %v, got %v", errMakeDir, err)
	}
}

func TestMarkdownExportWriteFileError(t *testing.T) {
	errWriteFile := errors.New("some error writing file")
	fs := mockFileSystem{
		makeDirFunc: func(path string) error {
			return nil
		},
		writeFileFunc: func(path string, data []byte) error {
			return errWriteFile
		},
	}
	events := []domain.WWDCEvent{
		{
			Title: "WWDC24",
			Year:  2024,
			CoverURL: "https://example.com/wwdc24.jpg",
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
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, errWriteFile) {
		t.Fatalf("expected error to be %v, got %v", errWriteFile, err)
	}
}
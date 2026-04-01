package exporter

import (
	"errors"
	"strconv"
	"strings"
	"path/filepath"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

type MarkdownExporter struct {
	fs domain.FileSystem
	omitTitle bool
}

func NewMarkdownExporter(fs domain.FileSystem, omitTitle bool) *MarkdownExporter {
	return &MarkdownExporter{fs: fs, omitTitle: omitTitle}
}

func (m *MarkdownExporter) Export(events []domain.WWDCEvent) error {
	for _, event := range events {
		year, found := strings.CutPrefix(strconv.Itoa(event.Year), "20")
		if !found {
			return errors.New("Unable to build year")
		}
		err := m.fs.MakeDir("WWDC" + year)
		if err != nil {
			return err
		}
		for _, video := range event.Videos {
			sanitizedTitle := strings.ReplaceAll(video.Title, "/", "_")
			path := filepath.Join("WWDC" + year, sanitizedTitle + ".md")
			content := ""
			if !m.omitTitle {
				content = content + "# " + video.Title + "\n\n"
			}
			content = content + "<video controls src=\"" + video.VideoURL + "\"></video>\n\n" + video.Content
			data := []byte(content)
			err := m.fs.WriteFile(path, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

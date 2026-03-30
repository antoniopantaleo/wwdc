package exporter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

type MarkdownExporter struct {
	fs domain.FileSystem
}

func NewMarkdownExporter(fs domain.FileSystem) *MarkdownExporter {
	return &MarkdownExporter{fs: fs}
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
			path := "WWDC" + year + "/" + video.Title + ".md"
			content := "# " + video.Title + "\n\n" + "<video controls src=\"" + video.VideoURL + "\"></video>\n\n" + video.Content
			data := []byte(content)
			err := m.fs.WriteFile(path, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

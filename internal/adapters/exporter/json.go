package exporter

import (
	"encoding/json"
	"io"

	"github.com/antoniopantaleo/wwdc/internal/domain"
)

type JSONWWDCEvent struct {
	Title    string          `json:"title"`
	Year     int             `json:"year"`
	CoverURL string          `json:"coverUrl"`
	Videos   []JSONWWDCVideo `json:"videos"`
}

type JSONWWDCVideo struct {
	Title    string `json:"title,omitempty"`
	VideoURL string `json:"videoUrl"`
	Content  string `json:"content"`
}

type jsonExportData struct {
	Events []JSONWWDCEvent `json:"events"`
}

type JSONExporter struct {
	writer io.Writer
	omitTitle bool
}

func NewJSONExporter(writer io.Writer, omitTitle bool) *JSONExporter {
	return &JSONExporter{writer: writer, omitTitle: omitTitle}
}

func (e *JSONExporter) Export(events []domain.WWDCEvent) error {
	jsonEvents := make([]JSONWWDCEvent, len(events))
	for i, event := range events {
		jsonEvents[i] = convertToJSONWWDCEvent(event, e.omitTitle)
	}
	exportData := jsonExportData{
		Events: jsonEvents,
	}
	encoder := json.NewEncoder(e.writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportData); err != nil {
		return err
	}
	return nil
}

func convertToJSONWWDCEvent(event domain.WWDCEvent, omitTitle bool) JSONWWDCEvent {
	jsonVideos := make([]JSONWWDCVideo, len(event.Videos))
	for i, video := range event.Videos {
		jsonVideos[i] = convertToJSONWWDCVideo(video, omitTitle)
	}
	return JSONWWDCEvent{
		Title:    event.Title,
		Year:     event.Year,
		CoverURL: event.CoverURL,
		Videos:   jsonVideos,
	}
}

func convertToJSONWWDCVideo(video domain.WWDCVideo, omitTitle bool) JSONWWDCVideo {
	title := video.Title
	if omitTitle {
		title = ""
	}
	return JSONWWDCVideo{
		Title:    title,
		VideoURL: video.VideoURL,
		Content:  video.Content,
	}
}

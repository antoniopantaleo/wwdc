package cmd

import (
	"fmt"

	"github.com/antoniopantaleo/wwdc/internal/adapters/exporter"
	"github.com/antoniopantaleo/wwdc/internal/adapters/filesystem"
	"github.com/antoniopantaleo/wwdc/internal/adapters/scraper"
	"github.com/antoniopantaleo/wwdc/internal/domain"
	"github.com/antoniopantaleo/wwdc/internal/usecases"
	"github.com/spf13/cobra"
)

func NewExportCommand() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use: "export",
		Short: "Scrape and export WWDC session videos",
		Long: `
  Scrape WWDC session videos from the Apple Developer website and export them in the specified format.

  Supported formats are:
  - json: Exports the data as JSON to standard output.
  - markdown|md: Exports the data as Markdown files, creating a directory for each year and a Markdown file for each session video.
`,
		Example: `
  wwdc export --format json
  wwdc export -f md
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			scraper := &scraper.StubScraper{
				Events: []domain.WWDCEvent{
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
		},
			}
			var exp domain.WWDCExporter
			switch format {
			case "json":
				exp = exporter.NewJSONExporter(cmd.OutOrStdout())
			case "markdown", "md":
				fs := filesystem.NewOSFileSystem("./WWDC")
				exp = exporter.NewMarkdownExporter(fs)
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}
			usecase := usecases.NewScrapeAndExportUseCase(scraper, exp)
			return usecase.Execute()
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "", "Export format")
	cmd.MarkFlagRequired("format")
	return cmd
}

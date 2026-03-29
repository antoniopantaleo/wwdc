package cmd

import (
	"fmt"

	"github.com/antoniopantaleo/wwdc/internal/adapters/exporter"
	"github.com/antoniopantaleo/wwdc/internal/adapters/scraper"
	"github.com/antoniopantaleo/wwdc/internal/domain"
	"github.com/antoniopantaleo/wwdc/internal/usecases"
	"github.com/spf13/cobra"
)

type stubFileSystem struct {
	makeDirFunc   func(path string) error
	writeFileFunc func(path string, data []byte) error
}

func (s *stubFileSystem) MakeDir(path string) error {
	return s.makeDirFunc(path)
}

func (s *stubFileSystem) WriteFile(path string, data []byte) error {
	return s.writeFileFunc(path, data)
}

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
						Title: "WWDC26",
						Year: 2026,
						CoverURL: "https://example.com/wwdc26.jpg",
						Videos: []domain.WWDCVideo{
							{
								Title: "Session 1",
								VideoURL: "https://example.com/session1.mp4",
								Content: "This is the content of session 1.",
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
				fs := &stubFileSystem{
					makeDirFunc: func(path string) error {
						fmt.Printf("MakeDir called with path: %s\n", path)
						return nil
					},
					writeFileFunc: func(path string, data []byte) error {
						fmt.Printf("WriteFile called with path: %s and data: %s\n", path, string(data))
						return nil
					},
				}
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

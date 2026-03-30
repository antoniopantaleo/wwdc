package cmd

import (
	"fmt"

	"github.com/antoniopantaleo/wwdc/internal/adapters/exporter"
	"github.com/antoniopantaleo/wwdc/internal/domain"
	"github.com/antoniopantaleo/wwdc/internal/usecases"
	"github.com/spf13/cobra"
)

func NewExportCommand(d *dependencies) *cobra.Command {
	var format string
	var omitTitle bool
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
			scraper := d.Scraper
			var exp domain.WWDCExporter
			switch format {
			case "json":
				exp = exporter.NewJSONExporter(cmd.OutOrStdout())
			case "markdown", "md":
				fs := d.FileSystem
				exp = exporter.NewMarkdownExporter(fs, omitTitle)
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}
			usecase := usecases.NewScrapeAndExportUseCase(scraper, exp)
			return usecase.Execute()
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "", "Export format. Currently supported formats are: json, markdown (or md).")
	cmd.MarkFlagRequired("format")

	cmd.Flags().BoolVar(&omitTitle, "omit-title", false, "Do not write title")
	return cmd
}

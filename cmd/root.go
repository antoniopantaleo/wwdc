package cmd

import (
	"fmt"
	"os"

	"github.com/antoniopantaleo/wwdc/internal/adapters/exporter"
	"github.com/antoniopantaleo/wwdc/internal/adapters/scraper"
	"github.com/antoniopantaleo/wwdc/internal/domain"
	"github.com/antoniopantaleo/wwdc/internal/usecases"
	"github.com/spf13/cobra"
)

func NewRootCommand() (*cobra.Command, error) {
	var format string
	cmd := &cobra.Command{
		Use: "wwdc",
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "json" {
				return fmt.Errorf("unsupported format: %s", format)
			}
			scraper := &scraper.StubScraper{
				Events: []domain.WWDCEvent{},
			}
			out := os.Stdout
			exporter := exporter.NewJSONExporter(out)
			usecase := usecases.NewScrapeAndExportUseCase(scraper, exporter)
			err := usecase.Execute()
			if err != nil {
				return err
			}
			fmt.Fprintln(out)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "", "Export format (only json is supported)")
	err := cmd.MarkFlagRequired("format")
	return cmd, err
}

package cmd

import (
	"bytes"
	"fmt"

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
			scraper := &scraper.StubScraper{
				Events: []domain.WWDCEvent{},
			}
			var buf bytes.Buffer
			exporter := exporter.NewJSONExporter(&buf)
			usecase := usecases.NewScrapeAndExportUseCase(scraper, exporter)
			err := usecase.Execute()
			if err != nil {
				return err
			}
			fmt.Println(buf.String())
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "", "Export format (only json is supported)")
	err := cmd.MarkFlagRequired("format")
	return cmd, err
}

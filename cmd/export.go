package cmd

import (
	"os"

	"github.com/antoniopantaleo/wwdc/internal/adapters/exporter"
	"github.com/antoniopantaleo/wwdc/internal/usecases"
	"github.com/spf13/cobra"
)

func NewExportCommand(d *dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use: "export",
		Short: "Scrape and export WWDC session videos",
		Long: `
  Scrape WWDC session videos from the Apple Developer website and export them in the specified format.
`,
	}

	cmd.AddCommand(newJSONExportCommand(d))
	cmd.AddCommand(newMarkdownExportCommand(d))

	return cmd
}

func newJSONExportCommand(d *dependencies) *cobra.Command {
	var outputPath string
	cmd := &cobra.Command{
		Use: "json",
		Short: "Short description",
		RunE: func(cmd *cobra.Command, args []string) error {
			scraper := d.Scraper
			writer := cmd.OutOrStdout()
			if outputPath != "" {
				f, err := os.Create(outputPath)
				if err != nil {
					return err
				}
				defer f.Close()
				writer = f
			}
			exporter := exporter.NewJSONExporter(writer)
			usecase := usecases.NewScrapeAndExportUseCase(scraper, exporter)
			return usecase.Execute()
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path")
	return cmd
}

func newMarkdownExportCommand(d *dependencies) *cobra.Command {
	var omitTitle bool
	cmd := &cobra.Command{
		Use: "markdown",
		Short: "Short description",
		Aliases: []string{"md"},
		RunE: func(cmd *cobra.Command, args []string) error {
			scraper := d.Scraper
			fs := d.FileSystem
			exporter := exporter.NewMarkdownExporter(fs, omitTitle)
			usecase := usecases.NewScrapeAndExportUseCase(scraper, exporter)
			return usecase.Execute()
		},
	}
	cmd.Flags().BoolVar(&omitTitle, "omit-title", false, "Do not write title")
	return cmd
}

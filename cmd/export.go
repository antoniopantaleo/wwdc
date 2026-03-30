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
		Short: "Scrape WWDC session videos from the Apple Developer website and export them in the specified format.",
	}

	cmd.AddCommand(newJSONExportCommand(d))
	cmd.AddCommand(newMarkdownExportCommand(d))

	return cmd
}

func newJSONExportCommand(d *dependencies) *cobra.Command {
	var outputPath string
	cmd := &cobra.Command{
		Use: "json",
		Short: "Export events in JSON format.",
		Long: "By default the JSON is printed out to standard output. Use `--output` option to save JSON on disk.",
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
		Short: "Export events in .md files, organized in folders, one per event.",
		Long: "Some markdown editors/visualizers use the file name as the title. You can use `--omit-title` flag to ignore title heading in the resulting markdown files.",
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

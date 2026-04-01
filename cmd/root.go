package cmd

import (
	"github.com/spf13/cobra"
	"github.com/antoniopantaleo/wwdc/internal/domain"
)

type dependencies struct {
	FileSystem domain.FileSystem
	Scraper domain.WWDCScraper
}

var version = "dev"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Version: version,
		Use: "wwdc",
		Short: "A CLI tool to scrape and export WWDC session videos",
	}
	d := createDependencies()
	cmd.AddCommand(NewExportCommand(d))
	return cmd
}

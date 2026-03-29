package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Version: "0.1.0",
		Use: "wwdc",
		Short: "A CLI tool to scrape and export WWDC session videos",
	}
	cmd.AddCommand(NewExportCommand())
	return cmd
}

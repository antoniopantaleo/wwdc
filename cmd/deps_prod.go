//go:build !dev

package cmd

import (
	"github.com/antoniopantaleo/wwdc/internal/adapters/filesystem"
	"github.com/antoniopantaleo/wwdc/internal/adapters/reporter"
	"github.com/antoniopantaleo/wwdc/internal/adapters/scraper"
)

func createDependencies() *dependencies {
	return &dependencies{
		FileSystem: filesystem.NewOSFileSystem("./WWDC"),
		Scraper:    scraper.NewCollyScraper("https://developer.apple.com", reporter.NewStderrReporter()),
	}
}

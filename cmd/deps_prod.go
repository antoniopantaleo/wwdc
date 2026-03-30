//go:build !dev

package cmd

import (
	"github.com/antoniopantaleo/wwdc/internal/domain"
	"github.com/antoniopantaleo/wwdc/internal/adapters/scraper"
	"github.com/antoniopantaleo/wwdc/internal/adapters/filesystem"
)

func createDependencies() *dependencies {
	return &dependencies{
		FileSystem: filesystem.NewOSFileSystem("./WWDC"),
		Scraper: &scraper.StubScraper{
			Events: []domain.WWDCEvent{},
		},
	}
}

//go:build dev

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
			Events: []domain.WWDCEvent{
				{
					Title: "WWDC26",
					Year: 2026,
					CoverURL: "https://cover-url.com",
					Videos: []domain.WWDCVideo{
						{
							Title: "A title",
							VideoURL: "https://a-video.mp4",
							Content: "This is a content",
						},
					},
				},
			},
		},
	}
}

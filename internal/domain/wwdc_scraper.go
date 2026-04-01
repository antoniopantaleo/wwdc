package domain

type WWDCScraper interface {
	Scrape() ([]WWDCEvent, error)
}

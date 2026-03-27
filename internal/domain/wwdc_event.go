package domain

type WWDCEvent struct {
	Title    string
	Year     int
	CoverURL string
	Videos   []WWDCVideo
}

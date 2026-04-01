package domain

type WWDCExporter interface {
	Export(events []WWDCEvent) error
}

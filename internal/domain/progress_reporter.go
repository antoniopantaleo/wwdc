package domain

type ProgressReporter interface {
	Info(message string)
	Warning(message string)
}

package reporter

import (
	"fmt"
	"os"
)

type StderrReporter struct{}

func NewStderrReporter() *StderrReporter {
	return &StderrReporter{}
}

func (r *StderrReporter) Info(message string) {
	fmt.Fprintln(os.Stderr, message)
}

func (r *StderrReporter) Warning(message string) {
	fmt.Fprintf(os.Stderr, "WARNING: %s\n", message)
}

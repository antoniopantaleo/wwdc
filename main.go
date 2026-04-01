package main

import (
	"os"

	"github.com/antoniopantaleo/wwdc/cmd"
)

func main() {
	cmd := cmd.NewRootCommand()
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

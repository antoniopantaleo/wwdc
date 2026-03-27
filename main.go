package main

import (
	"os"

	"github.com/antoniopantaleo/wwdc/cmd"
)

func main() {
	cmd, err := cmd.NewRootCommand()
	if err != nil {
		os.Exit(1)
	}
	err = cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

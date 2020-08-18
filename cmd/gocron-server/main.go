package main

import (
	"github.com/x893675/gocron/cmd/gocron-server/app"
	"os"
)

func main() {
	cmd := app.NewGoCronServerCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

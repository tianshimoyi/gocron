package main

import (
	"github.com/x893675/gocron/cmd/gocron-agent/app"
	"os"
)

func main() {
	cmd := app.NewGoCronAgentCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

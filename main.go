package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	app := &cli.App{
		Name:  "PASS download Service",
		Usage: "Provides HTTP endpoints for looking up DOIs and downloading their manuscripts",
		Commands: []*cli.Command{
			serve(),
		},
	}

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

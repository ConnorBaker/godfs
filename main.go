package main

import (
	"log"
	"os"

	"github.com/connorbaker/godfs/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "gocoding",
		Version: "0.0.1",
		Authors: []*cli.Author{
			{
				Name:  "Connor Baker",
				Email: "connorbaker01@gmail.com",
			},
		},
		Commands: cmd.Commands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

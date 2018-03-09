package main

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "makoto"
	app.Usage = "minimalist migration tool for PostgreSQL"
	app.Commands = []cli.Command{
		{
			Name: "init",
			Action: func(c *cli.Context) error {
				initMigrationDir()
				return nil
			},
		},
		{
			Name: "collect",
			Action: func(c *cli.Context) error {
				collectMigrationScrips()
				return nil
			},
		},
	}

	app.Run(os.Args)
}

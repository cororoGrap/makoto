package main

import (
	"fmt"
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
		{
			Name:  "new",
			Usage: "Create new migration sql script",
			Action: func(c *cli.Context) error {
				if c.NArg() == 1 {
					name := c.Args()[0]
					createNewScript(name)
				} else {
					fmt.Println("Missing file name")
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}

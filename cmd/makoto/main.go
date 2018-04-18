package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cororoGrap/makoto"
	"github.com/cororoGrap/makoto/cmd/makoto/db"
	"github.com/olekukonko/tablewriter"
	cli "gopkg.in/urfave/cli.v1"
)

const version = "0.0.1"

var (
	database   string
	configPath string
)

func main() {

	app := cli.NewApp()
	app.Name = "makoto"
	app.Version = version
	app.Usage = "minimalist migration tool for PostgreSQL"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "database",
			Destination: &database,
		},
		cli.StringFlag{
			Name:        "config",
			Destination: &configPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "version",
			Action: func(c *cli.Context) error {
				fmt.Println("makoto version: ", version)
				return nil
			},
		},
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
		{
			Name: "status",
			Action: func(c *cli.Context) error {
				configureDBUri()
				db := db.ConnectPostgres(database)
				r, err := makoto.GetAllRecords(db)
				if err != nil {
					panic(err)
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Version", "Script", "Create Date"})
				for _, record := range r {
					date := record.CreatedAt.Format(time.RFC3339)
					table.Append([]string{record.Version, record.Filename, date})
				}
				table.Render()
				return nil
			},
		},
		{
			Name: "up",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "version",
				},
			},
			Action: func(c *cli.Context) error {
				configureDBUri()
				db := db.ConnectPostgres(database)
				collection := processMigrationCollection(getMigrationDir())
				migrator := makoto.GetMigrator(db, collection)

				version := c.String("version")
				if len(version) == 0 {
					migrator.Up()
				} else {
					migrator.EnsureSchema(version)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func configureDBUri() {
	if len(database) == 0 {
		err := loadDBJson()
		if err != nil {
			panic(err)
		}
	}
}

func getConfigPath() string {
	if len(strings.TrimSpace(configPath)) == 0 {
		return getDefaultConfigPath()
	}
	fmt.Println("this is config path: ", configPath)
	return configPath
}

func loadDBJson() error {
	path := getConfigPath()

	file, err := os.Open(path)
	logError(err)

	config := dbConfig{}
	configSt, err := ioutil.ReadAll(file)
	err = json.Unmarshal(configSt, &config)
	logError(err)

	if len(config.Database) == 0 || config.Database == "PostgreSQL" {
		pg := config.PostgreSQL
		database = fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
			pg.User, pg.Password, pg.Host, pg.Port, pg.DBName)
	} else {
		panic("Unsupported database")
	}

	return nil
}

func getDefaultConfigPath() string {
	return filepath.Join(getMigrationDir(), "config.json")
}

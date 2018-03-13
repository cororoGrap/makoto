package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const collectionFilename = "collection.go"

func GenerateCollection(path string) {
	fmt.Println("Collect migration scripts:")

	buffer := bytes.NewBuffer(nil)
	fmt.Fprint(buffer, `package migration

import "github.com/cororoGrap/makoto"

func GetCollection() []makoto.MigrateStatement {
	return []makoto.MigrateStatement{
	`)

	collection := processMigrationCollection(path)
	migration := collection.Head()
	for {
		st := migration.Statement()
		upSt, _ := json.Marshal(st.UpStatement)
		downSt, _ := json.Marshal(st.DownStatement)

		fmt.Fprintf(buffer, `
		{"%v", "%v", %v, %v, "%v"},
		`, st.Version, st.Filename, string(upSt), string(downSt), st.Checksum)

		fmt.Printf("%v\n", st.Filename)

		if migration.Next() != nil {
			migration = migration.Next()
			continue
		}
		break
	}

	fmt.Fprint(buffer, `
	}
}`)

	dest := filepath.Join(path, collectionFilename)
	if err := ioutil.WriteFile(dest, buffer.Bytes(), 0644); err != nil {
		panic(err)
	}
}

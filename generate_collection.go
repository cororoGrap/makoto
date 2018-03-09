package makoto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const collectionFilename = "collection.go"

func GenerateCollection(path string) {
	buffer := bytes.NewBuffer(nil)
	fmt.Fprint(buffer, `package migration

import "github.com/cororoGrap/makoto"

func GetCollection() []MigrateStatement {
	return []MigrateStatement{
	`)

	collection := processMigrationCollection(path)
	migration := collection.head
	for {
		st := migration.statement
		upSt, _ := json.Marshal(st.UpStatement)
		downSt, _ := json.Marshal(st.DownStatement)

		fmt.Fprintf(buffer, `
		{"%v", "%v", %v, %v, "%v"},
		`, st.Version, st.Filename, string(upSt), string(downSt), st.Checksum)

		if migration.nextNode != nil {
			migration = migration.nextNode
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

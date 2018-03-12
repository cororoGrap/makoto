package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cororoGrap/makoto"
)

const SQLFileExtension = ".sql"

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func processMigrationCollection(path string) *makoto.MigrationCollection {
	files, err := readSQLMigrationScript(path)
	logError(err)

	collection := &makoto.MigrationCollection{}
	for _, f := range files {
		fullPath := filepath.Join(path, f.Name())
		file, err := os.Open(fullPath)
		logError(err)

		migration, err := parseMigration(file)
		logError(err)

		migration.Filename = f.Name()
		migration.Version = parseFilenameVersion(f.Name())

		collection.Add(migration)
	}

	return collection
}

func readSQLMigrationScript(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(path)
	logError(err)

	files, err := dir.Readdir(0)
	logError(err)

	result := []os.FileInfo{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if filepath.Ext(f.Name()) != SQLFileExtension {
			continue
		}
		result = append(result, f)
	}
	return result, nil
}

func parseFilenameVersion(filename string) string {
	r, err := regexp.Compile("v[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return r.FindString(filename)
}

func parseMigration(r io.Reader) (*makoto.MigrateStatement, error) {
	var buf bytes.Buffer
	isDown := false

	migration := makoto.MigrateStatement{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		// line cannot be longer than 65536 characters
		line := scanner.Text()

		buf.WriteString(line)

		if strings.HasPrefix(line, "-- Down") {
			isDown = true
			continue
		}
		if strings.HasPrefix(line, "-- Up") {
			isDown = false
			continue
		}

		if isDown {
			migration.DownStatement += line + "\n"
		} else {
			migration.UpStatement += line + "\n"
		}
	}

	migration.Checksum = getMD5SumString(buf.Bytes())

	return &migration, nil
}

func getMD5SumString(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

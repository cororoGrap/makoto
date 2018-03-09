package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cororoGrap/makoto"
)

const migrationPath = "migration"

func initMigrationDir() {
	migrationPath, err := getMigrationDir()
	if err != nil {
		return
	}
	os.Mkdir(migrationPath, os.ModePerm)
}

func collectMigrationScrips() {
	migrationPath, err := getMigrationDir()
	if err != nil {
		log.Fatal("Migration not yet initialized")
	}

	makoto.GenerateCollection(migrationPath)
}

func getMigrationDir() (string, error) {
	dir := currentDir()
	fullPath := filepath.Join(dir, migrationPath)
	return fullPath, nil
}

func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cororoGrap/makoto"
)

const migrationPath = "migration"

func initMigrationDir() {
	dir := currentDir()
	path := filepath.Join(dir, migrationPath)
	if exists(path) {
		fmt.Println("Migration directory already exists")
		return
	}
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		fmt.Println("Created migration directory")
	}
}

func collectMigrationScrips() {
	migrationPath := getMigrationDir()
	GenerateCollection(migrationPath)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Fatal(err)
	return true
}

func getMigrationDir() string {
	dir := currentDir()
	if strings.HasSuffix(dir, migrationPath) {
		return dir
	}
	fullPath := filepath.Join(dir, migrationPath)
	if exists(fullPath) {
		return fullPath
	}
	log.Fatal("Unknow migration directory")
	return ""
}

func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func createNewScript(name string) {
	dir := getMigrationDir()
	version := getNewScriptVersion()

	filename := fmt.Sprintf("v%v_%s.sql", version, name)
	fullPath := filepath.Join(dir, filename)
	fmt.Println("Create new migration script: ", filename)
	os.Create(fullPath)
}

func getNewScriptVersion() string {
	collection := initCollection()
	if st := collection.LastStatement(); st != nil {
		v, err := strconv.Atoi(st.Version[1:])
		if err != nil {
			log.Fatal(err)
		}
		v++
		return strconv.Itoa(v)
	}
	return "1"
}

// func displayMigrati

func initCollection() *makoto.MigrationCollection {
	dir := getMigrationDir()
	return processMigrationCollection(dir)
}

package main

type dbConfig struct {
	Database   string
	PostgreSQL postgres
}

type postgres struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  bool
}

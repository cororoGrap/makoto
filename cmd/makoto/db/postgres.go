package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectPostgres(uri string) *sqlx.DB {
	con, err := sqlx.Connect("postgres", uri)
	if err != nil {
		panic(err)
	}
	return con
}

package makoto

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type migrator struct {
	db         *sqlx.DB
	collection *migrationCollection
}

func New(db *sql.DB, driverName string) *migrator {
	xDB := sqlx.NewDb(db, driverName)

	err := createSchemaVersionTable(xDB)
	if err != nil {
		log.Fatal(err)
	}

	return &migrator{
		db: xDB,
	}
}

func (m *migrator) AddCollection(sts []MigrateStatement) {
	if m.collection != nil {
		return
	}

	collection := migrationCollection{}
	for i := range sts {
		collection.Add(&sts[i])
	}
	m.collection = &collection
}

func (m *migrator) Up() {
	tx := m.db.MustBegin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Rollback migration, Error: ", r)
		}
	}()
	migrate := MigrateUp{
		tx:         tx,
		collection: m.collection,
	}
	migrate.Up()
	tx.Commit()
}

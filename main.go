package makoto

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type migrator struct {
	db         *sqlx.DB
	collection *MigrationCollection
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

	collection := MigrationCollection{}
	for i := range sts {
		collection.Add(&sts[i])
	}
	m.collection = &collection
}

func (m *migrator) EnsureSchema(targetVersion string) {
	record := m.getCurrentRecord()
	if record.Version == targetVersion {
		return
	}
	if record.Version < targetVersion {
		node := m.getCurrentNode()
		m.upto(node, targetVersion)
	}
}

func (m *migrator) getCurrentNode() *migrationItem {
	record := m.getCurrentRecord()
	return m.collection.Find(record.Version)
}

func (m *migrator) upto(currentNode *migrationItem, targetVersion string) {
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

	migrate.UpTo(currentNode, targetVersion)
	tx.Commit()
}

func (m *migrator) Up() {
	node := m.getCurrentNode()
	lastVersion := m.collection.LastStatement().Version
	if node.statement.Version < lastVersion {
		m.upto(node, lastVersion)
	}
}

func (m *migrator) getCurrentRecord() *MigrationRecord {
	tx := m.db.MustBegin()

	record, err := getLastRecord(tx)
	if err != nil && err != ErrRecordNotFound {
		log.Fatal(err)
	}
	log.Println("record: ", record)

	return record
}

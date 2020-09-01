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

func GetMigrator(db *sqlx.DB, collection *MigrationCollection) *migrator {
	return &migrator{
		db:         db,
		collection: collection,
	}
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

func (m *migrator) GetCollection() *MigrationCollection {
	if m.collection != nil {
		return m.collection
	}
	log.Fatal("Migration collection not found")
	return nil
}

func (m *migrator) SetCollection(sts []MigrateStatement) {
	collection := MigrationCollection{}
	for i := range sts {
		collection.Add(&sts[i])
	}
	m.collection = &collection
}

func (m *migrator) EnsureSchema(targetVersion string) {
	currentNode, err := m.getCurrentNode()

	if err != nil && err != ErrRecordNotFound {
		log.Fatal(err)
	}

	if err == ErrRecordNotFound {
		currentNode = m.GetCollection().Head()
		m.upto(currentNode, targetVersion)
		return
	}

	st := currentNode.Statement()
	if st.Version == targetVersion {
		return
	}
	if v(st.Version) < v(targetVersion) {
		log.Println("start migrate")
		m.upto(currentNode.nextNode, targetVersion)
	}
}

func (m *migrator) getCurrentNode() (*migrationItem, error) {
	record, err := getLastRecord(m.db)
	if err != nil {
		return nil, err
	}
	if v(record.Version) > v(m.GetCollection().LastStatement().Version) {
		return m.GetCollection().Tail(), nil
	}
	return m.GetCollection().Find(record.Version), nil
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
	lastVersion := m.GetCollection().LastStatement().Version
	m.EnsureSchema(lastVersion)
}

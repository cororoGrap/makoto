package makoto

import (
	"time"
)

type MigrationRecord struct {
	ID        int
	Version   string
	Filename  string
	Checksum  string
	CreatedAt time.Time `db:"created_at"`
}

type MigrateStatement struct {
	Version       string
	Filename      string
	UpStatement   string
	DownStatement string
	Checksum      string
}

// a simple sorted linkedlist

type migrationItem struct {
	statement    MigrateStatement
	previousNode *migrationItem
	nextNode     *migrationItem
}

type migrationCollection struct {
	head *migrationItem
}

func (m *migrationCollection) Add(st *MigrateStatement) {
	newItem := &migrationItem{
		statement: *st,
	}

	if m.head == nil {
		m.head = newItem
		return
	}

	migration := m.head
	for {
		if st.Version < migration.statement.Version {
			if migration.previousNode != nil {
				migration.previousNode.nextNode = newItem
			} else {
				m.head = newItem
			}
			migration.previousNode = newItem
			newItem.nextNode = migration
			break
		}
		if migration.nextNode == nil {
			migration.nextNode = newItem
			newItem.previousNode = migration
			break
		}
		migration = migration.nextNode
	}
}

func (m *migrationCollection) Find(version string) *migrationItem {
	for {
		migration := m.head
		if migration == nil {
			return nil
		}
		if migration.statement.Version == version {
			return migration
		}
		migration = migration.nextNode
	}
}

func (m *migrationCollection) FindStatement(version string) *MigrateStatement {
	item := m.Find(version)
	if item == nil {
		return nil
	}
	return &item.statement
}

func (m *migrationCollection) LastStatement() *MigrateStatement {
	if m.head == nil {
		return nil
	}

	migration := m.head
	for {
		if migration.nextNode != nil {
			migration = migration.nextNode
		} else {
			return &migration.statement
		}
	}
}

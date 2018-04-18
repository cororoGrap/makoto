package makoto

import (
	"strconv"
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

func (m *migrationItem) Statement() *MigrateStatement {
	return &m.statement
}

func (m *migrationItem) Next() *migrationItem {
	return m.nextNode
}

func (m *migrationItem) Previous() *migrationItem {
	return m.previousNode
}

type MigrationCollection struct {
	head *migrationItem
}

func (m *MigrationCollection) Head() *migrationItem {
	return m.head
}

func (m *MigrationCollection) Add(st *MigrateStatement) {
	newItem := &migrationItem{
		statement: *st,
	}

	if m.head == nil {
		m.head = newItem
		return
	}

	migration := m.head
	for {
		if v(st.Version) < v(migration.statement.Version) {
			if migration.previousNode != nil {
				migration.previousNode.nextNode = newItem
				newItem.previousNode = migration.previousNode
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

func (m *MigrationCollection) Find(version string) *migrationItem {
	migration := m.head
	for {
		if migration == nil {
			return nil
		}
		if migration.statement.Version == version {
			return migration
		}
		migration = migration.nextNode
	}
}

func (m *MigrationCollection) FindStatement(version string) *MigrateStatement {
	item := m.Find(version)
	if item == nil {
		return nil
	}
	return &item.statement
}

func (m *MigrationCollection) LastStatement() *MigrateStatement {
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

func v(v string) int {
	val, err := strconv.Atoi(v[1:])
	if err != nil {
		panic(err)
	}
	return val
}

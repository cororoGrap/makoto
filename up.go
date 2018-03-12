package makoto

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type MigrateUp struct {
	tx         *sqlx.Tx
	collection *MigrationCollection
}

func (m *MigrateUp) UpTo(node *migrationItem, targetVersion string) {
	tx := m.tx

	currentNode := node
	for {
		statement := currentNode.statement
		if statement.Version <= targetVersion {
			_, err := tx.Exec(statement.UpStatement)
			if err != nil {
				log.Println("Fail to run migration script: ", statement.Filename)
				log.Fatal(err)
			}
			log.Println("Migrated script: ", statement.Filename)
			err = addRecord(tx, statement.Version, statement.Filename, statement.Checksum)
			if err != nil {
				log.Fatal(err)
			}
			if currentNode.nextNode == nil {
				break
			}
			currentNode = currentNode.nextNode
		} else {
			break
		}
	}
}

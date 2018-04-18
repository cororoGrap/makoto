package makoto

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

func createSchemaVersionTable(db *sqlx.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS schema_version (
		id serial PRIMARY KEY,
		version text,
		filename text,
		checksum text,
		created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
	)
	`
	tx := db.MustBegin()
	_, err := tx.Exec(sql)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func addRecord(tx *sqlx.Tx, version, filename, checksum string) error {
	sql := `
	INSERT INTO schema_version (version, filename, checksum) 
	VALUES (:version, :filename, :checksum)
	`
	_, err := tx.NamedExec(sql, map[string]interface{}{
		"version":  version,
		"filename": filename,
		"checksum": checksum,
	})
	if err != nil {
		return err
	}

	return nil
}

func getLastRecord(db *sqlx.DB) (*MigrationRecord, error) {
	sql := `
	SELECT * FROM schema_version
	ORDER by id desc
	LIMIT 1
	`
	row, err := db.Queryx(sql)
	if err != nil {
		return nil, err
	}

	record := MigrationRecord{}
	if row.Next() {
		err = row.StructScan(&record)
		if err != nil {
			return nil, err
		}
		return &record, nil
	}
	return nil, ErrRecordNotFound
}

func GetAllRecords(db *sqlx.DB) ([]MigrationRecord, error) {
	sql := `
	SELECT * FROM schema_version
	`
	tx := db.MustBegin()
	rows, err := tx.Queryx(sql)
	if err != nil {
		return nil, err
	}

	records := []MigrationRecord{}
	for rows.Next() {
		record := MigrationRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, tx.Commit()
}

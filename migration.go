package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Migration struct {
	db             *sql.DB
	currentVersion int
}

func NewMigration(db *sql.DB) *Migration {
	return &Migration{
		db:             db,
		currentVersion: GetLastMigrationVersion(db),
	}
}

func (m *Migration) Up() {
	m.currentVersion++

	migrationFile := GetMigrationFile("up", m.currentVersion)

	m.Migrate(migrationFile)

	if _, err := m.db.Exec("INSERT INTO migrations (version) VALUES ($1)", m.currentVersion); err != nil {
		log.Fatal(errors.Join(ErrDBExecutionFailed, err))
	}

}

func (m *Migration) Down() {
	migrationFile := GetMigrationFile("down", m.currentVersion)

	m.Migrate(migrationFile)

	if _, err := m.db.Exec("DELETE FROM migrations WHERE version = $1", m.currentVersion); err != nil {
		log.Fatal(errors.Join(ErrDBExecutionFailed, err))
	}

	m.currentVersion--
}

func (m *Migration) Migrate(migrationFile string) {
	fmt.Println("Running file: ", migrationFile)

	query := GetQuery(migrationFile)

	m.ExecuteMigration(query)
}

func (m *Migration) ExecuteMigration(query string) {
	if _, err := m.db.Exec("BEGIN"); err != nil {
		log.Fatal(errors.Join(ErrDBExecutionFailed, err))
	}

	if _, err := m.db.Exec(query); err != nil {

		if _, err := m.db.Exec("ROLLBACK"); err != nil {
			log.Fatal(errors.Join(ErrDBExecutionFailed, err))
		}

		log.Fatal(errors.Join(ErrDBExecutionFailed, err))
	}

	if _, err := m.db.Exec("COMMIT"); err != nil {
		log.Fatal(errors.Join(ErrDBExecutionFailed, err))
	}
}

package main

import (
	"bufio"
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	db := SetupDB()
	defer db.Close()

	command := GetCommand()
	migration := NewMigration(db)
	migrationVersionLeft := GetMigrationVersionLeft(migration.currentVersion, command)

	if migrationVersionLeft == 0 {
		log.Println("Migration complete")
		return
	}

	for range migrationVersionLeft {
		if command == "up" {
			migration.Up()
		}

		if command == "down" {
			migration.Down()
		}
	}

	log.Println("Migration complete")

}

func GetCommand() string {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal(ErrCommandNotFound)
	}

	return args[0]
}

func GetQuery(filepath string) string {
	var query string

	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(errors.Join(ErrNoMigrationFile, err))
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		query += scanner.Text() + "\n"
	}

	return query
}

func GetMigrationFile(command string, version int) string {
	matches, err := filepath.Glob("migrations/" + strconv.Itoa(version) + "*/" + "*-" + command + ".sql")

	if err != nil {
		log.Fatal(errors.Join(ErrNoMigrationFile, err))
	}

	if len(matches) == 0 {
		log.Fatal(ErrNoMigrationFile)
	}

	return matches[0]
}
func GetMigrationVersionLeft(currentVersion int, command string) int {
	if command == "down" {
		return currentVersion
	}

	dirs, err := os.ReadDir("migrations")

	if err != nil {
		log.Fatal(ErrNoMigrationFolder)
	}

	return len(dirs) - currentVersion
}

func GetLastMigrationVersion(db *sql.DB) int {
	var lastMigration int

	row := db.QueryRow("SELECT version FROM migrations ORDER BY id DESC LIMIT 1")

	row.Scan(&lastMigration)

	return lastMigration
}

func SetupDB() *sql.DB {
	url := GetDatabaseURL()

	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal(errors.Join(ErrDBSetupFailed, err))
	}

	if err := db.Ping(); err != nil {
		log.Fatal(errors.Join(ErrDBSetupFailed, err))
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id SERIAL PRIMARY KEY, version INT);"); err != nil {
		log.Fatal(errors.Join(ErrDBSetupFailed, err))
	}

	return db
}

func GetDatabaseURL() string {
	url := os.Getenv("DATABASE_URL")

	if url != "" {
		return url
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(ErrEnvNotFound)
	}

	url = os.Getenv("DATABASE_URL")

	if url == "" {
		log.Fatal(ErrEnvNotFound)
	}

	return url
}

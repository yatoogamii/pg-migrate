package main

import "fmt"

var (
	ErrEnvNotFound       = fmt.Errorf("environment variable not found, you must set DATABASE_URL")
	ErrCommandNotFound   = fmt.Errorf("command argument not found, you must specify up or down")
	ErrDBSetupFailed     = fmt.Errorf("failed to setup to database")
	ErrNoMigrationFile   = fmt.Errorf("no migration file found")
	ErrNoMigrationFolder = fmt.Errorf("no migrations folder found in current directory")
	ErrDBExecutionFailed = fmt.Errorf("database execution failed")
)

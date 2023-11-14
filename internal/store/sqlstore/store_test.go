package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL   string
	migrationPath string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("CENARIUS_DATABASEDSN")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=cenarius_test sslmode=disable"
	}
	migrationPath = os.Getenv("CENARIUS_MIGRATION_PATH")
	if migrationPath == "" {
		migrationPath = "migrations/"
	}
	os.Exit(m.Run())
}

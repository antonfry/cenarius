package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("CENARIUS_DATABASEDSN")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=cenarius_test sslmode=disable"
	}
	os.Exit(m.Run())
}

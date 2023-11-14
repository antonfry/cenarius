package sqlstore_test

import (
	"cenarius/internal/store/sqlstore"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestNewPGConn(t *testing.T) {
	tests := []struct {
		name        string
		databaseDsn string
		wantErr     bool
	}{
		{
			name:        "Valid",
			databaseDsn: os.Getenv("CENARIUS_DATABASEDSN"),
			wantErr:     false,
		},
		{
			name:        "Wrong Host",
			databaseDsn: "host=123 dbname=cenarius_test sslmode=disable",
			wantErr:     true,
		},
		{
			name:        "Wrong DB",
			databaseDsn: "host=localhost dbname=nonexist sslmode=disable",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sqlstore.NewPGConn(tt.databaseDsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPGConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

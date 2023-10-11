package sqlstore

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestStore(t *testing.T, databaseURL string) (*Store, func(...string)) {
	t.Helper()
	db, err := NewPGConn(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	store := NewStore(db)
	return store, func(tables ...string) {
		if len(tables) > 0 {
			ctx := context.Background()
			if _, err := store.db.ExecContext(ctx, fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}
		store.db.Close()
	}
}

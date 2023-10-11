package sqlstore

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPGConn(databaseDsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", databaseDsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}

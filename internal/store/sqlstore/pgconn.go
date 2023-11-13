package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

func NewPGConn(databaseDsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", databaseDsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	if err := migrateSQL(conn); err != nil {
		log.Error("Migration Fail: ", err.Error())
		return nil, err
	}
	fmt.Println("NegPGCon")
	return conn, nil
}

func migrateSQL(conn *sql.DB) error {
	driver, err := pgx.WithInstance(conn, &pgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"pgx", driver)
	if err != nil {
		return err
	}
	m.Up()
	return nil
}

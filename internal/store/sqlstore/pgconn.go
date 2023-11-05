package sqlstore

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

const schemaSQL = `
begin;
select pg_advisory_xact_lock(12345);
CREATE TABLE IF NOT EXISTS users(
    "id" bigserial not null primary key,
    "login" varchar not null unique,
    "encrypted_password" varchar not null
);

CREATE TABLE IF NOT EXISTS LoginWithPassword(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "login" varchar,
    "password" varchar,
    "created_at" timestamp default NOW()
);

CREATE TABLE IF NOT EXISTS CreditCard(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "owner_name" varchar,
    "owner_last_name" varchar,
    "number" varchar,
    "cvc" varchar,
    "created_at" timestamp default NOW()
);

CREATE TABLE IF NOT EXISTS SecretText(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "text" text,
    "created_at" timestamp default NOW()
);

CREATE TABLE IF NOT EXISTS SecretFile(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "path" varchar,
    "created_at" timestamp default NOW()
);
commit;`

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
	return conn, nil
}

func migrateSQL(conn *sql.DB) error {
	if _, err := conn.Exec(schemaSQL); err != nil {
		return err
	}
	return nil
}

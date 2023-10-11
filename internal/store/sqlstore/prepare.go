package sqlstore

import "database/sql"

const craeteUserTableSQL = "CREATE TABLE IF NOT EXISTS users ( " +
	"id bigserial not null primary key," +
	"login varchar not null unique," +
	"encrypted_password varchar not null" +
	");"

const craeteOrderTableSQL = "CREATE TABLE IF NOT EXISTS orders ( " +
	"id bigserial not null primary key," +
	"number varchar," +
	"status varchar," +
	"user_id bigserial," +
	"accrual double precision," +
	"uploaded_at timestamp default NOW()" +
	");"

const craeteAcccountTableSQL = "CREATE TABLE IF NOT EXISTS accounts ( " +
	"id bigserial not null primary key," +
	"user_id bigserial," +
	"current double precision," +
	"withdrawn double precision," +
	"created_at timestamp default NOW()" +
	");"

const craeteWithdrawalTableSQL = "CREATE TABLE IF NOT EXISTS withdrawals ( " +
	"id bigserial not null primary key," +
	"order_number varchar," +
	"sum double precision," +
	"processed_at timestamp default NOW()" +
	");"

func PrepareSQLTables(db *sql.DB) error {
	_, err := db.Exec(craeteUserTableSQL)
	if err != nil {
		return err
	}
	_, err = db.Exec(craeteOrderTableSQL)
	if err != nil {
		return err
	}
	_, err = db.Exec(craeteAcccountTableSQL)
	if err != nil {
		return err
	}
	_, err = db.Exec(craeteWithdrawalTableSQL)
	if err != nil {
		return err
	}
	return nil
}

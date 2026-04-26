package store

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenMySQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(15)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

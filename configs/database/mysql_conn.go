package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(`mysql`, dataSourceName)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

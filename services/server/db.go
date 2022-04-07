package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func connectDB(conninfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

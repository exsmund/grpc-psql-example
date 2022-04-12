package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go"
	_ "github.com/lib/pq"
)

type DBManager struct {
	ok bool
	db *sql.DB
}

func (m *DBManager) Connect(conf *PostgresConfig) {
	var conninfo string = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		conf.host, conf.port, conf.db, conf.username, conf.password)

	var err error
	defer func() {
		m.ok = err == nil
	}()

	err = retry.Do(
		func() error {
			log.Println("Trying to connect to the DB...")
			db, err := sql.Open("postgres", conninfo)
			if err != nil {
				return err
			}

			err = db.Ping()
			if err != nil {
				db.Close()
				return err
			}

			m.db = db
			return nil
		},
		retry.Attempts(5),
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
	)
}

func (m *DBManager) Close() {
	m.db.Close()
}

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

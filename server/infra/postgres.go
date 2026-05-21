package infra

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to postgress. Error:", err.Error())
	}

	return db
}

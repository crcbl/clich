package resources

import (
	"database/sql"
	"log"

    _ "github.com/lib/pq"
)

var database *sql.DB

func GetDB() (*sql.DB, error) {
	if database != nil {
		return database, nil
	}

	database, err := sql.Open("postgres", ("DB_URL"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return database, nil
}

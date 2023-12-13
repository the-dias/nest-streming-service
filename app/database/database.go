package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	user     string
	password string
	dbname   string
	host     string
	port     int
	conn     *sql.DB
}

const driverName = "postgres"

func New(user, password, dbname, host string, port int) *Database {

	db := Database{
		user:     user,
		password: password,
		dbname:   dbname,
		host:     host,
		port:     port,
	}

	return &db
}

func (database *Database) Open() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		database.host, database.port, database.user, database.password, database.dbname)

	db, err := sql.Open(driverName, psqlInfo)
	if err != nil {
		return nil, err
	}

	database.conn = db

	return db, nil
}

func (database *Database) Close() {
	if database.conn != nil {
		if err := database.conn.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}
}

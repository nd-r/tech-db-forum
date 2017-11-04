package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
)

const pqErrViolatesContraint = "23505"

var db *sqlx.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "docker"
	password = "docker"
	dbname   = "docker"

	schema   = "./database/schema.sql"
)

// DBPoolInit initializes sqlx db pool
func DBPoolInit() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db = sqlx.MustConnect("postgres", psqlInfo)

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	// db.SetConnMaxLifetime(0)
}

// InitDBSchema inits tables, indexes, etc.
func InitDBSchema() {
	buf, err := ioutil.ReadFile(schema)
	if err != nil {
		log.Fatalln(err)
	}

	schema := string(buf)
	db.MustExec(schema)
}

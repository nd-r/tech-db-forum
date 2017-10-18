package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	ioutil "io/ioutil"
	"log"
)

// DB - database connection pool
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(0)
}


// InitDBSchema inits tables, indexes, etc.
func InitDBSchema() {
	buf, err := ioutil.ReadFile(schema)

	if err != nil {
		log.Fatal(err)
	}

	schema := string(buf)

	_, err = db.Query(schema)

	if err != nil {
		log.Fatal(err)
	}
}

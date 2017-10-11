package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	ioutil "io/ioutil"
	"log"
)

// DB - database connection pool
var DB *sqlx.DB

const (
	host        = "localhost"
	port        = 25432
	user        = "docker"
	password    = "docker"
	dbname      = "docker"
	schema = "./database/schema.sql"
)

// DBPoolInit initializes sqlx db pool
func DBPoolInit() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	DB, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	DB.SetMaxIdleConns(100)
	DB.SetMaxOpenConns(100)
	DB.SetConnMaxLifetime(0)
}

// InitDBSchema inits tables, indexes, etc.
func InitDBSchema() {
	buf, err := ioutil.ReadFile(schema)

	if err != nil {
		log.Fatal(err)
	}

	schema := string(buf)

	_, err = DB.Query(schema)

	if err != nil{
		log.Fatal(err)
	}
}

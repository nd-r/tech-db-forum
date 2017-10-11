package services

import (
	"log"
	"fmt"
	"github.com/jmoiron/sqlx"
	ioutil "io/ioutil"
)

const (
	host     = "localhost"
	port     = 25432
	user     = "docker"
	password = "docker"
	dbname   = "docker"
	initsqlPath = "resources/init.sql"
)

// DBPoolInit initializes sqlx db pool
func DBPoolInit() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	db.DB.SetMaxIdleConns(1000)
	db.DB.SetMaxOpenConns(1000)
	db.DB.SetConnMaxLifetime(0)

	return db, err
}

// InitDBSchema inits tables, indexes, etc.
func InitDBSchema(db *sqlx.DB) {
	buf, err := ioutil.ReadFile(initsqlPath)
	if err != nil{
		log.Fatal(err)
	}

	schema := string(buf)

	db.MustExec(schema)
}

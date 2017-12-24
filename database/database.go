package database

import (
	"io/ioutil"
	"log"

	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/resources"
)

var db *pgx.ConnPool

const schema = "./resources/schema.sql"

var pgConfig = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	User:     "docker",
	Password: "docker",
	Database: "docker",
}

// InitDBSchema initializes tables, indexes, etc.
func InitDBSchema(conn *pgx.Conn) error {
	tx, err := conn.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback()

	buf, err := ioutil.ReadFile(schema)
	if err != nil {
		log.Fatalln(err)
	}

	schema := string(buf)

	if _, err = tx.Exec(schema); err != nil {
		log.Fatalln(err)
	}

	resources.PrepareForumQueries(tx)
	resources.PrepareForumUsersQueries(tx)
	resources.PreparePostQueries(tx)
	resources.PrepateThreadQueries(tx)
	resources.PrepareUsersQueries(tx)
	resources.PrepareVotesQureies(tx)
	tx.Commit()
	return nil
}

// DBPoolInit initializes pgx db pool
func DBPoolInit() {
	var err error
	db, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgConfig,
		MaxConnections: 50,
		AfterConnect:   InitDBSchema,
	})

	if err != nil {
		log.Fatalln(err)
	}
}

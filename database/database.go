package database

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/resources"
)

var db *pgx.ConnPool

const schema = "./resources/schema.sql"

var pgConfig = pgx.ConnConfig{
	Host:     "/var/run/postgresql/",
	Port:     5432,
	User:     "docker",
	Password: "docker",
	Database: "docker",
}

func VacuumAnalyze() {
	time.Sleep(15 * time.Second)
	log.Println(`VACUUMED`)
	db.Exec("VACUUM ANALYZE")
}

func GoVacuum() {
	go VacuumAnalyze()
}

// InitDBSchema initializes tables, indexes, etc.
func InitDBSchema() {
	tx, err := db.Begin()
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
		log.Println(err)
		tx.Rollback()
	}
	tx.Commit()

	resources.PrepareForumQueries(db)
	resources.PrepareForumUsersQueries(db)
	resources.PreparePostQueries(db)
	resources.PrepateThreadQueries(db)
	resources.PrepareUsersQueries(db)
	resources.PrepareVotesQureies(db)
}

// DBPoolInit initializes pgx db pool
func DBPoolInit() {
	var err error
	db, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgConfig,
		MaxConnections: 50,
	})

	if err != nil {
		log.Fatalln(err)
	}
}

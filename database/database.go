package database

import (
	"io/ioutil"
	"log"

	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/resources"
	"time"
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

func TxMustBegin() *pgx.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	return tx
}

func Vaccuum() {
	time.Sleep(20 * time.Second)
	//go func() {
	db.Exec("CLUSTER forum_users USING forum_users_forum_id_nickname_index")
	db.Exec("CLUSTER users USING users_nickname_index")
	db.Exec("CLUSTER post USING parent_tree_3_1")
	db.Exec("CLUSTER thread USING thread_forum_id_created_index")
	db.Exec("CLUSTER forum USING forum_slug_id_index")
	db.Exec("VACUUM ANALYZE")
	//}()
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

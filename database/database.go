package database

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"log"
)

var db *pgx.ConnPool

const schema = "./database/schema.sql"

var pgConfig = pgx.ConnConfig{
	Host:"localhost",
	Port: 5432,
	User:"docker",
	Password: "docker",
	Database: "docker",
}

var zero = 0

// DBPoolInit initializes sqlx db pool
func DBPoolInit() {
	var err error
	db, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgConfig,
		MaxConnections: 50,
	})

	log.Println(err)

	if err != nil{
		log.Fatalln(err)
	}
}

// InitDBSchema inits tables, indexes, etc.
func InitDBSchema() {
	tx, err := db.Begin()
	if err != nil{
		log.Fatalln(err)
	}

	buf, err := ioutil.ReadFile(schema)
	if err != nil {
		log.Fatalln(err)
	}

	schema := string(buf)

	_, err = tx.Exec(schema)
	if err != nil{
		tx.Rollback()		
		log.Fatalln(err)
	}

	tx.Commit()

	_, err = db.Prepare("getThreadBySlug", getThreadBySlug)
	log.Println(err)
	_, err = db.Prepare("getUserProfileQuery", getUserProfileQuery)
	log.Println(err)
	_, err = db.Prepare("putVoteByThrID", putVoteByThrID)
	log.Println(err)
	_, err = db.Prepare("putVoteByThrSLUG", putVoteByThrSLUG)
	log.Println(err)
	_, err = db.Prepare("insertPost", insertPost)
	if err != nil{
		log.Fatalln(err, 1)
	}
	_, err = db.Prepare("generateNextIDs", generateNextIDs)
	log.Println(err)
	_, err = db.Prepare("getUserProfileQuery", getUserProfileQuery)
	log.Println(err)
	_, err = db.Prepare("selectParentAndParents", selectParentAndParents)
	log.Println(err)
	
	_, err = db.Prepare("insertIntoForumUsers", insertIntoForumUsers)
	log.Println(err)
	_, err = db.Prepare("selectParentAndParents", selectParentAndParents)
	log.Println(err)
	
}


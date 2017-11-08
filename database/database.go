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

	_, err = db.Prepare("InsertIntoThread", createThreadQuery)	
	_, err = db.Prepare("getThreadBySlug", getThreadBySlug)
	_, err = db.Prepare("putVoteByThrID", putVoteByThrID)
	_, err = db.Prepare("putVoteByThrSLUG", putVoteByThrSLUG)
	_, err = db.Prepare("getThreadIdBySlug", getThreadIdBySlug)
	_, err = db.Prepare("getThreadIdById", getThreadIdById)
	_, err = db.Prepare("insertPost", insertPost)
	if err != nil{
		log.Fatalln(err, 1)
	}
	_, err = db.Prepare("insertForumUsers", insertForumUsers)
	_, err = db.Prepare("generateNextIDs", generateNextIDs)
	_, err = db.Prepare("claimInfoWithoutParent", claimInfoWithoutParent)
	_, err = db.Prepare("claimInfoWithParent", claimInfoWithParent)
	_, err = db.Prepare("getUserProfileQuery", getUserProfileQuery)
	_, err = db.Prepare("selectParentAndParents", selectParentAndParents)
	_, err = db.Prepare("getForumIdBySlug", getForumIdBySlug)

}


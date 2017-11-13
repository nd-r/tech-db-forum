package database

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"log"
)

var db *pgx.ConnPool

const schema = "./database/schema.sql"

var pgConfig = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	User:     "docker",
	Password: "docker",
	Database: "docker",
}

var zero = 0

// DBPoolInit initializes sqlx db pool
func DBPoolInit() {
	var err error
	db, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgConfig,
		MaxConnections: 50,
	})

	log.Println(err)

	if err != nil {
		log.Fatalln(err)
	}
}

// InitDBSchema inits tables, indexes, etc.
func InitDBSchema() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	buf, err := ioutil.ReadFile(schema)
	if err != nil {
		log.Fatalln(err)
	}

	schema := string(buf)

	_, err = tx.Exec(schema)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	tx.Commit()

	_, err = db.Prepare("getThreadBySlug", getThreadBySlug)
	log.Println(err)
	_, err = db.Prepare("putVoteByThrID", putVoteByThrID)
	log.Println(err)
	_, err = db.Prepare("putVoteByThrSLUG", putVoteByThrSLUG)
	log.Println(err)
 
	_, err = db.Prepare("generateNextIDs", generateNextIDs)
	log.Println(err)
	_, err = db.Prepare("getUserProfileQuery", getUserProfileQuery)
	log.Println(err)
	_, err = db.Prepare("selectParentAndParents", selectParentAndParents)
	log.Println(err)

	_, err = db.Prepare("selectParentAndParents", selectParentAndParents)
	log.Println(err)

	//SELECTING FORUM DATA 1
	_, err = db.Prepare("selectForumQuery", selectForumQuery)
	if err != nil{
		log.Fatalln(selectForumQuery, err)		
	}

	//INSERTING INTO FORUM 1
	_, err = db.Prepare("createForumQuery", createForumQuery)
	if err != nil{
		log.Fatalln(createForumQuery, err)		
	}
	
	//Claiming full user info 1
	_, err = db.Prepare("claimUserInfo", claimUserInfo)
	if err != nil{
		log.Fatalln(claimUserInfo, err)		
	}

	//Claiming forum id and slug by slug 1
	_, err = db.Prepare("getForumIDAndSlugBySlug", getForumIDAndSlugBySlug)
	if err != nil{
		log.Fatalln(getForumIDAndSlugBySlug, err)		
	}
	//inserting into thread 1
	_, err = db.Prepare("insertIntoThread", insertIntoThread)
	if err != nil{
		log.Fatalln(insertIntoThread, err)		
	}
	//inserting into forum_users 1
	_, err = db.Prepare("insertIntoForumUsers", insertIntoForumUsers)
	if err != nil{
		log.Fatalln(insertIntoForumUsers, err)		
	}
	//inserting claiming forum id by slug 2
	_, err = db.Prepare("getForumIDBySlug", getForumIDBySlug)
	if err != nil{
		log.Fatalln(getForumIDBySlug, err)		
	}


	// PERF TEST

	//GET FORUM THREADS
	_, err = db.Prepare("gftLimit", gftLimit)
	if err != nil{
		log.Fatalln(gftLimit, err)		
	}
	_, err = db.Prepare("gftLimitDesc", gftLimitDesc)
	if err != nil{
		log.Fatalln(gftLimitDesc, err)		
	}
	_, err = db.Prepare("gftCreatedLimit", gftCreatedLimit)
	if err != nil{
		log.Fatalln(gftCreatedLimit, err)		
	}
	_, err = db.Prepare("gftCreatedLimitDesc", gftCreatedLimitDesc)
	if err != nil{
		log.Fatalln(gftCreatedLimitDesc, err)		
	}

	// GET FORUM USERS
	_, err = db.Prepare("gfuLimit", gfuLimit)
	if err != nil{
		log.Fatalln(gfuLimit, err)		
	}
	_, err = db.Prepare("gfuLimitDesc", gfuLimitDesc)
	if err != nil{
		log.Fatalln(gfuLimitDesc, err)		
	}
	_, err = db.Prepare("gfuSinceLimit", gfuSinceLimit)
	if err != nil{
		log.Fatalln(gfuSinceLimit, err)		
	}
	_, err = db.Prepare("gfuSinceLimitDesc", gfuSinceLimitDesc)
	if err != nil{
		log.Fatalln(gfuSinceLimitDesc, err)		
	}

	//claiming threads
	_, err = db.Prepare("getThreadBySlug", getThreadBySlug)
	if err != nil{
		log.Fatalln(getThreadBySlug, err)		
	}
	_, err = db.Prepare("getThreadById", getThreadById)
	if err != nil{
		log.Fatalln(getThreadById, err)		
	}
	_, err = db.Prepare("checkThreadIdBySlug", checkThreadIdBySlug)
	if err != nil{
		log.Fatalln(checkThreadIdBySlug, err)		
	}
	_, err = db.Prepare("checkThreadIdById", checkThreadIdById)
	if err != nil{
		log.Fatalln(checkThreadIdById, err)		
	}

	//Getting posts
	_, err = db.Prepare("getPostDetailsQuery", getPostDetailsQuery)
	if err != nil{
		log.Fatalln(getPostDetailsQuery, err)		
	}


}

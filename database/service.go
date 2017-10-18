package database

import (
	"github.com/nd-r/tech-db-forum/models"
)

const statusQuery = "SELECT (SELECT count(*) FROM forum) as forum, (SELECT count(*) FROM post) as post, (SELECT count(*) FROM users) as user, (SELECT count(*) FROM thread) as thread"

func GetDBStatus() *models.Status {
	tx := db.MustBegin()
	defer tx.Commit()

	status := models.Status{}
	tx.Get(&status, statusQuery)

	return &status
}

const deleteDBQuery = "TRUNCATE forum_users, vote, post, thread, forum, users RESTART IDENTITY CASCADE"

func DeleteDB(){
	tx := db.MustBegin()
	defer tx.Commit()

	tx.Exec(deleteDBQuery)
}


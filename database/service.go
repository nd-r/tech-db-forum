package database

import (
	"log"

	"github.com/nd-r/tech-db-forum/models"
)

const statusQuery = `SELECT
	(SELECT count(*) FROM forum) as forum,
	(SELECT count(*) FROM post) as post,
	(SELECT count(*) FROM users) as user,
	(SELECT count(*) FROM thread) as thread`

func GetDBStatus() *models.Status {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	status := models.Status{}
	tx.QueryRow(statusQuery).Scan(&status.Forum, &status.Post, &status.User, &status.Thread)

	return &status
}

const deleteDBQuery = `TRUNCATE forum_users,
	vote,
	post,
	thread,
	forum,
	users
RESTART IDENTITY
CASCADE`

func DeleteDB() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	tx.Exec(deleteDBQuery)
}

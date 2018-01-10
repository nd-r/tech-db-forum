package database

import (
	"log"
	"strings"

	"github.com/nd-r/tech-db-forum/models"
)

func GetPostDetails(id *string, related []byte) (*models.PostDetails, int) {

	postDetails := models.PostDetails{}
	postDetails.PostDetails = &models.Post{}

	conn, _ := db.Acquire()
	defer db.Release(conn)

	err := conn.QueryRow("getPostDetailsQuery", id).
		Scan(&postDetails.PostDetails.Id, &postDetails.PostDetails.User_nick,
		&postDetails.PostDetails.Message, &postDetails.PostDetails.Created,
		&postDetails.PostDetails.Forum_slug, &postDetails.PostDetails.Thread_id,
		&postDetails.PostDetails.Is_edited, &postDetails.PostDetails.Parent)
	if err != nil {
		return nil, 404
	}

	if related == nil {
		return &postDetails, 200
	}

	relatedArr := strings.Split(string(related), ",")

	for _, val := range relatedArr {
		switch val {
		case "user":
			postDetails.AuthorDetails = &models.User{}
			conn.QueryRow("getUserProfileQuery", &postDetails.PostDetails.User_nick).
				Scan(&postDetails.AuthorDetails.Nickname, &postDetails.AuthorDetails.Email,
				&postDetails.AuthorDetails.About, &postDetails.AuthorDetails.Fullname)
		case "forum":
			postDetails.ForumDetails = &models.Forum{}
			conn.QueryRow("selectForumQuery", postDetails.PostDetails.Forum_slug).
				Scan(&postDetails.ForumDetails.Slug, &postDetails.ForumDetails.Title,
				&postDetails.ForumDetails.Posts, &postDetails.ForumDetails.Threads,
				&postDetails.ForumDetails.Moderator)
		case "thread":
			postDetails.ThreadDetails = &models.Thread{}
			conn.QueryRow("getThreadById", postDetails.PostDetails.Thread_id).
				Scan(&postDetails.ThreadDetails.Id, &postDetails.ThreadDetails.Slug,
				&postDetails.ThreadDetails.Title, &postDetails.ThreadDetails.Message,
				&postDetails.ThreadDetails.Forum_slug, &postDetails.ThreadDetails.User_nick,
				&postDetails.ThreadDetails.Created, &postDetails.ThreadDetails.Votes_count)
		}
	}
	return &postDetails, 200
}

const updatePostDetailsQuery = `UPDATE post
SET message=coalesce($2,message),
	is_edited=(CASE WHEN $2 IS NULL OR $2 = message
		THEN FALSE
		ELSE TRUE
		END)
WHERE ID=$1
RETURNING id,
	user_nick::TEXT,
	message,
	created,
	forum_slug::TEXT,
	thread_id,
	is_edited,
	parent`

func UpdatePostDetails(id *string, postUpd *models.PostUpdate) (*models.Post, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	postUpdated := models.Post{}

	err = tx.QueryRow(updatePostDetailsQuery, id, postUpd.Message).
		Scan(&postUpdated.Id, &postUpdated.User_nick, &postUpdated.Message,
		&postUpdated.Created, &postUpdated.Forum_slug, &postUpdated.Thread_id,
		&postUpdated.Is_edited, &postUpdated.Parent)
	if err != nil {
		return nil, 404
	}

	return &postUpdated, 200
}

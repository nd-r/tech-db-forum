package database

import (
	"github.com/nd-r/tech-db-forum/models"
	"log"
	"strings"
)

const getPostDetailsQuery = "SELECT id, user_nick, message, created, forum_slug, thread_id, is_edited, parent FROM post WHERE id=$1"
const getPostForumDetailsQuery = "SELECT f.slug, f.posts, f.threads, f.title, f.moderator FROM post p JOIN forum f ON p.forum_slug = f.slug WHERE p.id = $1"
const getPostAuthorDetailsQuery = "SELECT u.nickname, u.fullname, u.email, u.about FROM post p JOIN users u ON p.user_nick = u.nickname WHERE p.id = $1"
const getPostThreadDetailsQuery = "SELECT t.user_nick, t.created,t.forum_slug, t.id, t.message, t.slug, t.title, t.votes_count FROM post p JOIN thread t ON p.thread_id = t.id WHERE p.id = $1"

func GetPostDetails(id *string, related []byte) (*models.PostDetails, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	postDetails := models.PostDetails{}
	postDetails.PostDetails = &models.Post{}

	err = tx.QueryRow(getPostDetailsQuery, id).Scan(&postDetails.PostDetails.Id, &postDetails.PostDetails.User_nick, &postDetails.PostDetails.Message, &postDetails.PostDetails.Created, &postDetails.PostDetails.Forum_slug, &postDetails.PostDetails.Thread_id, &postDetails.PostDetails.Is_edited, &postDetails.PostDetails.Parent)
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
			tx.QueryRow( getPostAuthorDetailsQuery, id).Scan(&postDetails.AuthorDetails.Nickname, &postDetails.AuthorDetails.Fullname,&postDetails.AuthorDetails.Email,&postDetails.AuthorDetails.About)
		case "forum":
			postDetails.ForumDetails = &models.Forum{}
			tx.QueryRow( getPostForumDetailsQuery, id).Scan(&postDetails.ForumDetails.Slug,&postDetails.ForumDetails.Posts,&postDetails.ForumDetails.Threads,&postDetails.ForumDetails.Title,&postDetails.ForumDetails.Moderator)
		case "thread":
			postDetails.ThreadDetails = &models.Thread{}
			tx.QueryRow( getPostThreadDetailsQuery, id).Scan(&postDetails.ThreadDetails.User_nick,&postDetails.ThreadDetails.Created,&postDetails.ThreadDetails.Forum_slug,&postDetails.ThreadDetails.Id,&postDetails.ThreadDetails.Message,&postDetails.ThreadDetails.Slug,&postDetails.ThreadDetails.Title,&postDetails.ThreadDetails.Votes_count)
		}
	}
	return &postDetails, 200
}

const updatePostDetailsQuery = "UPDATE post SET message=coalesce($2,message), is_edited=(CASE WHEN $2 IS NULL OR $2 = message THEN FALSE ELSE TRUE END) WHERE ID=$1 RETURNING id, user_nick, message, created, forum_slug, thread_id, is_edited, parent"

func UpdatePostDetails(id *string, postUpd *models.PostUpdate) (*models.Post, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	postUpdated := models.Post{}

	err = tx.QueryRow(updatePostDetailsQuery, id, postUpd.Message).Scan(&postUpdated.Id, &postUpdated.User_nick, &postUpdated.Message, &postUpdated.Created, &postUpdated.Forum_slug, &postUpdated.Thread_id, &postUpdated.Is_edited, &postUpdated.Parent)
	if err != nil {
		return nil, 404
	}

	return &postUpdated, 200
}

package database

import (
	"github.com/nd-r/tech-db-forum/models"
	"strings"
)

const getPostDetailsQuery = "SELECT * FROM post WHERE id=$1"
const getPostForumDetailsQuery = "SELECT f.slug, f.posts, f.threads, f.title, f.moderator FROM post p JOIN forum f ON p.forum_slug = f.slug WHERE p.id = $1"
const getPostAuthorDetailsQuery = "SELECT u.nickname, u.fullname, u.email, u.about FROM post p JOIN users u ON p.user_nick = u.nickname WHERE p.id = $1"
const getPostThreadDetailsQuery = "SELECT t.user_nick, t.created,t.forum_slug, t.id, t.message, t.slug, t.title, t.votes_count FROM post p JOIN thread t ON p.thread_id = t.id WHERE p.id = $1"

func GetPostDetails(id *string, related []byte) (*models.PostDetails, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	postDetails := models.PostDetails{}
	postDetails.PostDetails = &models.Post{}

	err := tx.Get(postDetails.PostDetails, getPostDetailsQuery, id)
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
			tx.Get(postDetails.AuthorDetails, getPostAuthorDetailsQuery, id)
		case "forum":
			postDetails.ForumDetails = &models.Forum{}
			tx.Get(postDetails.ForumDetails, getPostForumDetailsQuery, id)
		case "thread":
			postDetails.ThreadDetails = &models.Thread{}
			tx.Get(postDetails.ThreadDetails, getPostThreadDetailsQuery, id)
		}
	}
	return &postDetails, 200
}

const updatePostDetailsQuery = "UPDATE post SET message=coalesce($2,message), is_edited=(CASE WHEN $2 IS NULL OR $2 = message THEN FALSE ELSE TRUE END) WHERE ID=$1 RETURNING *"

func UpdatePostDetails(id *string, postUpd *models.PostUpdate) (*models.Post, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	postUpdated := models.Post{}

	err := tx.Get(&postUpdated, updatePostDetailsQuery, id, postUpd.Message)
	if err != nil {
		return nil, 404
	}

	return &postUpdated, 200
}

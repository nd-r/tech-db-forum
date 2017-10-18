package database

import (
	"github.com/lib/pq"
	"github.com/nd-r/tech-db-forum/models"
	"log"
	"strconv"
	"strings"
	"time"
)

const getThreadIdBySlug = "SELECT id FROM thread WHERE lower(slug)=lower($1)"
const getThreadIdById = "SELECT id FROM thread WHERE id=$1"

const insertPost = "INSERT INTO post (id, user_nick, message, created, forum_slug, thread_id, parent, parents)" +
	"VALUES($6::INTEGER, (SELECT nickname FROM users WHERE lower(nickname) = lower($1)), $2, $3, (SELECT forum_slug FROM thread WHERE id=$4), (CASE WHEN $5=0 THEN $4 ELSE CASE WHEN $4=(SELECT thread_id FROM post WHERE id=$5) THEN $4 ELSE NULL END END), $5, ((SELECT parents FROM post WHERE id = $5) || $6::INTEGER))" +
	"RETURNING id, user_nick, message, created, forum_slug,thread_id,is_edited, parent"

const getNextId = "SELECT nextval('post_id_seq')"

const UpdateForumPosts = "UPDATE forum SET posts=posts + $2 WHERE slug=$1"

func CreatePosts(slugOrID string, postsArr *models.PostArr) (*models.PostArr, int) {
	tx := db.MustBegin()
	defer tx.Rollback()
	currTime := time.Now().Format("2006-01-02T15:04:05.000000Z")

	var ID int
	var err error
	if ID, err = strconv.Atoi(slugOrID); err != nil {
		err = tx.Get(&ID, getThreadIdBySlug, slugOrID)
		if err != nil {
			return nil, 404
		}
	} else {
		err = tx.Get(&ID, getThreadIdById, ID)
		if err != nil {
			return nil, 404
		}
	}

	if len(*postsArr) == 0 {
		return nil, 201
	}

	var post models.Post
	postsAdded := models.PostArr{}
	prep, err := tx.Preparex(insertPost)
	if err != nil {
		log.Fatalln(err)
	}

	for _, val := range *postsArr {
		var nextId int
		tx.Get(&nextId, getNextId)
		if err != nil{
			log.Fatalln(err)
		}
		if val.Parent == nil {
			a := 0
			val.Parent = &a
		}
		err = prep.Get(&post, val.User_nick, val.Message, currTime, ID, val.Parent, nextId)
		if err != nil {
			log.Println(err, ID)
			if err.(*pq.Error).Column == "user_nick" {
				return nil, 404
			}
			return nil, 409
		}
		postsAdded = append(postsAdded, post)
	}
	rows, err := tx.Queryx(UpdateForumPosts, post.Forum_slug, len(postsAdded))
	if err != nil {
		log.Fatalln(err)
	}
	rows.Close()
	err = tx.Commit()
	if err != nil{
		log.Fatalln(err)
	}
	return &postsAdded, 201
}

const putVote = "INSERT INTO vote (nickname , thread_id, voice) VALUES ((SELECT nickname FROM users WHERE lower(nickname)=lower($1)), $2, $3)"

const getThreadBySlug = "SELECT * FROM thread WHERE lower(slug)=lower($1)"
const getThreadById = "SELECT * FROM thread WHERE id=$1"
const getPrevVote = "SELECT * FROM vote WHERE lower(nickname)=lower($1) AND thread_id=$2"
const updateVote = "UPDATE TABLE vote SET voice=$1 WHERE id=$2"

func PutVote(slugOrID string, vote *models.Vote) (*models.Thread, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	thread := models.Thread{}

	var ID int
	var err error
	if ID, err = strconv.Atoi(slugOrID); err != nil {
		err = tx.Get(&thread, getThreadBySlug, slugOrID)
		if err != nil {
			tx.Rollback()
			return nil, 404
		}
	} else {
		err = tx.Get(&thread, getThreadById, ID)
		if err != nil {
			tx.Rollback()
			return nil, 404
		}
	}

	prevVote := models.VoteDB{}
	err = tx.Get(&prevVote, getPrevVote, vote.Nickname, thread.Id)
	if err != nil {
		_, err = tx.Queryx(putVote, vote.Nickname, thread.Id, vote.Voice)
		if err != nil {
			tx.Rollback()
			return nil, 404
		}
		*thread.Votes_count += vote.Voice
		return &thread, 200
	}
	if prevVote.Voice != vote.Voice {
		_, err = tx.Queryx(updateVote, vote.Voice, prevVote.ID)
		*thread.Votes_count += vote.Voice * 2
		return &thread, 200
	}
	return &thread, 200
}

func GetThread(slugOrID *string) (*models.Thread, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	thread := models.Thread{}

	var ID int
	var err error

	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		err = tx.Get(&thread, getThreadBySlug, slugOrID)
		if err != nil {
			tx.Rollback()
			return nil, 404
		}
	} else {
		err = tx.Get(&thread, getThreadById, ID)
		if err != nil {
			tx.Rollback()
			return nil, 404
		}
	}

	return &thread, 200
}

const checkThreadId = "SELECT id FROM thread WHERE id=$1"

const getPostsFlat = "SELECT id, user_nick, message, created, forum_slug,thread_id,is_edited, parent" +
	" FROM post WHERE thread_id=$1 AND id >COALESCE($3::INTEGER,0) " +
	" ORDER BY id $4" +
	" LIMIT $2"

const getPostsTree = "SELECT id, user_nick, message, created, forum_slug, thread_id, is_edited, parent FROM post" +
	" WHERE thread_id = $1  AND parents > COALESCE((SELECT parents FROM post WHERE id = $3::INTEGER), '{0}')" +
	" ORDER BY parents $4" +
	" LIMIT $2;"

const getPostsParentTree = "WITH sub AS (" +
	"SELECT parents FROM post" +
	" WHERE parent=0 AND thread_id = $1 AND parents > COALESCE((SELECT parents FROM post WHERE id = $3::INTEGER), '{0}')" +
	" ORDER BY post.parents $4" +
	" LIMIT $2)" +
	" SELECT p.id, p.user_nick, p.forum_slug, p.created, p.message, p.thread_id, p.parent, p.is_edited" +
	" FROM post p" +
	"  JOIN sub ON sub.parents <@ p.parents" +
	" ORDER BY p.parents $4;"

func GetThreadPosts(slugOrID *string, limit []byte, since []byte, sort []byte, desc []byte) (*models.PostArr, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	var ID int
	var err error
	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		err = tx.Get(&ID, getThreadIdBySlug, slugOrID)
		if err != nil {
			return nil, 404
		}
	} else {
		err = tx.Get(&ID, checkThreadId, ID)
		if err != nil {
			return nil, 404
		}
	}

	var query string

	switch string(sort) {
	case "tree":
		query = getPostsTree
	case "parent_tree":
		query = getPostsParentTree
	default:
		query = getPostsFlat
	}

	if string(desc) == "true" {
		query = strings.Replace(query, "$4", " DESC", -1)
		if since != nil {
			query = strings.Replace(query, ">", " <", -1)
		}
	} else {
		query = strings.Replace(query, "$4", " ASC", -1)
	}
	var posts models.PostArr

	if limit != nil {
		if since != nil {
			err = tx.Select(&posts, query, ID, limit, since)
		} else {
			err = tx.Select(&posts, query, ID, limit, nil)
		}
	} else {
		if since != nil {
			err = tx.Select(&posts, query, ID, nil, since)
		} else {
			err = tx.Select(&posts, query, ID, nil, nil)
		}
	}

	if err != nil {
		log.Println(err)
	}

	return &posts, 200
}

const threadUpdateQuery = "UPDATE thread SET message = coalesce($1, message), title = coalesce($2,title) WHERE id = $3 RETURNING *"

func UpdateThreadDetails(slugOrID *string, thrUpdate *models.ThreadUpdate) (*models.Thread, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	var ID int
	var err error
	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		err = tx.Get(&ID, getThreadIdBySlug, slugOrID)
		if err != nil {
			return nil, 404
		}
	}

	var thread models.Thread

	err = tx.Get(&thread, threadUpdateQuery, thrUpdate.Message, thrUpdate.Title, ID)

	if err != nil {
		return nil, 404
	}
	return &thread, 200
}

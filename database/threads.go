package database

import (
	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"log"
	"strconv"
	"strings"
	"time"
)

const getThreadIdAndForumSlugBySlug = "SELECT id, forum_slug::TEXT FROM thread WHERE slug=$1"
const getThreadIdAndForumSlugById = "SELECT id, forum_slug::TEXT FROM thread WHERE id=$1"

const insertPost = "INSERT INTO post (id, user_nick, message, forum_slug, thread_id, parent, parents, created) VALUES ($1 :: INTEGER, $2, $3, $4, $5, $6, $8::BIGINT[] || $1 :: BIGINT, $7)"

const UpdateForumPosts = "UPDATE forum SET posts=posts+$2 WHERE slug=$1"

const generateNextIDs = "SELECT array_agg(nextval('post_id_seq')::BIGINT) FROM generate_series(1,$1);"

const selectParentAndParents = "SELECT thread_id, parents FROM post WHERE id = $1"

var created, _ = time.Parse("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:05.010000Z")

const claimUserInfoAndInsertIntoFU = "WITH insertion AS (INSERT INTO forum_users (SELECT $1, nickname, email, about, fullname FROM users WHERE nickname=$2) ON CONFLICT DO NOTHING) " +
	"SELECT nickname::TEXT, email::TEXT, about, fullname FROM users WHERE nickname = $2"

func CreatePosts(slugOrID interface{}, postsArr *models.PostArr) (*models.PostArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback()

	var forumID, threadID int
	var forumSlug string

	//Claiming thread ID
	threadID, err = strconv.Atoi(slugOrID.(string))
	if err != nil {
		if err = tx.QueryRow(getThreadIdAndForumSlugBySlug, slugOrID).Scan(&threadID, &forumSlug); err != nil {
			log.Println(err)
			return nil, dberrors.ErrThreadNotFound
		}
	} else {
		if err = tx.QueryRow(getThreadIdAndForumSlugById, threadID).Scan(&threadID, &forumSlug); err != nil {
			log.Println(err)
			return nil, dberrors.ErrThreadNotFound
		}
	}

	if len(*postsArr) == 0 {
		return nil, nil
	}

	//claiming forum id
	if err = tx.QueryRow(getForumIDBySlug, &forumSlug).Scan(&forumID); err != nil {
		log.Fatalln(err)
	}

	//claiming ids for further posts
	ids := make([]int64, 0, len(*postsArr))
	if err = tx.QueryRow("generateNextIDs", len(*postsArr)).Scan(&ids); err != nil {
		log.Fatalln(err)
	}

	//Inserting posts
	rowsToCopy := [][]interface{}{}
	for index, post := range *postsArr {
		var parentThreadID int64

		if post.Parent != 0 {
			if err = tx.QueryRow(selectParentAndParents, post.Parent).
				Scan(&parentThreadID, &post.Parents); err != nil {
				log.Println(err, post.Parent)
				return nil, dberrors.ErrPostsConflict
			}
			if parentThreadID != 0 && parentThreadID != int64(threadID) {
				log.Println(err)
				return nil, dberrors.ErrPostsConflict
			}
		}

		user := models.User{}
		if err = tx.QueryRow(claimUserInfoAndInsertIntoFU, forumID, post.User_nick).Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
			log.Println(err)
			return nil, dberrors.ErrUserNotFound
		}

		post.Id = int(ids[index])
		post.Thread_id = threadID
		post.Forum_slug = forumSlug
		post.Created = &created
		post.User_nick = user.Nickname
		post.Parents = append(post.Parents, ids[index])
		// _, err = tx.Exec("insertIntoForumUsers", forumID, &user.Nickname, &user.Email, &user.About, &user.Fullname)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		rowsToCopy = append(rowsToCopy, []interface{}{post.Id, post.User_nick, post.Message, post.Created, post.Forum_slug, post.Thread_id, post.Parent, post.Parents})
	}

	rowsCreated, err := tx.CopyFrom(pgx.Identifier{"post"}, []string{"id", "user_nick", "message", "created", "forum_slug", "thread_id", "parent", "parents"}, pgx.CopyFromRows(rowsToCopy))
	if err != nil {
		log.Fatalln(err)
	}
	if rowsCreated != len(*postsArr) {
		log.Println(err, rowsCreated)
	}

	_, err = tx.Exec(UpdateForumPosts, forumSlug, len(*postsArr))
	if err != nil {
		log.Fatalln(err)
	}

	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return postsArr, nil
}

const putVoteByThrID = "WITH sub AS (INSERT INTO vote (user_id, thread_id, voice) VALUES ((SELECT id FROM users WHERE nickname=$1), $2, $3) ON CONFLICT ON CONSTRAINT unique_user_and_thread DO UPDATE SET prev_voice = vote.voice ,voice = EXCLUDED.voice RETURNING prev_voice, voice, thread_id) UPDATE thread SET votes_count = votes_count - (SELECT prev_voice-voice FROM sub) WHERE id = $2 RETURNING id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count;"
const putVoteByThrSLUG = "WITH sub AS (INSERT INTO vote (user_id, thread_id, voice) VALUES ((SELECT id FROM users WHERE nickname=$1), (SELECT id FROM thread WHERE slug=$2), $3) ON CONFLICT ON CONSTRAINT unique_user_and_thread DO UPDATE SET prev_voice = vote.voice ,voice = EXCLUDED.voice RETURNING prev_voice, voice, thread_id) UPDATE thread SET votes_count = votes_count - (SELECT prev_voice-voice FROM sub) WHERE slug=$2 RETURNING id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count;"

func PutVote(slugOrID interface{}, vote *models.Vote) (*models.Thread, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	_, err = strconv.Atoi(slugOrID.(string))

	thread := models.Thread{}

	if err != nil {
		err = tx.QueryRow("putVoteByThrSLUG", vote.Nickname, slugOrID, vote.Voice).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	} else {
		err = tx.QueryRow("putVoteByThrID", vote.Nickname, slugOrID, vote.Voice).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	}

	if err != nil {
		tx.Rollback()
		log.Println(err)
		return nil, err
	}
	return &thread, nil
}

const getThreadById = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread WHERE id=$1"

const getThreadBySlug = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE slug=$1"

func GetThread(slugOrID interface{}) (*models.Thread, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	thread := models.Thread{}

	_, err = strconv.Atoi(slugOrID.(string))

	if err != nil {
		err = tx.QueryRow(getThreadBySlug, slugOrID).
			Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
		log.Println(err)
		return &thread, err
	}

	err = tx.QueryRow(getThreadById, slugOrID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	log.Println(err)
	return &thread, err
}

const checkThreadId = "SELECT id FROM thread WHERE id=$1"

const getPostsFlat = "SELECT id, user_nick::TEXT, message, created, forum_slug::TEXT,thread_id,is_edited, parent" +
	" FROM post WHERE thread_id=$1 AND id >COALESCE($3::TEXT::INTEGER,0) " +
	" ORDER BY id $4" +
	" LIMIT $2::TEXT::BIGINT"

const getPostsTree = "SELECT id, user_nick::TEXT, message, created, forum_slug::TEXT, thread_id, is_edited, parent FROM post" +
	" WHERE thread_id = $1  AND parents > COALESCE((SELECT parents FROM post WHERE id = $3::TEXT::INTEGER), '{0}')" +
	" ORDER BY parents $4" +
	" LIMIT $2::TEXT::BIGINT;"

const getPostsParentTree = "WITH sub AS (" +
	"SELECT parents FROM post" +
	" WHERE parent=0 AND thread_id = $1 AND parents > COALESCE((SELECT parents FROM post WHERE id = $3::TEXT::INTEGER), '{0}')" +
	" ORDER BY post.parents $4" +
	" LIMIT $2::TEXT::INTEGER)" +
	" SELECT p.id, p.user_nick::TEXT, p.message, p.created, p.forum_slug::TEXT, p.thread_id, p.is_edited, p.parent" +
	" FROM post p" +
	"  JOIN sub ON sub.parents <@ p.parents" +
	" ORDER BY p.parents $4;"

func GetThreadPosts(slugOrID *string, limit []byte, since []byte, sort []byte, desc []byte) (*models.PostArr, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	var fs string
	if _, err = strconv.Atoi(*slugOrID); err != nil {
		if err = tx.QueryRow(getThreadIdAndForumSlugBySlug, slugOrID).Scan(&ID, &fs); err != nil {

			return nil, 404
		}
	} else {
		if err = tx.QueryRow(checkThreadId, slugOrID).Scan(&ID); err != nil {
			log.Println(err)
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
	var rows *pgx.Rows
	if limit != nil {
		if since != nil {
			rows, err = tx.Query(query, ID, limit, since)
		} else {
			rows, err = tx.Query(query, ID, limit, nil)
		}
	} else {
		if since != nil {
			rows, err = tx.Query(query, ID, nil, since)
		} else {
			rows, err = tx.Query(query, ID, nil, nil)
		}
	}

	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(&post.Id, &post.User_nick, &post.Message, &post.Created, &post.Forum_slug, &post.Thread_id, &post.Is_edited, &post.Parent); err != nil {
			log.Fatalln(err)
		}
		posts = append(posts, &post)
	}

	if err != nil {
		log.Println(err)
	}

	return &posts, 200
}

const threadUpdateQuery = "UPDATE thread SET message = coalesce($1, message), title = coalesce($2,title) WHERE id = $3 RETURNING  id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count "

func UpdateThreadDetails(slugOrID *string, thrUpdate *models.ThreadUpdate) (*models.Thread, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	var fs string
	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		err = tx.QueryRow(getThreadIdAndForumSlugBySlug, slugOrID).Scan(&ID, &fs)
		if err != nil {
			log.Println(err)
			return nil, 404
		}
	}

	var thread models.Thread

	err = tx.QueryRow(threadUpdateQuery, thrUpdate.Message, thrUpdate.Title, ID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	log.Println(err)
	if err != nil {
		return nil, 404
	}
	return &thread, 200
}

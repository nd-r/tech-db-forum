package database

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"time"
	// 	"github.com/lib/pq"
	"github.com/nd-r/tech-db-forum/models"
	"log"
	"strconv"
	"strings"
	// 	"time"
)

const getThreadIdBySlug = "SELECT id, forum_slug FROM thread WHERE lower(slug)=lower($1)"
const getThreadIdById = "SELECT id, forum_slug FROM thread WHERE id=$1"

const insertPost = "INSERT INTO post (id, user_nick, message, forum_slug, thread_id, parent, parents, created) VALUES ($1 :: INTEGER, $2, $3, $4, (CASE WHEN $6 != 0 THEN CASE WHEN (SELECT thread_id FROM post WHERE id = $6) = $5 THEN $5 ELSE NULL END ELSE $5 END), $6, (SELECT parents FROM post WHERE id = $6) || $1 :: BIGINT, $7)"

// const insertPost = "INSERT INTO post (id, user_nick, message, forum_slug, thread_id, parent, parents, created) VALUES ($1::INTEGER, $2, $3, $4, $5, $6, $7::BIGINT[] || $1::BIGINT, $8)"

const insertForumUsers = "INSERT INTO forum_users VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING"

const UpdateForumPosts = "UPDATE forum SET posts=posts + $2 WHERE lower(slug)=lower($1)"

const generateNextIDs = "SELECT array_agg(nextval('post_id_seq')::BIGINT) FROM generate_series(1,$1);"

const claimInfoWithParent = "SELECT users.nickname, users.about, users.email, users.fullname, post.parents, post.thread_id FROM users, post WHERE lower(users.nickname) = lower($1) AND post.id = $2"
const claimInfoWithoutParent = "SELECT users.nickname, users.about, users.email, users.fullname, NULL, NULL FROM users WHERE lower(users.nickname) = lower($1)"
const selectParentAndParents = "SELECT thread_id, parents FROM post WHERE id = $1"

var zero = 0
const getForumIdBySlug = "SELECT id FROM forum WHERE lower(slug) = lower($1)"

func CreatePosts(slugOrID interface{}, postsArr *models.PostArr) (*models.PostArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var forumID int
	var forumSlug string

	created := time.Now()

	//Claiming thread ID
	threadID, err := strconv.Atoi(slugOrID.(string))
	if err != nil {
		if err = tx.QueryRow("getThreadIdBySlug", slugOrID).Scan(&threadID, &forumSlug); err != nil {
			return nil, dberrors.ErrThreadNotFound
		}
	} else {
		if err = tx.QueryRow("getThreadIdById", threadID).Scan(&threadID, &forumSlug); err != nil {
			return nil, dberrors.ErrThreadNotFound
		}
	}

	if len(*postsArr) == 0 {
		return nil, nil
	}

	//claiming forum id
	if err = tx.QueryRow("getForumIdBySlug", &forumSlug).Scan(&forumID); err != nil {
		log.Fatalln(err)
	}

	//claiming ids for further posts
	ids := make([]int64, 0, len(*postsArr))
	if err = tx.QueryRow("generateNextIDs", len(*postsArr)).Scan(&ids); err != nil {
		log.Fatalln(err)
	}

	user := models.User{}
	if err = tx.QueryRow("getUserProfileQuery", (*postsArr)[0].User_nick).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
		return nil, dberrors.ErrUserNotFound
	}

	tx.Exec("insertForumUsers", forumID, user.Nickname, user.About, user.Email, user.Fullname)

	batchPostInserter := tx.BeginBatch()
	// batchForumUsersInserter := tx.BeginBatch()
	for index, post := range *postsArr {
		if post.Parent == nil {
			post.Parent = &zero
		}

		// var parentThreadID int64
		// var parents []int64

		if post.Parent != nil && *post.Parent != 0 {
			// if err = tx.QueryRow("selectParentAndParents", post.Parent).Scan(&parentThreadID, &parents); err != nil {
			// 	return nil, dberrors.ErrPostsConflict
			// }

			// if parentThreadID != 0 && parentThreadID != int64(threadID) {
			// 	return nil, dberrors.ErrPostsConflict
			// }
		} else {
			post.Parent = &zero
		}

		pID := int(ids[index])

		post.Id = &pID
		post.Thread_id = &threadID
		post.Forum_slug = forumSlug
		post.Created = &created

		batchPostInserter.Queue("insertPost", []interface{}{ids[index], &user.Nickname, &post.Message, &forumSlug, &threadID, post.Parent, created}, nil, nil)
	}

	if err = batchPostInserter.Send(context.Background(), nil); err != nil {
		log.Fatalln(err)
	}

	res, err := batchPostInserter.ExecResults()
	if int(res.RowsAffected()) != 1 {
		return nil, dberrors.ErrPostsConflict
	}
	batchPostInserter.Close()

	tx.Exec(UpdateForumPosts, forumSlug, len(*postsArr))
	return postsArr, nil
}

const putVoteByThrID = "WITH sub AS (INSERT INTO vote (user_id, thread_id, voice) VALUES ((SELECT id FROM users WHERE lower(nickname) = lower($1)), $2, $3) ON CONFLICT ON CONSTRAINT unique_user_and_thread DO UPDATE SET prev_voice = vote.voice ,voice = EXCLUDED.voice RETURNING prev_voice, voice, thread_id) UPDATE thread SET votes_count = votes_count - (SELECT prev_voice-voice FROM sub) WHERE id = $2 RETURNING id, slug, title, message, forum_slug, user_nick, created, votes_count;"
const putVoteByThrSLUG = "WITH sub AS (INSERT INTO vote (user_id, thread_id, voice) VALUES ((SELECT id FROM users WHERE lower(nickname) = lower($1)), (SELECT id FROM thread WHERE lower(slug) = lower($2)), $3) ON CONFLICT ON CONSTRAINT unique_user_and_thread DO UPDATE SET prev_voice = vote.voice ,voice = EXCLUDED.voice RETURNING prev_voice, voice, thread_id) UPDATE thread SET votes_count = votes_count - (SELECT prev_voice-voice FROM sub) WHERE lower(slug) = lower($2) RETURNING id, slug, title, message, forum_slug, user_nick, created, votes_count;"

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
		return nil, err
	}

	return &thread, nil
}

const getThreadById = "SELECT id, slug, title, message, forum_slug, user_nick, created, votes_count FROM thread WHERE id=$1"
const getThreadBySlug = "SELECT id, slug, title, message, forum_slug, user_nick, created, votes_count FROM thread WHERE lower(slug)=lower($1)"

func GetThread(slugOrID interface{}) (*models.Thread, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	thread := models.Thread{}

	_, err = strconv.Atoi(slugOrID.(string))

	if err != nil {
		err = tx.QueryRow(getThreadBySlug, slugOrID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
		return &thread, err
	}

	err = tx.QueryRow(getThreadById, slugOrID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	return &thread, err
}

const checkThreadId = "SELECT id FROM thread WHERE id=$1"

const getPostsFlat = "SELECT id, user_nick, message, created, forum_slug,thread_id,is_edited, parent" +
	" FROM post WHERE thread_id=$1 AND id >COALESCE($3::TEXT::INTEGER,0) " +
	" ORDER BY id $4" +
	" LIMIT $2::TEXT::BIGINT"

const getPostsTree = "SELECT id, user_nick, message, created, forum_slug, thread_id, is_edited, parent FROM post" +
	" WHERE thread_id = $1  AND parents > COALESCE((SELECT parents FROM post WHERE id = $3::TEXT::INTEGER), '{0}')" +
	" ORDER BY parents $4" +
	" LIMIT $2::TEXT::BIGINT;"

const getPostsParentTree = "WITH sub AS (" +
	"SELECT parents FROM post" +
	" WHERE parent=0 AND thread_id = $1 AND parents > COALESCE((SELECT parents FROM post WHERE id = $3::TEXT::INTEGER), '{0}')" +
	" ORDER BY post.parents $4" +
	" LIMIT $2::TEXT::INTEGER)" +
	" SELECT p.id, p.user_nick, p.message, p.created, p.forum_slug, p.thread_id, p.is_edited, p.parent" +
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
		if err = tx.QueryRow(getThreadIdBySlug, slugOrID).Scan(&ID, &fs); err != nil {

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

const threadUpdateQuery = "UPDATE thread SET message = coalesce($1, message), title = coalesce($2,title) WHERE id = $3 RETURNING  id, slug, title, message, forum_slug, user_nick, created, votes_count "

func UpdateThreadDetails(slugOrID *string, thrUpdate *models.ThreadUpdate) (*models.Thread, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	var fs string
	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		err = tx.QueryRow(getThreadIdBySlug, slugOrID).Scan(&ID, &fs)
		if err != nil {
			return nil, 404
		}
	}

	var thread models.Thread

	err = tx.QueryRow(threadUpdateQuery, thrUpdate.Message, thrUpdate.Title, ID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)

	if err != nil {
		return nil, 404
	}
	return &thread, 200
}

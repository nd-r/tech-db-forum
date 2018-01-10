package database

import (
	"bytes"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
)

var created, _ = time.Parse("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:05.010000Z")

type forumUser struct {
	userNickname *string
	userEmail    *string
	userAbout    *string
	userFullname *string
}
type forumUserArr []forumUser

func (fu forumUserArr) Len() int {
	return len(fu)
}
func (fu forumUserArr) Swap(i, j int) {
	fu[i], fu[j] = fu[j], fu[i]
}
func (fu forumUserArr) Less(i, j int) bool {
	return *(fu[i].userNickname) < *(fu[j].userNickname)
}

const generateNextIDs = `SELECT
	array_agg(nextval('post_id_seq')::BIGINT)
FROM generate_series(1,$1)`

const selectParentAndParents = `SELECT thread_id,
	parents
FROM post
WHERE id = $1`

const getThreadIdAndForumSlugBySlug = `SELECT id,
	forum_slug::TEXT
FROM thread
WHERE slug=$1`

const getThreadIdAndForumSlugById = `SELECT id,
	forum_slug::TEXT
FROM thread
WHERE id=$1`

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
	if err = tx.QueryRow("getForumIDBySlug", &forumSlug).Scan(&forumID); err != nil {
		log.Fatalln(err)
	}

	//claiming ids for further posts
	ids := make([]int64, 0, len(*postsArr))
	if err = tx.QueryRow(generateNextIDs, len(*postsArr)).Scan(&ids); err != nil {
		log.Fatalln(err)
	}


	var allFu forumUserArr
	//Inserting posts
	var rowsToCopy [][]interface{}
	for index, post := range *postsArr {
		var parentThreadID int64

		if post.Parent != 0 {
			if err = tx.QueryRow(selectParentAndParents, post.Parent).
				Scan(&parentThreadID, &post.Parents); err != nil {
				return nil, dberrors.ErrPostsConflict
			}
			if parentThreadID != 0 && parentThreadID != int64(threadID) {
				return nil, dberrors.ErrPostsConflict
			}
		}

		user := models.User{}
		if err = tx.QueryRow("getUserProfileQuery", post.User_nick).Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
			log.Println(err)
			return nil, dberrors.ErrUserNotFound
		}

		post.Id = int(ids[index])
		post.Thread_id = threadID
		post.Forum_slug = forumSlug
		post.Created = &created
		post.User_nick = user.Nickname
		post.Parents = append(post.Parents, ids[index])

		allFu = append(allFu, forumUser{&user.Nickname, &user.Email, &user.About, &user.Fullname})

		rowsToCopy = append(rowsToCopy, []interface{}{post.Id, post.User_nick, post.Message, post.Created, post.Forum_slug, post.Thread_id, post.Parent, post.Parents, post.Parents[0]})
	}

	sort.Sort(allFu)

	for _, i := range allFu {
		_, err = tx.Exec("insertIntoForumUsers", forumID, i.userNickname, i.userEmail, i.userAbout, i.userFullname)
		if err != nil {
			log.Fatalln(err)
		}
	}

	rowsCreated, err := tx.CopyFrom(pgx.Identifier{"post"}, []string{"id", "user_nick", "message", "created", "forum_slug", "thread_id", "parent", "parents", "main_parent"}, pgx.CopyFromRows(rowsToCopy))
	if err != nil {
		log.Fatalln(err)
	}
	if rowsCreated != len(*postsArr) {
		log.Println(err, rowsCreated)
	}

	_, err = tx.Exec(`UPDATE forum SET posts=posts+$2 WHERE slug=$1`, forumSlug, len(*postsArr))
	if err != nil {
		log.Fatalln(err)
	}

	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	return postsArr, nil
}

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

func GetThread(slugOrID interface{}) (*models.Thread, error) {
	thread := models.Thread{}

	conn, _ := db.Acquire()
	defer db.Release(conn)

	_, err := strconv.Atoi(slugOrID.(string))

	if err != nil {
		err = conn.QueryRow("getThreadBySlug", slugOrID).
			Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
		return &thread, err
	}

	err = conn.QueryRow("getThreadById", slugOrID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	return &thread, err
}

const threadUpdateQuery = `UPDATE thread
SET message = coalesce($1, message),
	title = coalesce($2,title)
WHERE id = $3
RETURNING  id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count `

func UpdateThreadDetails(slugOrID *string, thrUpdate *models.ThreadUpdate) (*models.Thread, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	var fs string
	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		if err = tx.QueryRow(getThreadIdAndForumSlugBySlug, slugOrID).Scan(&ID, &fs);
			err != nil {
			return nil, 404
		}
	}

	var thread models.Thread

	if err = tx.QueryRow(threadUpdateQuery, thrUpdate.Message, thrUpdate.Title, ID).
		Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug,
		&thread.User_nick, &thread.Created, &thread.Votes_count);
		err != nil {
		return nil, 404
	}
	return &thread, 200
}

func getThreadPostsFlat(conn *pgx.Conn, ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	defer db.Release(conn)

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsFlatSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = conn.Query("getPostsFlatSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsFlatSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = conn.Query("getPostsFlatSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsFlatLimitDesc", ID, limit)
			} else {
				rows, err = conn.Query("getPostsFlatLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsFlatLimitDesc", ID, nil)
			} else {
				rows, err = conn.Query("getPostsFlatLimit", ID, nil)
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var posts models.PostArr

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(&post.Id, &post.User_nick, &post.Message,
			&post.Created, &post.Forum_slug, &post.Thread_id,
			&post.Is_edited, &post.Parent);
			err != nil {
			log.Fatalln(err)
		}
		posts = append(posts, &post)
	}
	rows.Close()

	return &posts, 200
}

func getThreadPostsTree(conn *pgx.Conn, ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	defer db.Release(conn)

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = conn.Query("getPostsTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = conn.Query("getPostsTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsTreeLimitDesc", ID, limit)
			} else {
				rows, err = conn.Query("getPostsTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsTreeLimitDesc", ID, nil)
			} else {
				rows, err = conn.Query("getPostsTreeLimit", ID, nil)
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var posts models.PostArr

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(&post.Id, &post.User_nick, &post.Message,
			&post.Created, &post.Forum_slug, &post.Thread_id,
			&post.Is_edited, &post.Parent);
			err != nil {
			log.Fatalln(err)
		}
		posts = append(posts, &post)
	}
	rows.Close()

	return &posts, 200
}

func getThreadPostsParentTree(conn *pgx.Conn, ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	defer db.Release(conn)

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsParentTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = conn.Query("getPostsParentTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsParentTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = conn.Query("getPostsParentTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsParentTreeLimitDesc", ID, limit)
			} else {
				rows, err = conn.Query("getPostsParentTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = conn.Query("getPostsParentTreeLimitDesc", ID, nil)
			} else {
				rows, err = conn.Query("getPostsParentTreeLimit", ID, nil)
			}
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	var posts models.PostArr

	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(&post.Id, &post.User_nick, &post.Message,
			&post.Created, &post.Forum_slug, &post.Thread_id,
			&post.Is_edited, &post.Parent);
			err != nil {
			log.Fatalln(err)
		}
		posts = append(posts, &post)
	}
	rows.Close()

	return &posts, 200
}

func GetThreadPosts(slugOrID *string, limit []byte, since []byte, sort []byte, desc []byte) (*models.PostArr, int) {
	var ID int
	var err error

	conn, _ := db.Acquire()

	if _, err = strconv.Atoi(*slugOrID); err != nil {
		if err = db.QueryRow("checkThreadIdBySlug", slugOrID).Scan(&ID); err != nil {
			db.Release(conn)
			return nil, 404
		}
	} else {
		if err = db.QueryRow("checkThreadIdById", slugOrID).Scan(&ID); err != nil {
			db.Release(conn)
			return nil, 404
		}
	}

	switch true {
	case bytes.Equal([]byte("tree"), sort):
		return getThreadPostsTree(conn, ID, limit, since, desc)
	case bytes.Equal([]byte("parent_tree"), sort):
		return getThreadPostsParentTree(conn, ID, limit, since, desc)
	default:
		return getThreadPostsFlat(conn, ID, limit, since, desc)
	}
}

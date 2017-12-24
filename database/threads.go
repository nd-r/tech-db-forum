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
func (fu forumUserArr) Less(i, j int) bool{
	return *(fu[i].userNickname) < *(fu[j].userNickname)
}

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
		if err = tx.QueryRow("getThreadIdAndForumSlugBySlug", slugOrID).Scan(&threadID, &forumSlug); err != nil {
			log.Println(err)
			return nil, dberrors.ErrThreadNotFound
		}
	} else {
		if err = tx.QueryRow("getThreadIdAndForumSlugById", threadID).Scan(&threadID, &forumSlug); err != nil {
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
	if err = tx.QueryRow("generateNextIDs", len(*postsArr)).Scan(&ids); err != nil {
		log.Fatalln(err)
	}

	var allFu forumUserArr
	//Inserting posts
	var rowsToCopy [][]interface{}
	for index, post := range *postsArr {
		var parentThreadID int64

		if post.Parent != 0 {
			if err = tx.QueryRow("selectParentAndParents", post.Parent).
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

	_, err = tx.Exec("updateForumPosts", forumSlug, len(*postsArr))
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
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	thread := models.Thread{}

	_, err = strconv.Atoi(slugOrID.(string))

	if err != nil {
		err = tx.QueryRow("getThreadBySlug", slugOrID).
			Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
		return &thread, err
	}

	err = tx.QueryRow("getThreadById", slugOrID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
	return &thread, err
}

// const getPostsParentTree = "WITH sub AS (" +
// 	"SELECT parents FROM post" +
// 	" WHERE parent=0 AND thread_id = $1 AND parents > COALESCE((SELECT parents FROM post WHERE id = $3::TEXT::INTEGER), '{0}')" +
// 	" ORDER BY post.parents $4" +
// 	" LIMIT $2::TEXT::INTEGER)" +
// 	" SELECT p.id, p.user_nick::TEXT, p.message, p.created, p.forum_slug::TEXT, p.thread_id, p.is_edited, p.parent" +
// 	" FROM post p" +
// 	"  JOIN sub ON sub.parents <@ p.parents" +
// 	" ORDER BY p.parents $4;"

func UpdateThreadDetails(slugOrID *string, thrUpdate *models.ThreadUpdate) (*models.Thread, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	var fs string
	if ID, err = strconv.Atoi(*slugOrID); err != nil {
		if err = tx.QueryRow("getThreadIdAndForumSlugBySlug", slugOrID).Scan(&ID, &fs);
			err != nil {
			return nil, 404
		}
	}

	var thread models.Thread

	if err = tx.QueryRow("threadUpdateQuery", thrUpdate.Message, thrUpdate.Title, ID).
		Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug,
		&thread.User_nick, &thread.Created, &thread.Votes_count);
		err != nil {
		return nil, 404
	}
	return &thread, 200
}

func getThreadPostsFlat(ID int, limit []byte, since []byte, desc []byte, tx *pgx.Tx) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsFlatSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = tx.Query("getPostsFlatSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsFlatSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = tx.Query("getPostsFlatSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsFlatLimitDesc", ID, limit)
			} else {
				rows, err = tx.Query("getPostsFlatLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsFlatLimitDesc", ID, nil)
			} else {
				rows, err = tx.Query("getPostsFlatLimit", ID, nil)
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

func getThreadPostsTree(ID int, limit []byte, since []byte, desc []byte, tx *pgx.Tx) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = tx.Query("getPostsTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = tx.Query("getPostsTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsTreeLimitDesc", ID, limit)
			} else {
				rows, err = tx.Query("getPostsTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsTreeLimitDesc", ID, nil)
			} else {
				rows, err = tx.Query("getPostsTreeLimit", ID, nil)
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

func getThreadPostsParentTree(ID int, limit []byte, since []byte, desc []byte, tx *pgx.Tx) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsParentTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = tx.Query("getPostsParentTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsParentTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = tx.Query("getPostsParentTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsParentTreeLimitDesc", ID, limit)
			} else {
				rows, err = tx.Query("getPostsParentTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = tx.Query("getPostsParentTreeLimitDesc", ID, nil)
			} else {
				rows, err = tx.Query("getPostsParentTreeLimit", ID, nil)
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
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	if _, err = strconv.Atoi(*slugOrID); err != nil {
		if err = tx.QueryRow("checkThreadIdBySlug", slugOrID).Scan(&ID); err != nil {
			return nil, 404
		}
	} else {
		if err = tx.QueryRow("checkThreadIdById", slugOrID).Scan(&ID); err != nil {
			return nil, 404
		}
	}

	switch true {
	case bytes.Equal([]byte("tree"), sort):
		return getThreadPostsTree(ID, limit, since, desc, tx)
	case bytes.Equal([]byte("parent_tree"), sort):
		return getThreadPostsParentTree(ID, limit, since, desc, tx)
	default:
		return getThreadPostsFlat(ID, limit, since, desc, tx)
	}
}

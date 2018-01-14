package database

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/emirpasic/gods/sets/treeset"
	"strings"
	"context"
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

const getThreadIdAndForumSlugBySlug = `SELECT id,
	forum_slug::TEXT
FROM thread
WHERE slug=$1`

const getThreadIdAndForumSlugById = `SELECT id,
	forum_slug::TEXT
FROM thread
WHERE id=$1`

func usersComparator(a, b interface{}) int {
	u1 := a.(*models.User)
	u2 := b.(*models.User)

	return u1.Compare(u2)
}

func StringsCompare(a, b interface{}) int {
	return strings.Compare(a.(string), b.(string))
}

func CreatePosts(slugOrID interface{}, postsArr *models.PostArr) (*models.PostArr, error) {
	tx := TxMustBegin()
	defer tx.Rollback()

	batch := tx.BeginBatch()
	defer batch.Close()

	var err error
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

	/**
	Check all parents for existence and same thread id
	 */
	var postsWaitingParents []int
	userNicknameSet := treeset.NewWith(StringsCompare)

	for i, post := range *postsArr {
		userNicknameSet.Add(strings.ToLower(post.User_nick))

		if post.Parent != 0 {
			postsWaitingParents = append(postsWaitingParents, i)
			batch.Queue("selectParentAndParents", []interface{}{int(post.Parent)}, nil, nil)
		}
	}

	userNicknameOrderedSet := userNicknameSet.Values()

	for _, userNickname := range userNicknameOrderedSet {
		batch.Queue("getUserProfileQuery", []interface{}{userNickname}, nil, nil)
	}

	var parentThreadID int64
	if err = batch.Send(context.Background(), nil);
		err != nil {
		log.Fatalln(err)
	}

	for _, postIdx := range postsWaitingParents {
		if err = batch.QueryRowResults().
			Scan(&parentThreadID, &(*postsArr)[postIdx].Parents);
			err != nil {
			return nil, dberrors.ErrPostsConflict
		}
		if parentThreadID != 0 && parentThreadID != int64(threadID) {
			return nil, dberrors.ErrPostsConflict
		}
	}

	userRealNicknameMap := make(map[string]string)
	var userModelsOrderedSet models.UsersArr

	for _, userNickname := range userNicknameOrderedSet {
		user := models.User{}
		if err = batch.QueryRowResults().
			Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname);
			err != nil {
			return nil, dberrors.ErrUserNotFound
		}
		userModelsOrderedSet = append(userModelsOrderedSet, &user)
		userRealNicknameMap[userNickname.(string)] = user.Nickname
	}

	/**
	end check
	 */

	for index, post := range *postsArr {
		post.Id = int(ids[index])
		post.Thread_id = threadID
		post.Forum_slug = forumSlug
		post.Created = &created
		post.User_nick = userRealNicknameMap[strings.ToLower(post.User_nick)]
		post.Parents = append(post.Parents, int32(ids[index]))

		batch.Queue("insertIntoPost", []interface{}{post.Id, post.User_nick, post.Message, post.Created, post.Forum_slug, post.Thread_id, post.Parent, post.Parents, post.Parents[0]}, nil, nil)
	}

	for _, user := range userModelsOrderedSet {
		batch.Queue("insertIntoForumUsers", []interface{}{forumID, user.Nickname, user.Email, user.About, user.Fullname}, nil, nil)
	}

	if err = batch.Send(context.Background(), nil);
		err != nil {
		log.Fatalln(err)
	}

	for range *postsArr {
		if _, err := batch.ExecResults(); err != nil {
			log.Fatalln(err)
		}
	}

	for range userModelsOrderedSet {
		if _, err := batch.ExecResults(); err != nil {
			log.Fatalln(err)
		}
	}

	_, err = tx.Exec(`UPDATE forum SET posts=posts+$2 WHERE slug=$1`, forumSlug, len(*postsArr))
	if err != nil {
		log.Fatalln(err)
	}

	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	if ids[len(ids) - 1] == 1500000 {
		Vaccuum()
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

	_, err := strconv.Atoi(slugOrID.(string))

	if err != nil {
		err = db.QueryRow("getThreadBySlug", slugOrID).
			Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
		return &thread, err
	}

	err = db.QueryRow("getThreadById", slugOrID).Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message, &thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count)
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

func getThreadPostsFlat(ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsFlatSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = db.Query("getPostsFlatSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsFlatSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = db.Query("getPostsFlatSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsFlatLimitDesc", ID, limit)
			} else {
				rows, err = db.Query("getPostsFlatLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsFlatLimitDesc", ID, nil)
			} else {
				rows, err = db.Query("getPostsFlatLimit", ID, nil)
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

func getThreadPostsTree(ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = db.Query("getPostsTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = db.Query("getPostsTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsTreeLimitDesc", ID, limit)
			} else {
				rows, err = db.Query("getPostsTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsTreeLimitDesc", ID, nil)
			} else {
				rows, err = db.Query("getPostsTreeLimit", ID, nil)
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

func getThreadPostsParentTree(ID int, limit []byte, since []byte, desc []byte) (*models.PostArr, int) {
	var err error
	var rows *pgx.Rows

	if since != nil {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsParentTreeSinceLimitDesc", ID, limit, since)
			} else {
				rows, err = db.Query("getPostsParentTreeSinceLimit", ID, limit, since)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsParentTreeSinceLimitDesc", ID, nil, since)
			} else {
				rows, err = db.Query("getPostsParentTreeSinceLimit", ID, nil, since)
			}
		}
	} else {
		if limit != nil {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsParentTreeLimitDesc", ID, limit)
			} else {
				rows, err = db.Query("getPostsParentTreeLimit", ID, limit)
			}
		} else {
			if bytes.Equal(desc, []byte("true")) {
				rows, err = db.Query("getPostsParentTreeLimitDesc", ID, nil)
			} else {
				rows, err = db.Query("getPostsParentTreeLimit", ID, nil)
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

	if _, err = strconv.Atoi(*slugOrID); err != nil {
		if err = db.QueryRow("checkThreadIdBySlug", slugOrID).Scan(&ID); err != nil {
			return nil, 404
		}
	} else {
		if err = db.QueryRow("checkThreadIdById", slugOrID).Scan(&ID); err != nil {
			return nil, 404
		}
	}

	switch true {
	case bytes.Equal([]byte("tree"), sort):
		return getThreadPostsTree(ID, limit, since, desc)
	case bytes.Equal([]byte("parent_tree"), sort):
		return getThreadPostsParentTree(ID, limit, since, desc)
	default:
		return getThreadPostsFlat(ID, limit, since, desc)
	}
}

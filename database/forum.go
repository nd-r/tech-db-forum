package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nd-r/tech-db-forum/models"
	"log"
	"strings"
)

const createForumQuery = "INSERT INTO forum (slug, title, moderator) " +
	"VALUES ($1, $2, (SELECT nickname FROM users WHERE lower(nickname) = lower($3)))" +
	"RETURNING *"

func CreateForum(forum *models.Forum) *pq.Error {
	tx := db.MustBegin()
	defer tx.Rollback()

	err := tx.Get(forum, createForumQuery, forum.Slug, forum.Title, forum.Moderator)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Fatalln(err)
		}
		return pqErr
	}

	tx.Commit()
	return nil
}

const selectForumQuery = "SELECT * FROM forum WHERE lower(slug)=lower($1)"

func GetForumDetails(slug interface{}) (*models.Forum, error) {
	tx := db.MustBegin()
	defer tx.Commit()

	forum := models.Forum{}

	err := tx.Get(&forum, selectForumQuery, slug)
	if err != nil {
		return nil, err
	}

	return &forum, nil
}

const createThreadQuery = "INSERT INTO thread (slug, title, message, forum_slug, created, user_nick) " +
	"VALUES ($1, $2, $3, (SELECT slug FROM forum WHERE lower(slug)=lower($4)), COALESCE($5::TIMESTAMPTZ, current_timestamp), " +
	"(SELECT nickname FROM users WHERE lower(nickname) = lower($6))) " +
	"RETURNING  id, slug, title, message, forum_slug, user_nick, created, votes_count"

func CreateThread(slug *interface{}, threadDetails *models.Thread) *pq.Error {
	tx := db.MustBegin()
	defer tx.Commit()

	err := tx.Get(threadDetails, createThreadQuery, threadDetails.Slug, threadDetails.Title, threadDetails.Message,
		slug, threadDetails.Created, threadDetails.User_nick)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Fatalln(err)
		}

		tx.Rollback()
		return pqErr
	}
	return nil
}

const getForumThreads = "SELECT * FROM thread WHERE lower(forum_slug)=lower($1)"

func GetForumThreads(slug *interface{}, limit []byte, since []byte, desc []byte) (*models.TreadArr, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	isDesc := string(desc) == "true"

	var err error
	var query string
	var threads models.TreadArr

	if since != nil {
		if !isDesc {
			query = "WITH findForum AS (SELECT lower(slug) as l_forum_slug FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT * FROM thread WHERE lower(forum_slug)=(SELECT l_forum_slug FROM findForum) AND created >= $2 ORDER BY created LIMIT $3) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT l_forum_slug FROM findForum), ''), '', NULL, NULL)"
		} else {
			query = "WITH findForum AS (SELECT lower(slug) as l_forum_slug FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT * FROM thread WHERE lower(forum_slug)=(SELECT l_forum_slug FROM findForum) AND created <= $2 ORDER BY created DESC LIMIT $3) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT l_forum_slug FROM findForum), ''), '', NULL, NULL)"
		}
		if limit != nil {
			err = tx.Select(&threads, query, slug, since, limit)
		} else {
			err = tx.Select(&threads, query, slug, since, nil)
		}
	} else {
		if !isDesc {
			query = "WITH findForum AS (SELECT lower(slug) as l_forum_slug FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT * FROM thread WHERE lower(forum_slug)=(SELECT l_forum_slug FROM findForum) ORDER BY created LIMIT $2) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT l_forum_slug FROM findForum), ''), '', NULL, NULL)"
		} else {
			query = "WITH findForum AS (SELECT lower(slug) as l_forum_slug FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT * FROM thread WHERE lower(forum_slug)=(SELECT l_forum_slug FROM findForum) ORDER BY created DESC LIMIT $2) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT l_forum_slug FROM findForum), ''), '', NULL, NULL)"
		}
		if limit != nil {
			err = tx.Select(&threads, query, slug, limit)
		} else {
			err = tx.Select(&threads, query, slug, nil)
		}
	}

	if err != nil {
		log.Fatalln(err)
		tx.Rollback()
	}

	if len(threads) == 1 {
		if threads[0].Forum_slug != "" {
			return nil, 200
		}
		return nil, 404
	}
	
	threads = threads[:len(threads)-1]
	return &threads, 200
}

const getForumUsersQuery = "SELECT us.about, us.email, us.fullname, us.nickname FROM forum_users f JOIN users us ON us.id = f.userID WHERE f.forumid = (SELECT id FROM forum WHERE lower(slug) = lower($1)) AND lower(nickname) > lower(coalesce($2, '')) ORDER BY lower(nickname) $4 LIMIT $3:: INTEGER"
const checkForumSlug = "SELECT slug FROM forum WHERE lower(slug)=lower($1)"

func GetForumUsers(slug *string, limit []byte, since []byte, desc []byte) (*models.UsersArr, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	err := tx.Get(&slug, checkForumSlug, slug)
	if err != nil {
		return nil, 404
	}
	users := models.UsersArr{}

	query := getForumUsersQuery

	if string(desc) == "true" {
		if since != nil {
			query = strings.Replace(query, ">", "<", -1)
		}
		query = strings.Replace(query, "$4", " DESC", -1)
	} else {
		query = strings.Replace(query, "$4", " ASC", -1)
	}

	if limit == nil {
		tx.Select(&users, query, slug, since, nil)
	} else {
		tx.Select(&users, query, slug, since, limit)
	}
	return &users, 200
}

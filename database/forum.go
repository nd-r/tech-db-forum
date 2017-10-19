package database

import (
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
	query := getForumThreads

	if since != nil {
		if !isDesc {
			query += " AND created >= $2"
		} else {
			query += " AND created <= $2"
		}
	}

	query += " ORDER BY created"
	if isDesc {
		query += " DESC"
	}

	if limit != nil && since != nil {
		query += " LIMIT $3"
	} else if limit != nil {
		query += " LIMIT $2"
	}

	var threads models.TreadArr
	if since != nil {
		tx.Select(&threads, query, slug, since, limit)
		if len(threads) == 0 {
			tx.Select(&threads, getForumThreads, slug)
			if len(threads) == 0 {
				return nil, 404
			}
			return nil, 200
		}
		return &threads, 200
	}
	tx.Select(&threads, query, slug, limit)
	if len(threads) == 0 {
		tx.Select(&threads, getForumThreads, slug)
		if len(threads) == 0 {
			return nil, 404
		}
		return nil, 200
	}
	return &threads, 200
}

const getForumUsersQuery = "SELECT us.about, us.email, us.fullname, us.nickname FROM forum_users f JOIN forum fo ON fo.id = f.forumId JOIN users us ON us.id = f.userID WHERE lower(fo.slug) = lower($1) AND lower(nickname) > lower(coalesce($2, '')) GROUP BY f.forumid, f.userid, us.about, us.email, us.fullname, us.nickname ORDER BY lower(nickname) $4 LIMIT $3::INTEGER"
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

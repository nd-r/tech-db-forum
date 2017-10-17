package database

import (
	"github.com/nd-r/tech-db-forum/models"
)

const createForumQuery = "INSERT INTO forum " +
	"(slug, title, moderator) " +
	"VALUES ($1, $2, " +
	"(SELECT nickname FROM users " +
	"WHERE lower(nickname) = lower($3))) RETURNING slug, title, moderator, posts, threads"

const selectForumQuery = "SELECT slug, title, posts, threads, moderator FROM forum WHERE lower(slug)=lower($1)"

func CreateForum(forum *models.Forum) (*models.Forum, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	forumExisting := models.Forum{}
	err := tx.Get(&forumExisting, selectForumQuery, forum.Slug)
	if err == nil {
		return &forumExisting, 409
	}

	err = tx.Get(forum, createForumQuery, forum.Slug, forum.Title, forum.Moderator)
	if err != nil {
		return nil, 404
	}
	return nil, 201
}

func GetForumDetails(slug string) (*models.Forum, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	forum := models.Forum{}

	err := tx.Get(&forum, selectForumQuery, slug)
	if err == nil {
		return &forum, 200
	}

	return nil, 404
}

const getThreadDetailsBySlug = "SELECT * FROM thread WHERE lower(slug)=lower($1)"

const createThreadQuery = "INSERT INTO thread " +
	"(slug, title, message, forum_title, created, user_nick) " +
	"VALUES ($1, $2, $3, (SELECT slug FROM forum WHERE lower(slug)=lower($4)), COALESCE($5::TIMESTAMPTZ, current_timestamp)," +
	"(SELECT nickname FROM users " +
	"WHERE lower(nickname) = lower($6))) RETURNING  id, slug, title, message, forum_title, user_nick, created, votes_count"

func CreateThread(thread *models.Thread) (*models.Thread, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	threadExisting := models.Thread{}

	if thread.Slug != nil {
		err := tx.Get(&threadExisting, getThreadDetailsBySlug, *thread.Slug)
		if err == nil {
			return &threadExisting, 409
		}
	}

	err := tx.Get(thread, createThreadQuery, thread.Slug, thread.Title, thread.Message, thread.Forum_title, thread.Created, thread.User_nick)
	if err != nil {
		tx.Rollback()
		return nil, 404
	}
	return nil, 201
}

const getForumThreads = "SELECT * FROM thread WHERE lower(forum_title)=lower($1)"

func GetForumThreads(slug *string, limit []byte, since []byte, desc []byte) (*models.TreadArr, int) {
	tx := db.MustBegin()
	defer tx.Commit()
	query := getForumThreads

	if since != nil {
		if string(desc) != "true" {
			query += " AND created >= $2"
		} else {
			query += " AND created <= $2"
		}
	}

	query += " ORDER BY created"
	if string(desc) == "true" {
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

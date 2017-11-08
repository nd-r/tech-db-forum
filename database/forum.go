package database

import (
	"github.com/jackc/pgx"
	// 	"github.com/jackc/pgx"
	// 	"github.com/lib/pq"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"log"
		"strings"
)

const createForumQuery = "INSERT INTO forum (slug, title, moderator) " +
	"VALUES ($1, $2, (SELECT nickname FROM users WHERE lower(nickname) = lower($3)))" +
	"RETURNING id, slug, title, posts, threads, moderator"

func CreateForum(forum *models.Forum) (*models.Forum, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback()

	forumExisting := models.Forum{}

	if err = tx.QueryRow(selectForumQuery, &forum.Slug).
		Scan(&forumExisting.Id, &forumExisting.Slug, &forumExisting.Title,
			&forumExisting.Posts, &forumExisting.Threads, &forumExisting.Moderator); err == nil {
		return &forumExisting, dberrors.ErrForumExists
	}

	if err = tx.QueryRow(createForumQuery, forum.Slug, forum.Title, forum.Moderator).
		Scan(&forumExisting.Id, &forumExisting.Slug, &forumExisting.Title,
			&forumExisting.Posts, &forumExisting.Threads, &forumExisting.Moderator); err != nil {
		return nil, dberrors.ErrUserNotFound
	}

	tx.Commit()
	return &forumExisting, nil
}

const selectForumQuery = "SELECT id, slug, title, posts, threads, moderator FROM forum WHERE lower(slug)=lower($1)"

func GetForumDetails(slug interface{}) (*models.Forum, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	forum := models.Forum{}
	err = tx.QueryRow(selectForumQuery, slug).
		Scan(&forum.Id, &forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forum.Moderator)

	return &forum, err
}

const createThreadQuery = "WITH select_forum_details AS (SELECT id, slug FROM forum WHERE lower(slug) = lower($4))," +
	"select_user_details AS (SELECT about, fullname, nickname, email FROM users WHERE lower(nickname) = lower($6))," +
	"upd_forum_users AS (INSERT INTO forum_users (forumid, nickname, about, email, fullname) VALUES ((SELECT id FROM select_forum_details),(SELECT nickname FROM select_user_details),(SELECT about FROM select_user_details),(SELECT email FROM select_user_details),(SELECT fullname FROM select_user_details)) ON CONFLICT DO NOTHING)," +
	"upd_forum AS (UPDATE forum SET threads = threads + 1 WHERE id = (SELECT id FROM select_forum_details))" +
	"INSERT INTO thread (slug, title, message, forum_id, forum_slug, created, user_nick)" +
	"VALUES ($1,$2,$3,(SELECT id FROM select_forum_details), (SELECT slug FROM select_forum_details), COALESCE($5::TIMESTAMPTZ, current_timestamp), (SELECT nickname FROM select_user_details))" +
	"RETURNING id, forum_slug, user_nick"

func CreateThread(forumSlug interface{}, threadDetails *models.Thread) (*models.Thread,error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback()

	if threadDetails.Slug != nil {
		existingThread := models.Thread{}

		if err = tx.QueryRow("getThreadBySlug", threadDetails.Slug).
			Scan(&existingThread.Id, &existingThread.Slug, &existingThread.Title,
				&existingThread.Message, &existingThread.Forum_slug, &existingThread.User_nick, &existingThread.Created,
				&existingThread.Votes_count); err == nil {

			return &existingThread, dberrors.ErrThreadExists
		}
	}

	if err = tx.QueryRow("InsertIntoThread", threadDetails.Slug, &threadDetails.Title,
		&threadDetails.Message, forumSlug, threadDetails.Created, &threadDetails.User_nick).
		Scan(&threadDetails.Id, &threadDetails.Forum_slug, &threadDetails.User_nick); err != nil {

		if pqErr, ok := err.(pgx.PgError); ok {
			switch pqErr.ColumnName {
			case "forum_slug", "forum_id":
				return nil, dberrors.ErrForumNotFound
			case "user_nick":
				return nil, dberrors.ErrUserNotFound
			}
		} else {
			log.Fatalln(err)
		}
	}

	tx.Commit()
	return threadDetails, nil
}

func GetForumThreads(slug interface{}, limit []byte, since []byte, desc []byte) (*models.TreadArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	isDesc := string(desc) == "true"

	var query string
	var threads models.TreadArr

	var rows *pgx.Rows

	if since != nil {
		if !isDesc {
			query = "WITH findForum AS (SELECT id FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT id, slug, title, message, forum_slug, user_nick, created, votes_count FROM thread WHERE forum_id=(SELECT id FROM findForum) AND created >= $2::TEXT::TIMESTAMPTZ ORDER BY created LIMIT $3::TEXT::INTEGER) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT id::TEXT FROM findForum), ''), '', '2017-11-05 22:08:13.326059 +03:00', NULL)"
		} else {
			query = "WITH findForum AS (SELECT id FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT id, slug, title, message, forum_slug, user_nick, created, votes_count FROM thread WHERE forum_id=(SELECT id FROM findForum) AND created <= $2::TEXT::TIMESTAMPTZ ORDER BY created DESC LIMIT $3::TEXT::INTEGER) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT id::TEXT FROM findForum), ''), '', '2017-11-05 22:08:13.326059 +03:00', NULL)"
		}
		if limit != nil {
			rows, err = tx.Query(query, slug, since, limit)
		} else {
			rows, err = tx.Query(query, slug, since, nil)
		}
	} else {
		if !isDesc {
			query = "WITH findForum AS (SELECT id FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT id, slug, title, message, forum_slug, user_nick, created, votes_count FROM thread WHERE forum_id=(SELECT id FROM findForum) ORDER BY created LIMIT $2::TEXT::INTEGER) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT id::TEXT FROM findForum), ''), '', '2017-11-05 22:08:13.326059 +03:00', NULL)"
		} else {
			query = "WITH findForum AS (SELECT id FROM forum WHERE lower(slug) = lower($1)) SELECT * FROM (SELECT id, slug, title, message, forum_slug, user_nick, created, votes_count FROM thread WHERE forum_id=(SELECT id FROM findForum) ORDER BY created DESC LIMIT $2::TEXT::INTEGER) s UNION ALL (SELECT 0, '', '','', COALESCE((SELECT id::TEXT FROM findForum), ''), '', '2017-11-05 22:08:13.326059 +03:00', NULL)"
		}
		if limit != nil {
			rows, err = tx.Query(query, slug, limit)
		} else {
			rows, err = tx.Query(query, slug, nil)
		}
	}

	if err != nil {
		log.Fatalln(err)
		tx.Rollback()
	}
	defer rows.Close()	

	for rows.Next() {
		thread := models.Thread{}

		if err = rows.Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message,
			&thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count); err != nil {
			log.Fatalln(err)
		}

		threads = append(threads, &thread)
	}

	if len(threads) == 1 {
		if threads[0].Forum_slug != "" {
			return nil, nil
		}
		return nil, dberrors.ErrForumNotFound
	}

	threads = threads[:len(threads)-1]
	return &threads, nil
}

const getForumUsersQuery = "SELECT about, email, fullname, nickname FROM forum_users f WHERE f.forumid = $1 AND lower(nickname) > lower(coalesce($2, '')) ORDER BY lower(nickname) $4 LIMIT $3::TEXT:: INTEGER"
const checkForumSlug = "SELECT id FROM forum WHERE lower(slug)=lower($1)"

func GetForumUsers(slug *string, limit []byte, since []byte, desc []byte) (*models.UsersArr, int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	var ID int
	if err = tx.QueryRow(checkForumSlug, slug).Scan(&ID); err != nil {
		return nil, 404
	}

	var users models.UsersArr

	query := getForumUsersQuery

	if string(desc) == "true" {
		if since != nil {
			query = strings.Replace(query, ">", "<", -1)
		}
		query = strings.Replace(query, "$4", " DESC", -1)
	} else {
		query = strings.Replace(query, "$4", " ASC", -1)
	}

	var rows *pgx.Rows
	if limit == nil {
		rows, err = tx.Query(query, ID, since, nil)
	} else {
		rows, err = tx.Query(query, ID, since, limit)
	}

	if err != nil{ 
		log.Fatalln(err)
	}
	for rows.Next(){
		user := models.User{}
		if err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil{
			log.Fatalln(err)
		}
		users = append(users, &user)
	}
	return &users, 200
}

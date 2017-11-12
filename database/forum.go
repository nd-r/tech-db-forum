package database

import (
	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"log"
)

const createForumQuery = "INSERT INTO forum (slug, title, moderator) " +
	"VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname=$3)) " +
	"RETURNING moderator::TEXT"

const selectForumQuery = "SELECT slug::TEXT, title, posts, threads, moderator::TEXT " +
	"FROM forum " +
	"WHERE slug=$1"

func CreateForum(forum *models.Forum) (*models.Forum, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	forumExisting := models.Forum{}

	if err = tx.QueryRow(selectForumQuery, &forum.Slug).
		Scan(&forumExisting.Slug, &forumExisting.Title,
			&forumExisting.Posts, &forumExisting.Threads, &forumExisting.Moderator); err == nil {

		tx.Rollback()
		return &forumExisting, dberrors.ErrForumExists
	}

	if err := tx.QueryRow(createForumQuery, &forum.Slug, &forum.Title, &forum.Moderator).Scan(&forum.Moderator); err != nil {
		tx.Rollback()
		return nil, dberrors.ErrUserNotFound
	}

	tx.Commit()
	return forum, nil
}

func GetForumDetails(slug interface{}) (*models.Forum, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	forum := models.Forum{}
	err = tx.QueryRow(selectForumQuery, &slug).
		Scan(&forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forum.Moderator)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &forum, err
}

const claimUserInfo = "SELECT id, nickname::text, email::text, about, fullname FROM users WHERE nickname = $1"

const getForumIDAndSlugBySlug = "SELECT id, slug::text FROM forum WHERE slug = $1"

const insertIntoThread = "INSERT INTO thread (slug, title, message, forum_id, forum_slug, user_id, user_nick, created) " +
	"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING RETURNING id"

const insertIntoForumUsers = "INSERT INTO forum_users (forumid, nickname, email, about, fullname) " +
	"VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING"

func CreateThread(forumSlug interface{}, threadDetails *models.Thread) (*models.Thread, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	user := models.User{}
	var userID, forumID int
	var realForumSlug string

	if err = tx.QueryRow(claimUserInfo, &threadDetails.User_nick).
		Scan(&userID, &user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {

		tx.Rollback()
		log.Println(err)
		
		return nil, dberrors.ErrUserNotFound
	}

	if err = tx.QueryRow(getForumIDAndSlugBySlug, &forumSlug).
		Scan(&forumID, &realForumSlug); err != nil {

		tx.Rollback()
		log.Println(err)
		
		return nil, dberrors.ErrForumNotFound
	}

	if err = tx.QueryRow(insertIntoThread, threadDetails.Slug, &threadDetails.Title,
		&threadDetails.Message, forumID, &realForumSlug, userID, &user.Nickname, &threadDetails.Created).
		Scan(&threadDetails.Id); err != nil {

		existingThread := models.Thread{}

		if err = tx.QueryRow(getThreadBySlug, threadDetails.Slug).
			Scan(&existingThread.Id, &existingThread.Slug, &existingThread.Title,
				&existingThread.Message, &existingThread.Forum_slug, &existingThread.User_nick, &existingThread.Created,
				&existingThread.Votes_count); err == nil {

			tx.Rollback()
			return &existingThread, dberrors.ErrThreadExists
		}

		tx.Rollback()
		log.Fatalln(err)
	}

	if _, err = tx.Exec(insertIntoForumUsers, forumID, user.Nickname, user.Email, user.About, user.Fullname); err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	threadDetails.Forum_slug = realForumSlug
	threadDetails.User_nick = user.Nickname

	tx.Commit()
	return threadDetails, nil
}

const getForumIDBySlug = "SELECT id FROM forum WHERE slug=$1"

const gft = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 " +
	"ORDER BY created"

const gftDesc = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 " +
	"ORDER BY created DESC"

const gftLimit = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 " +
	"ORDER BY created " +
	"LIMIT $2::TEXT::INTEGER"

const gftLimitDesc = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 " +
	"ORDER BY created DESC " +
	"LIMIT $2::TEXT::INTEGER"

const gftCreated = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 AND created >= $2::TEXT::TIMESTAMPTZ " +
	"ORDER BY created "

const gftCreatedDesc = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 AND created <= $2::TEXT::TIMESTAMPTZ " +
	"ORDER BY created DESC"

const gftCreatedLimit = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 AND created >= $2::TEXT::TIMESTAMPTZ " +
	"ORDER BY created " +
	"LIMIT $3::TEXT::INTEGER"

const gftCreatedLimitDesc = "SELECT id, slug::TEXT, title, message, forum_slug::TEXT, user_nick::TEXT, created, votes_count FROM thread " +
	"WHERE forum_id = $1 AND created <= $2::TEXT::TIMESTAMPTZ " +
	"ORDER BY created DESC " +
	"LIMIT $3::TEXT::INTEGER"

func GetForumThreads(slug interface{}, limit []byte, since []byte, desc []byte) (*models.TreadArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	isDesc := string(desc) == "true"

	var forumID int

	if err = tx.QueryRow(getForumIDBySlug, &slug).Scan(&forumID); err != nil {
		tx.Rollback()
		return nil, dberrors.ErrForumNotFound
	}

	var rows *pgx.Rows

	if limit == nil {
		if since == nil {
			if isDesc {
				// no limit, no since, DESC
				rows, err = tx.Query(gftDesc, forumID)
			} else {
				// no limit, no since, ASC
				rows, err = tx.Query(gft, forumID)
			}
		} else {
			if isDesc {
				// no limit, yes since, DESC
				rows, err = tx.Query(gftCreatedDesc, forumID, since)
			} else {
				// no limit, yes since, ASC
				rows, err = tx.Query(gftCreated, forumID, since)
			}
		}
	} else {
		if since == nil {
			if isDesc {
				// yes limit, no since, DESC
				rows, err = tx.Query(gftLimitDesc, forumID, limit)
			} else {
				// yse limit, no since, ASC
				rows, err = tx.Query(gftLimit, forumID, limit)
			}
		} else {
			if isDesc {
				// yes limit, yes since, DESC
				rows, err = tx.Query(gftCreatedLimitDesc, forumID, since, limit)
			} else {
				// yes limit, yes since, ASC
				rows, err = tx.Query(gftCreatedLimit, forumID, since, limit)
			}
		}
	}

	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	var threads models.TreadArr

	for rows.Next() {
		thread := models.Thread{}

		if err = rows.Scan(&thread.Id, &thread.Slug, &thread.Title, &thread.Message,
			&thread.Forum_slug, &thread.User_nick, &thread.Created, &thread.Votes_count); err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}

		threads = append(threads, &thread)
	}

	rows.Close()
	tx.Commit()
	return &threads, nil
}

const getForumUsersQuery = "SELECT about, email, fullname, nickname FROM forum_users f WHERE f.forumid = $1 AND lower(nickname) > lower(coalesce($2, '')) ORDER BY lower(nickname) $4 LIMIT $3::TEXT:: INTEGER"

const gfu = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 " +
	"ORDER BY lower(nickname)"

const gfuDesc = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 " +
	"ORDER BY lower(nickname) DESC"

const gfuLimit = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 " +
	"ORDER BY lower(nickname) " +
	"LIMIT $2::TEXT::INTEGER"

const gfuLimitDesc = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 " +
	"ORDER BY lower(nickname) DESC " +
	"LIMIT $2::TEXT::INTEGER"

const gfuSince = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 AND nickname > $2::TEXT::CITEXT " +
	"ORDER BY lower(nickname)"

const gfuSinceDesc = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 AND nickname < $2::TEXT::CITEXT " +
	"ORDER BY lower(nickname) DESC"

const gfuSinceLimit = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 AND nickname > $2::TEXT::CITEXT " +
	"ORDER BY lower(nickname) " +
	"LIMIT $3::TEXT::INTEGER"

const gfuSinceLimitDesc = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM forum_users " +
	"WHERE forumid = $1 AND nickname < $2::TEXT::CITEXT " +
	"ORDER BY lower(nickname) DESC " +
	"LIMIT $3::TEXT::INTEGER"

func GetForumUsers(slug interface{}, limit []byte, since []byte, desc []byte) (*models.UsersArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	var forumID int
	if err = tx.QueryRow(getForumIDBySlug, &slug).Scan(&forumID); err != nil {
		tx.Rollback()
		return nil, dberrors.ErrForumNotFound
	}

	isDesc := string(desc) == "true"
	var rows *pgx.Rows

	if limit == nil {
		if since == nil {
			if isDesc {
				// no limit, no since, desc
				rows, err = tx.Query(gfuDesc, forumID)
			} else {
				// no limit, no since, asc
				rows, err = tx.Query(gfu, forumID)				
			}
		} else {
			if isDesc {
				// no limit, yes since, desc
				rows, err = tx.Query(gfuSinceDesc, forumID, since)				
			} else {
				// no limit, yes since, asc
				rows, err = tx.Query(gfuSince, forumID, since)				
			}
		}
	} else {
		if since == nil {
			if isDesc {
				// yes limit, no since, desc
				rows, err = tx.Query(gfuLimitDesc, forumID, limit)				
			} else {
				// yes limit, no since, asc
				rows, err = tx.Query(gfuLimit, forumID, limit)				
			}
		} else {
			if isDesc {
				// yes limit, yes since, desc
				rows, err = tx.Query(gfuSinceLimitDesc, forumID, since, limit)				
			} else {
				// yes limit, yes since, asc
				rows, err = tx.Query(gfuSinceLimit, forumID, since, limit)				
			}
		}
	}

	if err != nil{
		rows.Close()
		tx.Rollback()
		log.Fatalln(err, limit == nil, since == nil, isDesc)
	}
	var users models.UsersArr

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
			rows.Close()
			tx.Rollback()
			log.Fatalln(err)
		}
		users = append(users, &user)
	}

	rows.Close()
	tx.Commit()
	return &users, nil
}

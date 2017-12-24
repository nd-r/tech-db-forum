package database

import (
	"log"

	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
)

func CreateForum(forum *models.Forum) (*models.Forum, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	forumExisting := models.Forum{}

	if err = tx.QueryRow("selectForumQuery", &forum.Slug).
		Scan(&forumExisting.Slug, &forumExisting.Title,
		&forumExisting.Posts, &forumExisting.Threads, &forumExisting.Moderator); err == nil {

		tx.Rollback()
		return &forumExisting, dberrors.ErrForumExists
	}

	if err := tx.QueryRow("createForumQuery", &forum.Slug, &forum.Title, &forum.Moderator).Scan(&forum.Moderator); err != nil {
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
	err = tx.QueryRow("selectForumQuery", &slug).
		Scan(&forum.Slug, &forum.Title, &forum.Posts, &forum.Threads, &forum.Moderator)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &forum, err
}

func CreateThread(forumSlug interface{}, threadDetails *models.Thread) (*models.Thread, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	user := models.User{}
	var userID, forumID int
	var realForumSlug string

	if err = tx.QueryRow("claimUserInfo", &threadDetails.User_nick).
		Scan(&userID, &user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {

		tx.Rollback()

		return nil, dberrors.ErrUserNotFound
	}

	if err = tx.QueryRow("getForumIDAndSlugBySlug", &forumSlug).
		Scan(&forumID, &realForumSlug); err != nil {

		tx.Rollback()

		return nil, dberrors.ErrForumNotFound
	}

	if err = tx.QueryRow("insertIntoThread", threadDetails.Slug, &threadDetails.Title,
		&threadDetails.Message, forumID, &realForumSlug, userID, &user.Nickname, &threadDetails.Created).
		Scan(&threadDetails.Id); err != nil {

		existingThread := models.Thread{}

		if err = tx.QueryRow("getThreadBySlug", threadDetails.Slug).
			Scan(&existingThread.Id, &existingThread.Slug, &existingThread.Title,
			&existingThread.Message, &existingThread.Forum_slug, &existingThread.User_nick, &existingThread.Created,
			&existingThread.Votes_count); err == nil {

			tx.Rollback()
			return &existingThread, dberrors.ErrThreadExists
		}

		tx.Rollback()
		log.Fatalln(err)
	}

	if _, err = tx.Exec("insertIntoForumUsers", forumID, user.Nickname, user.Email, user.About, user.Fullname); err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	threadDetails.Forum_slug = realForumSlug
	threadDetails.User_nick = user.Nickname

	tx.Commit()
	return threadDetails, nil
}

func GetForumThreads(slug interface{}, limit []byte, since []byte, desc []byte) (*models.TreadArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	isDesc := string(desc) == "true"

	var forumID int

	if err = tx.QueryRow("getForumIDBySlug", &slug).Scan(&forumID); err != nil {
		tx.Rollback()
		return nil, dberrors.ErrForumNotFound
	}

	var rows *pgx.Rows

	if limit == nil {
		if since == nil {
			if isDesc {
				// no limit, no since, DESC
				rows, err = tx.Query("gftLimitDesc", forumID, nil)
			} else {
				// no limit, no since, ASC
				rows, err = tx.Query("gftLimit", forumID, nil)
			}
		} else {
			if isDesc {
				// no limit, yes since, DESC
				rows, err = tx.Query("gftCreatedLimitDesc", forumID, since, nil)
			} else {
				// no limit, yes since, ASC
				rows, err = tx.Query("gftCreatedLimit", forumID, since, nil)
			}
		}
	} else {
		if since == nil {
			if isDesc {
				// yes limit, no since, DESC
				rows, err = tx.Query("gftLimitDesc", forumID, limit)
			} else {
				// yse limit, no since, ASC
				rows, err = tx.Query("gftLimit", forumID, limit)
			}
		} else {
			if isDesc {
				// yes limit, yes since, DESC
				rows, err = tx.Query("gftCreatedLimitDesc", forumID, since, limit)
			} else {
				// yes limit, yes since, ASC
				rows, err = tx.Query("gftCreatedLimit", forumID, since, limit)
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

func GetForumUsers(slug interface{}, limit []byte, since []byte, desc []byte) (*models.UsersArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	var forumID int
	if err = tx.QueryRow("getForumIDBySlug", &slug).Scan(&forumID); err != nil {
		tx.Rollback()
		return nil, dberrors.ErrForumNotFound
	}

	isDesc := string(desc) == "true"
	var rows *pgx.Rows

	if limit == nil {
		if since == nil {
			if isDesc {
				// no limit, no since, desc
				rows, err = tx.Query("gfuLimitDesc", forumID, nil)
			} else {
				// no limit, no since, asc
				rows, err = tx.Query("gfuLimit", forumID, nil)
			}
		} else {
			if isDesc {
				// no limit, yes since, desc
				rows, err = tx.Query("gfuSinceLimitDesc", forumID, since, nil)
			} else {
				// no limit, yes since, asc
				rows, err = tx.Query("gfuSinceLimit", forumID, since, nil)
			}
		}
	} else {
		if since == nil {
			if isDesc {
				// yes limit, no since, desc
				rows, err = tx.Query("gfuLimitDesc", forumID, limit)
			} else {
				// yes limit, no since, asc
				rows, err = tx.Query("gfuLimit", forumID, limit)
			}
		} else {
			if isDesc {
				// yes limit, yes since, desc
				rows, err = tx.Query("gfuSinceLimitDesc", forumID, since, limit)
			} else {
				// yes limit, yes since, asc
				rows, err = tx.Query("gfuSinceLimit", forumID, since, limit)
			}
		}
	}

	if err != nil {
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

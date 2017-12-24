package resources

import (
	"log"

	"github.com/jackc/pgx"
)

/**
Thread queries
 */

//
// INSERTING
//
const insertIntoThread = `INSERT INTO thread (slug,
	title,
	message,
	forum_id,
	forum_slug,
	user_id,
	user_nick,
	created)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT DO NOTHING
RETURNING id`

//
// UPDATING
//

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

//
// CLAIMING
//

const getThreadById = `SELECT id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count
FROM thread
WHERE id=$1`

const getThreadBySlug = `SELECT id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count
FROM thread
WHERE slug=$1`

const checkThreadIdById = `SELECT id
FROM thread
WHERE id=$1`

const checkThreadIdBySlug = `SELECT id
FROM thread
WHERE slug=$1`

const getThreadIdAndForumSlugBySlug = `SELECT id,
	forum_slug::TEXT
FROM thread
WHERE slug=$1`

const getThreadIdAndForumSlugById = `SELECT id,
	forum_slug::TEXT
FROM thread
WHERE id=$1`

const gftLimit = `SELECT id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count
FROM thread
WHERE forum_id = $1
ORDER BY created
LIMIT $2::TEXT::INTEGER`

const gftLimitDesc = `SELECT id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count
FROM thread
WHERE forum_id = $1
ORDER BY created DESC
LIMIT $2::TEXT::INTEGER`

const gftCreatedLimit = `SELECT id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count
FROM thread
WHERE forum_id = $1 AND created >= $2::TEXT::TIMESTAMPTZ
ORDER BY created
LIMIT $3::TEXT::INTEGER`

const gftCreatedLimitDesc = `SELECT id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count
FROM thread
WHERE forum_id = $1 AND created <= $2::TEXT::TIMESTAMPTZ
ORDER BY created DESC
LIMIT $3::TEXT::INTEGER`

func PrepateThreadQueries(tx *pgx.ConnPool) {
	if _, err := tx.Prepare("insertIntoThread", insertIntoThread); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("threadUpdateQuery", threadUpdateQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getThreadById", getThreadById); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getThreadBySlug", getThreadBySlug); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("checkThreadIdById", checkThreadIdById); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("checkThreadIdBySlug", checkThreadIdBySlug); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getThreadIdAndForumSlugBySlug", getThreadIdAndForumSlugBySlug); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getThreadIdAndForumSlugById", getThreadIdAndForumSlugById); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gftLimit", gftLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gftLimitDesc", gftLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gftCreatedLimit", gftCreatedLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gftCreatedLimitDesc", gftCreatedLimitDesc); err != nil {
		log.Fatalln(err)
	}

}

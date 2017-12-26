package resources

import (
	"log"

	"github.com/jackc/pgx"
)

/**
Thread queries
 */

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


const gftLimit = `SELECT id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count
FROM thread
WHERE forum_id = (SELECT id FROM forum WHERE slug=$1)
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
WHERE forum_id = (SELECT id FROM forum WHERE slug=$1)
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
WHERE forum_id = (SELECT id FROM forum WHERE slug=$1) AND created >= $2::TEXT::TIMESTAMPTZ
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
WHERE forum_id = (SELECT id FROM forum WHERE slug=$1) AND created <= $2::TEXT::TIMESTAMPTZ
ORDER BY created DESC
LIMIT $3::TEXT::INTEGER`

func PrepateThreadQueries(tx *pgx.ConnPool) {
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

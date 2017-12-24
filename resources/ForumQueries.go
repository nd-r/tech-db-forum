package resources

import (
	"log"

	"github.com/jackc/pgx"
)

/**
Forum queries
 */

//
// INSERTING
//
const createForumQuery = `INSERT INTO forum
(slug, title, moderator)
VALUES (
	$1,
	$2,
	(SELECT nickname FROM users WHERE nickname=$3))
RETURNING moderator::TEXT`

//
// UPDATING
//

const updateForumPosts = `UPDATE forum
SET posts=posts+$2
WHERE slug=$1`

//
// CLAIMING
//

const getForumIDAndSlugBySlug = `SELECT
	id,
	slug::text
FROM forum
WHERE slug = $1`

const selectForumQuery = `SELECT
	slug::TEXT,
	title,
	posts,
	threads,
	moderator::TEXT
FROM forum
WHERE slug=$1`

const getForumIDBySlug = `SELECT id
FROM forum
WHERE slug=$1`

func PrepareForumQueries(tx *pgx.Tx) {
	if _, err := tx.Prepare("createForumQuery", createForumQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("updateForumPosts", updateForumPosts); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getForumIDAndSlugBySlug", getForumIDAndSlugBySlug); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("selectForumQuery", selectForumQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getForumIDBySlug", getForumIDBySlug); err != nil {
		log.Fatalln(err)
	}
}

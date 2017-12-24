package resources

import (
	"log"

	"github.com/jackc/pgx"
)

/**
Forum - users queries
 */

//
//INSERTING
//
const insertIntoForumUsers = `INSERT INTO forum_users (forumid,
	nickname,
	email,
	about,
	fullname)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING`

//
// UPDATING
//

//
// CLAIMING
//

const gfuLimit = `SELECT nickname::TEXT,
	email::TEXT,
	about,
	fullname
FROM forum_users
WHERE forumid = $1
ORDER BY lower(nickname)
LIMIT $2::TEXT::INTEGER`

const gfuLimitDesc = `SELECT nickname::TEXT,
	email::TEXT,
	about,
	fullname
FROM forum_users
WHERE forumid = $1
ORDER BY lower(nickname) DESC
LIMIT $2::TEXT::INTEGER`

const gfuSinceLimit = `SELECT nickname::TEXT,
	email::TEXT,
	about,
	fullname
FROM forum_users
WHERE forumid = $1 AND lower(nickname) > lower($2::TEXT)
ORDER BY lower(nickname)
LIMIT $3::TEXT::INTEGER`

const gfuSinceLimitDesc = `SELECT nickname::TEXT,
	email::TEXT,
	about,
	fullname
FROM forum_users
WHERE forumid = $1 AND lower(nickname) < lower($2::TEXT)
ORDER BY lower(nickname) DESC
LIMIT $3::TEXT::INTEGER`

func PrepareForumUsersQueries(tx *pgx.Tx) {
	if _, err := tx.Prepare("insertIntoForumUsers", insertIntoForumUsers); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gfuLimit", gfuLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gfuLimitDesc", gfuLimitDesc); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gfuSinceLimit", gfuSinceLimit); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("gfuSinceLimitDesc", gfuSinceLimitDesc); err != nil {
		log.Fatalln(err)
	}
}

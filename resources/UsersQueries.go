package resources

import (
	"log"

	"github.com/jackc/pgx"
)

/**
User queries
 */

//
// INSERTING
//
const createUserQuery = `INSERT INTO users
	(about, email, fullname, nickname)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING`

//
// UPDATING
//
const updateUserProfileQuery = `UPDATE users
SET about = COALESCE($1, users.about),
	email = COALESCE($2, users.email),
	fullname = COALESCE($3, users.fullname)
WHERE nickname=$4
RETURNING
	nickname::TEXT,
	email::TEXT,
	about,
	fullname`

//
// CLAIMING
//

const claimUserInfo = `SELECT
	id,
	nickname::text,
	email::text,
	about,
	fullname
FROM users
WHERE nickname = $1`

const selectUsrByNickOrEmailQuery = `SELECT nickname::TEXT,
	email::TEXT,
	about,
	fullname
FROM users
WHERE nickname=$1 OR email=$2`

const getUserProfileQuery = `SELECT
	nickname::TEXT,
	email::TEXT,
	about,
	fullname
FROM users
WHERE nickname = $1`

func PrepareUsersQueries(tx *pgx.Tx) {
	if _, err := tx.Prepare("createUserQuery", createUserQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("updateUserProfileQuery", updateUserProfileQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("claimUserInfo", claimUserInfo); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("selectUsrByNickOrEmailQuery", selectUsrByNickOrEmailQuery); err != nil {
		log.Fatalln(err)
	}

	if _, err := tx.Prepare("getUserProfileQuery", getUserProfileQuery); err != nil {
		log.Fatalln(err)
	}
}

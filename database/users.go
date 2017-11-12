package database

import (
	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"log"
)

const createUserQuery = "INSERT INTO users (about, email, fullname, nickname) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING"
const selectUsrByNickOrEmailQuery = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM users WHERE nickname=$1 OR email=$2"

func CreateUser(user *models.User, nickname interface{}) (*models.UsersArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	res, err := tx.Exec(createUserQuery, &user.About, &user.Email, &user.Fullname, &nickname)
	if err != nil {
		log.Fatalln(err)
	}

	if res.RowsAffected() == 0 {
		existingUsers := models.UsersArr{}

		rows, err := tx.Query(selectUsrByNickOrEmailQuery, &nickname, &user.Email)
		if err != nil {
			log.Fatalln(err)
		}

		for rows.Next() {
			existingUser := models.User{}

			if err = rows.Scan(&existingUser.Nickname, &existingUser.Email, &existingUser.About, &existingUser.Fullname); err != nil {
				log.Fatalln(err)
			}

			existingUsers = append(existingUsers, &existingUser)
		}

		rows.Close()
		tx.Rollback()
		return &existingUsers, dberrors.ErrUserExists
	}

	tx.Commit()
	return nil, nil
}

const getUserProfileQuery = "SELECT nickname::TEXT, email::TEXT, about, fullname FROM users WHERE nickname = $1"

func GetUserProfile(nickname interface{}) (*models.User, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	user := models.User{}

	if err = tx.QueryRow(getUserProfileQuery, &nickname).
		Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
		tx.Rollback()
		return nil, dberrors.ErrUserNotFound
	}

	tx.Commit()
	return &user, nil
}

const updateUserProfileQuery = "UPDATE users SET about = COALESCE($1, users.about), email = COALESCE($2, users.email), fullname = COALESCE($3, users.fullname) WHERE nickname=$4 RETURNING nickname::TEXT, email::TEXT, about, fullname"

func UpdateUserProfile(newData *models.UserUpd, nickname interface{}) (*models.User, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	user := models.User{}

	if err = tx.QueryRow(updateUserProfileQuery, newData.About, newData.Email, newData.Fullname, &nickname).
		Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname); err != nil {
		if _, ok := err.(pgx.PgError); ok {
			return nil, dberrors.ErrUserConflict
		}
		return nil, dberrors.ErrUserNotFound
	}

	return &user, nil
}

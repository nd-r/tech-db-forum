package database

import (
	"github.com/jackc/pgx"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"log"
)

const createUserQuery = "INSERT INTO users (about, email, fullname, nickname) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id"
const selectUsrByNickOrEmailQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1) OR lower(email)=lower($2)"

func CreateUser(usr *models.User, nickname interface{}) (*models.UsersArr, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Rollback()

	var id int
	if err = tx.QueryRow(createUserQuery, &usr.About, &usr.Email, &usr.Fullname, &nickname).Scan(&id); err != nil {
		existingUsers := models.UsersArr{}

		rows, err := tx.Query(selectUsrByNickOrEmailQuery, &nickname, &usr.Email)
		if err != nil {
			log.Fatalln(err)
		}
		defer rows.Close()

		for rows.Next() {
			existingUser := models.User{}
			if err = rows.Scan(&existingUser.About, &existingUser.Email, &existingUser.Fullname, &existingUser.Nickname); err != nil {
				log.Fatalln(err)
			}

			existingUsers = append(existingUsers, &existingUser)
		}
		return &existingUsers, dberrors.ErrUserExists
	}

	tx.Commit()
	return nil, nil
}

const getUserProfileQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1)"

func GetUserProfile(nickname interface{}) (*models.User, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	user := models.User{}

	if err = tx.QueryRow(getUserProfileQuery, &nickname).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
		return nil, dberrors.ErrUserNotFound
	}

	return &user, nil
}

const updateUserProfileQuery = "UPDATE users SET about = COALESCE($1, users.about), email = COALESCE($2, users.email), fullname = COALESCE($3, users.fullname) WHERE lower(nickname)=lower($4) RETURNING about, email, fullname, nickname"

func UpdateUserProfile(newData *models.UserUpd, nickname interface{}) (*models.User, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit()

	user := models.User{}

	if err = tx.QueryRow(updateUserProfileQuery, newData.About, newData.Email, newData.Fullname, nickname).
		Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
		if _, ok := err.(pgx.PgError); ok {
			return nil, dberrors.ErrUserConflict
		}
		return nil, dberrors.ErrUserNotFound
	}

	return &user, nil
}

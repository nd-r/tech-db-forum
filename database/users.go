package database

import (
	"github.com/nd-r/tech-db-forum/models"
	"log"
	"strings"
)

const createUserQuery = "INSERT INTO users (about, email, fullname, nickname) VALUES ($1, $2, $3, $4)  RETURNING about, email, fullname, nickname"
const selectUsrByNickOrEmailQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1) OR lower(email)=lower($2)"

func CreateUser(usr *models.User) (*models.UsersArr, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	userArr := models.UsersArr{}
	user := models.User{}
	tx.Select(&userArr, selectUsrByNickOrEmailQuery, usr.Nickname, usr.Email)

	if len(userArr) == 0 {
		err := tx.Get(&user, createUserQuery, usr.About, usr.Email, usr.Fullname, usr.Nickname)
		if err != nil {
			log.Fatal(err)
		}
		usr = &user
		return nil, 201
	}
	return &userArr, 409
}

const getUserProfileQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1)"

func GetUserProfile(nickname interface{}) (*models.User, error) {
	tx := db.MustBegin()
	defer tx.Commit()
	user := models.User{}

	err := tx.Get(&user, getUserProfileQuery, nickname)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, err
}

const updateUserProfileQuery = "UPDATE users SET about=$1, email=$2, fullname=$3 WHERE lower(nickname)=lower($4) RETURNING about, email, fullname, nickname"

func UpdateUserProfile(newData *models.UserUpd) (*models.User, int) {
	tx := db.MustBegin()
	defer tx.Commit()
	userArr := models.UsersArr{}

	tx.Select(&userArr, selectUsrByNickOrEmailQuery, newData.Nickname, newData.Email)

	switch len(userArr) {
	case 0:
		return nil, 404
	case 1:
		if strings.ToLower(userArr[0].Nickname) == strings.ToLower(*newData.Nickname) {
			if newData.About != nil {
				userArr[0].About = *newData.About
			}
			if newData.Email != nil {
				userArr[0].Email = *newData.Email
			}
			if newData.Fullname != nil {
				userArr[0].Fullname = *newData.Fullname
			}
			
			tx.Select(&userArr[0], updateUserProfileQuery, userArr[0].About, userArr[0].Email, userArr[0].Fullname, userArr[0].Nickname)

			return &userArr[0], 200
		}
		return nil, 404
	case 2:
		return nil, 409
	default:
		return nil, 0
	}
}

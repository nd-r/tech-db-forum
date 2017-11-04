package database

import (
	"github.com/nd-r/tech-db-forum/models"
	"log"
	"strings"
)

const createUserQuery = "INSERT INTO users (about, email, fullname, nickname) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id"
const selectUsrByNickOrEmailQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1) OR lower(email)=lower($2)"

func CreateUser(usr *models.User, nickname interface{}) (*models.UsersArr, error) {
	tx := db.MustBegin()
	defer tx.Commit()

	var id int
	err := tx.Get(&id, createUserQuery, usr.About, usr.Email, usr.Fullname, nickname)
	if err != nil{
		var users models.UsersArr
		err := tx.Select(&users, selectUsrByNickOrEmailQuery, nickname, usr.Email)
		if err != nil{
			log.Fatalln(err);
		}
		return &users, nil;
	}
	return nil, nil;
}

const getUserProfileQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1)"

func GetUserProfile(nickname interface{}) (*models.User, error) {
	tx := db.MustBegin()
	defer tx.Commit()
	user := models.User{}

	err := tx.Get(&user, getUserProfileQuery, nickname)
	if err != nil {
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

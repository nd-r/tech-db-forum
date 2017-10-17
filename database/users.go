package database

import (
	"github.com/nd-r/tech-db-forum/models"
	"strings"
)

const createUserQuery = "INSERT INTO users (about, email, fullname, nickname) VALUES ($1, $2, $3, $4)"
const selectUsrByNickOrEmailQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1) OR lower(email)=lower($2)"

func CreateUser(usr *models.User) (*models.UsersArr, int) {
	tx := db.MustBegin()
	defer tx.Commit()

	userArr := models.UsersArr{}
	tx.Select(&userArr, selectUsrByNickOrEmailQuery, usr.Nickname.V, usr.Email)

	if len(userArr) == 0 {
		tx.Queryx(createUserQuery, usr.About.V, usr.Email, usr.Fullname, usr.Nickname.V)
		return nil, 201
	}
	return &userArr, 409
}

const getUserProfileQuery = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname)=lower($1)"

func GetUserProfile(nickname []byte) (*models.UserResp, error) {
	tx := db.MustBegin()
	defer tx.Commit()
	user := models.UserResp{}

	err := tx.Get(&user, getUserProfileQuery, nickname)
	if err != nil {
		return nil, err
	}

	return &user, err
}

const updateUserProfileQuery = "UPDATE users SET about=$1, email=$2, fullname=$3 WHERE lower(nickname)=lower($4)"

func UpdateUserProfile(newData *models.UserUpdProfile) (*models.UserResp, int) {
	tx := db.MustBegin()
	defer tx.Commit()
	userArr := models.UsersArr{}

	tx.Select(&userArr, selectUsrByNickOrEmailQuery, newData.Nickname.V, newData.Email.V)

	switch len(userArr) {
	case 0:
		return nil, 404
	case 1:
		if strings.ToLower(userArr[0].Nickname) == strings.ToLower(newData.Nickname.V) {
			if newData.About.Defined {
				userArr[0].About = newData.About.V
			}
			if newData.Email.Defined {
				userArr[0].Email = newData.Email.V
			}
			if newData.Fullname.Defined {
				userArr[0].Fullname = newData.Fullname.V
			}
			tx.Queryx(updateUserProfileQuery, userArr[0].About, userArr[0].Email, userArr[0].Fullname, userArr[0].Nickname)

			return &userArr[0], 200
		}
		return nil, 404
	case 2:
		return nil, 409
	default:
		return nil, 0
	}
}

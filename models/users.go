package models

import "strings"

//easyjson:json
type User struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

func (u *User) Compare(rhs *User) int {
	return strings.Compare(strings.ToLower(u.Nickname), strings.ToLower(rhs.Nickname))
}

//easyjson:json
type UserUpd struct {
	About    *string `json:"about,omitempty"`
	Email    *string `json:"email"`
	Fullname *string `json:"fullname"`
	Nickname *string `json:"nickname,omitempty"`
}

//easyjson:json
type UsersArr []*User

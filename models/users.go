package models

import (
	"github.com/mailru/easyjson/opt"
)

type User struct {
	About    opt.String `json:"about,omitempty"`
	Email    string     `json:"email"`
	Fullname string     `json:"fullname"`
	Nickname opt.String `json:"nickname,omitempty"`
}

type UserUpdProfile struct {
	About    opt.String `json:"about,omitempty"`
	Email    opt.String `json:"email"`
	Fullname opt.String `json:"fullname"`
	Nickname opt.String `json:"nickname,omitempty"`
}

//easyjson:json
type UsersArr []UserResp

type UserResp struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

package models

import (
	"time"
)

type Forum struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Moderator string `json:"user"`
	Threads   int    `json:"threads"`
	Posts     int    `json:"posts"`
}

type Thread struct {
	Id          *int      `json:"id"`
	Slug        *string   `json:"slug"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Forum_slug  string    `json:"forum"`
	User_nick   string    `json:"author"`
	Created     time.Time `json:"created"`
	Votes_count *int      `json:"votes"`
}

type ThreadUpdate struct {
	Message *string `json:"message"`
	Title   *string `json:"title"`
}

//easyjson:json
type TreadArr []*Thread

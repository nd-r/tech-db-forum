package models

import (
	"time"
)

type Post struct {
	Id         int        `json:"id"`
	User_nick  string     `json:"author"`
	Message    string     `json:"message"`
	Created    *time.Time `json:"created"`
	Forum_slug string     `json:"forum"`
	Thread_id  int        `json:"thread"`
	Is_edited  bool       `json:"isEdited"`
	Parent     int64      `json:"parent,omitempty"`
	Parents    []int64
}

type PostDetails struct {
	AuthorDetails *User   `json:"author,omitempty"`
	ForumDetails  *Forum  `json:"forum,omitempty"`
	PostDetails   *Post   `json:"post,omitempty"`
	ThreadDetails *Thread `json:"thread,omitempty"`
}

type PostUpdate struct {
	Message *string `json:"message"`
}

//easyjson:json
type PostArr []*Post

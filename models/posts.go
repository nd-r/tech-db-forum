package models

type Post struct {
	Id         *int    `json:"id"`
	User_nick  string  `json:"author"`
	Message    string  `json:"message"`
	Created    *string `json:"created"`
	Forum_slug string  `json:"forum"`
	Thread_id  *int    `json:"thread"`
	Is_edited  *bool   `json:"isEdited"`
	Parent     *int    `json:"parent,omitempty"`
	Parents    *string
}

//easyjson:json
type PostArr []Post

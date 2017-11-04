package models

type Forum struct {
	Id        int
	Posts     *int   `json:"posts"`
	Slug      string `json:"slug"`
	Threads   *int   `json:"threads"`
	Title     string `json:"title"`
	Moderator string `json:"user"`
}

type Thread struct {
	Id          *int    `json:"id"`
	Slug        *string `json:"slug"`
	Title       string  `json:"title"`
	Message     string  `json:"message"`
	Forum_slug  string  `json:"forum"`
	User_nick   string  `json:"author"`
	Created     *string `json:"created"`
	Votes_count *int    `json:"votes"`
}

type ThreadUpdate struct {
	Message *string `json:"message"`
	Title   *string `json:"title"`
}

//easyjson:json
type TreadArr []*Thread

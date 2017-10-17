package models

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type VoteDB struct{
	ID int 
	Nickname string
	Thread_id int
	Voice int
}

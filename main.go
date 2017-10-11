package main

import (
	services "./services"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type User struct {
	Nickname string
	Email    string
	Fullname string
	About    string
}


func main() {
	db, err := services.DBPoolInit()

	if err != nil {
		log.Fatalf("%s", err)
	}
	services.InitDBSchema(db)

	router := services.RouterInit()
	log.Println(http.ListenAndServe(":8000", router))
}

package services

import (
	// "github.com/nd-r/tech-db-forum/handlers"
	"../handlers"
	"github.com/julienschmidt/httprouter"
	"log"
)

// RouterInit inits the router and routes
func RouterInit() *httprouter.Router {
	r := httprouter.New()

	db, err := DBPoolInit()

	if err != nil {
		log.Fatalf("%s", err)
	}

	InitDBSchema(db)
		
	r.GET("/api/forum/:slug/:param", handlers.GETForumHandler)
	// r.POST("", )
	return r
}

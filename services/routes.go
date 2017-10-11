package services

import (
	handlers "../handlers"
	"github.com/julienschmidt/httprouter"
)

// RouterInit inits the router and routes
func RouterInit() *httprouter.Router {
	r := httprouter.New()

	// r.POST("/forum/*param", handlers.CreateForum)
	r.GET("/forum/:slug/*param", handlers.CreateForum)
	// r.GET("/forum/:slug/:param", handlers.GETForumHandler)
	return r
}

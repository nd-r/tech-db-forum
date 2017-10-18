package services

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/nd-r/tech-db-forum/handlers"
)

// RouterInit inits the router and routes
func RouterInit() *fasthttprouter.Router {
	r := fasthttprouter.New()

	//slug == create
	//yes
	r.POST("/api/forum/:slug", handlers.CreateForum)

	//yes
	r.POST("/api/forum/:slug/create", handlers.CreateThread)
	//yes
	r.GET("/api/forum/:slug/details", handlers.GetForumDetails)
	//yes
	r.GET("/api/forum/:slug/threads", handlers.GetForumThreads)
	r.GET("/api/forum/:slug/users", handlers.GetForumUsers)

	r.GET("/api/post/:id/details", handlers.GetPostDetails)
	r.POST("/api/post/:id/details", handlers.ChangePostDetails)

	r.POST("/api/service/clear", handlers.DBClear)
	r.GET("/api/service/status", handlers.DBStatus)

	//yes
	r.POST("/api/thread/:slug_or_id/create", handlers.CreateNewPosts)
	//yes
	r.GET("/api/thread/:slug_or_id/details", handlers.GetThreadDetails)
	//yes
	r.POST("/api/thread/:slug_or_id/details", handlers.UpdateThreadDetails)
	//yes
	r.GET("/api/thread/:slug_or_id/posts", handlers.GetThreadPosts)
	//yes
	r.POST("/api/thread/:slug_or_id/vote", handlers.VoteThread)

	//yes
	r.POST("/api/user/:nickname/create", handlers.CreateUser)
	//yes
	r.GET("/api/user/:nickname/profile", handlers.GetUserProfile)
	//yes
	r.POST("/api/user/:nickname/profile", handlers.UpdateUserProfile)

	return r
}

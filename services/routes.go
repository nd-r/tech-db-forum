package services

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/nd-r/tech-db-forum/handlers"
)

// RouterInit inits the router and routes
func RouterInit() *fasthttprouter.Router {
	r := fasthttprouter.New()

	//slug == create
	r.POST("/api/forum/:slug", handlers.CreateForum)

	r.POST("/api/forum/:slug/create", handlers.CreateThread)
	r.GET("/api/forum/:slug/details", handlers.GetForumDetails)
	r.GET("/api/forum/:slug/threads", handlers.GetForumThreads)
	r.GET("/api/forum/:slug/users", handlers.GetForumUsers)

	r.GET("/api/post/:id/details", handlers.GetPostDetails)
	r.POST("/api/post/:id/details", handlers.ChangePostDetails)

	r.POST("/api/service/clear", handlers.DBClear)
	r.GET("/api/service/status", handlers.DBStatus)

	r.POST("/api/thread/:slug_or_id/create", handlers.CreateThread)
	r.GET("/api/thread/:slug_or_id/details", handlers.GetThreadDetails)
	r.POST("/api/thread/:slug_or_id/details", handlers.UpdateThreadDetails)
	r.GET("/api/thread/:slug_or_id/posts", handlers.GetThreadPosts)
	r.POST("/api/thread/:slug_or_id/vote", handlers.VoteThread)

	r.POST("/api/user/:nickname/create", handlers.CreateUser)
	r.GET("/api/user/:nickname/profile", handlers.GetUserProfile)
	r.POST("/api/user/:nickname/profile", handlers.UpdateUserProfile)

	return r
}

package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	"log"
)

// CreateForum - handler создания форума
// POST /forum/create
func CreateForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	forum := models.Forum{}

	forum.UnmarshalJSON(ctx.PostBody())

	existingForum, statusCode := database.CreateForum(&forum)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 201:
		resp, _ := forum.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	case 409:
		resp, _ := existingForum.MarshalJSON()
		ctx.Write(resp)
	}
}

// CreateThread - handler создания ветки
// POST /forum/{slug}/create
func CreateThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	thread := models.Thread{}
	thread.UnmarshalJSON(ctx.PostBody())

	slug := ctx.UserValue("slug")
	if slug != nil {
		thread.Forum_title = ctx.UserValue("slug").(string)
	} else {
		log.Panicln(string(ctx.Path()))
	}

	threadExisting, statusCode := database.CreateThread(&thread)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 201:
		resp, _ := thread.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	case 409:
		resp, _ := threadExisting.MarshalJSON()
		ctx.Write(resp)
	}
}

// GetForumDetails - handler информации ветки
// GET /forum/{slug}/details
func GetForumDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slug").(string)

	forum, statusCode := database.GetForumDetails(slug)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		resp, _ := forum.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}
}

// GetForumThreads - handler получения списка ветвей данного форума
// GET /forum/{slug}/threads
func GetForumThreads(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")

	threadArr, statusCode := database.GetForumThreads(&slug, limit, since, desc)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		if threadArr == nil {
			ctx.Write([]byte("[]"))
			return
		}
		resp, _ := threadArr.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}
}

// GetForumUsers - handler получение списка пользователей
// GET /forum/{slug}/users
func GetForumUsers(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	// slug := ctx.UserValue("slug").(string)

	// users, statusCode := database.GetForumUsers(&slug)
	// ctx.SetStatusCode(statusCode)

	// switch statusCode {
	// case 200:
	// 	if len(users) != 0 {
	// 		resp, _ := users.MarshalJSON()
	// 		ctx.Write(resp)
	// 		return
	// 	}
	// 	ctx.Write([]byte("[]"))
	// case 404:
	// 	ctx.Write(models.ErrorMsg)
	// }
}

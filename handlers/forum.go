package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	// "log"
)

// CreateForum - handler создания форума
// POST /forum/create
func CreateForum(ctx *fasthttp.RequestCtx) {
	forum := models.Forum{}
	forum.UnmarshalJSON(ctx.PostBody())

	forumResp, err := database.CreateForum(&forum)

	var resp []byte

	switch err {
	case nil:
		ctx.SetStatusCode(201)
		resp, _ = forumResp.MarshalJSON()

	case dberrors.ErrForumExists:
		ctx.SetStatusCode(409)
		resp, _ = forumResp.MarshalJSON()

	case dberrors.ErrUserNotFound:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

// CreateThread - handler создания ветки
// POST /forum/{slug}/create
func CreateThread(ctx *fasthttp.RequestCtx) {
	threadDetails := models.Thread{}
	threadDetails.UnmarshalJSON(ctx.PostBody())

	slug := ctx.UserValue("slug")

	threadExisting, err := database.CreateThread(&slug, &threadDetails)

	var resp []byte

	switch err {
	case nil:
		ctx.SetStatusCode(201)
		resp, _ = threadExisting.MarshalJSON()

	case dberrors.ErrUserNotFound, dberrors.ErrForumNotFound:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg

	case dberrors.ErrThreadExists:
		ctx.SetStatusCode(409)
		resp, _ = threadExisting.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

// GetForumDetails - handler информации ветки
// GET /forum/{slug}/details
func GetForumDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug")

	forum, err := database.GetForumDetails(slug)

	var resp []byte

	switch err {
	case nil:
		resp, _ = forum.MarshalJSON()
	default:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

// GetForumThreads - handler получения списка ветвей данного форума
// GET /forum/{slug}/threads
func GetForumThreads(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	

	slug := ctx.UserValue("slug")
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")

	threadArr, error := database.GetForumThreads(&slug, limit, since, desc)

	var resp []byte
	switch error {
	case nil:
		if threadArr == nil {
			ctx.Write([]byte("[]"))
			return
		}
		resp, _ = threadArr.MarshalJSON()
	case dberrors.ErrForumNotFound:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg
	}

	ctx.Write(resp)
}

// GetForumUsers - handler получение списка пользователей
// GET /forum/{slug}/users
func GetForumUsers(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")

	users, statusCode := database.GetForumUsers(&slug, limit, since, desc)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		if len(*users) != 0 {
			resp, _ := users.MarshalJSON()
			ctx.Write(resp)
			return
		}
		ctx.Write([]byte("[]"))
	case 404:
		ctx.Write(models.ErrorMsg)
	}
}

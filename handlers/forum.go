package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	// "log"
	"github.com/nd-r/tech-db-forum/cache"
	"strings"
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
	var resp []byte

	cachedForum := cache.TheGreatForumCache.Get(strings.ToLower(slug.(string)))

	if cachedForum == nil {
		forum, err := database.GetForumDetails(slug)

		switch err {
		case nil:
			resp, _ = forum.MarshalJSON()
			cachedForum := cache.CachedForum{Model: *forum, Json: resp}
			cache.TheGreatForumCache.Push(strings.ToLower(slug.(string)), &cachedForum)
		default:
			ctx.SetStatusCode(404)
			resp = models.ErrorMsg
		}
	} else {
		resp = cachedForum.Json
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
		if len(*threadArr) == 0 {
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

	slug := ctx.UserValue("slug")
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	desc := ctx.QueryArgs().Peek("desc")

	users, err := database.GetForumUsers(&slug, limit, since, desc)

	var resp []byte

	switch err {
	case nil:
		ctx.SetStatusCode(200)
		if len(*users) != 0 {
			resp, _ = users.MarshalJSON()
		} else {
			resp = []byte("[]")
		}
	case dberrors.ErrForumNotFound:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

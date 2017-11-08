package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	"log"
)

// CreateNewPosts - создание новых постов
// POST /thread/{slug_or_id}/create
func CreateNewPosts(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id")

	postsArr := models.PostArr{}
	postsArr.UnmarshalJSON(ctx.PostBody())

	newPosts, err := database.CreatePosts(slugOrID, &postsArr)

	var resp []byte

	switch err {
	case nil:
		ctx.SetStatusCode(201)
		if newPosts != nil {
			resp, _ = newPosts.MarshalJSON()
		} else {
			resp = []byte("[]")
		}

	case dberrors.ErrThreadNotFound, dberrors.ErrUserNotFound:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg

	case dberrors.ErrPostsConflict:
		ctx.SetStatusCode(409)
		resp = models.ErrorMsg
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

// GetThreadDetails - получение информации о ветке обсуждения
// GET /thread/{slug_or_id}/details
func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slugOrID := ctx.UserValue("slug_or_id")

	thread, err := database.GetThread(slugOrID)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(models.ErrorMsg)
		return
	}

	var resp []byte
	resp, err = thread.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
	}

	ctx.Write(resp)
}

// UpdateThreadDetails - обновление ветки
// POST /thread/{slug_or_id}/details
func UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	thrUpdate := models.ThreadUpdate{}
	thrUpdate.UnmarshalJSON(ctx.PostBody())

	thread, statusCode := database.UpdateThreadDetails(&slugOrID, &thrUpdate)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		resp, _ := thread.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}
}

// GetThreadPosts - сообщения данной ветки
// GET /thread/{slug_or_id}/posts
func GetThreadPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slugOrID := ctx.UserValue("slug_or_id").(string)
	limit := ctx.QueryArgs().Peek("limit")
	since := ctx.QueryArgs().Peek("since")
	sort := ctx.QueryArgs().Peek("sort")
	desc := ctx.QueryArgs().Peek("desc")

	postArr, statusCode := database.GetThreadPosts(&slugOrID, limit, since, sort, desc)

	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		if len(*postArr) != 0 {
			resp, _ := postArr.MarshalJSON()
			ctx.Write(resp)
		} else {
			ctx.Write([]byte("[]"))
		}
	case 404:
		ctx.Write(models.ErrorMsg)
	}
}

// VoteThread - проголосовать за ветвь обсуждения
// POST /thread/{slug_or_id}/vote
func VoteThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	var vote models.Vote
	vote.UnmarshalJSON(ctx.PostBody())

	slugOrID := ctx.UserValue("slug_or_id")

	thread, err := database.PutVote(slugOrID, &vote)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(models.ErrorMsg)
		return
	}

	ctx.SetStatusCode(200)
	resp, _ := thread.MarshalJSON()
	ctx.Write(resp)
}

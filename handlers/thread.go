package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	"log"
)

// CreateNewPosts - создание новых постов
// POST /thread/{slug_or_id}/create
func CreateNewPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	if len(ctx.PostBody()) == 3 {
		ctx.SetStatusCode(201)
		ctx.Write([]byte("[]"))
		return
	}

	postsArr := models.PostArr{}
	postsArr.UnmarshalJSON(ctx.PostBody())

	newPosts, statusCode := database.CreatePosts(slugOrID, &postsArr)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 201:
		resp, _ := newPosts.MarshalJSON()
		ctx.Write(resp)
	default:
		ctx.Write(models.ErrorMsg)
	}
}

// GetThreadDetails - получение информации о ветке обсуждения
// GET /thread/{slug_or_id}/details
func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slugOrID := ctx.UserValue("slug_or_id").(string)

	thread, statusCode := database.GetThread(&slugOrID)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		resp, _ := thread.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}
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
			resp, err := postArr.MarshalJSON()
			if err != nil {
				log.Println(err)
			}
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
	vote := models.Vote{}
	vote.UnmarshalJSON(ctx.PostBody())

	slugOrID := ctx.UserValue("slug_or_id").(string)

	thread, statusCode := database.PutVote(slugOrID, &vote)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		resp, _ := thread.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}
}

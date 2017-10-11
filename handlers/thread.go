package handlers

import (
	"github.com/valyala/fasthttp"
)

// CreateNewPosts - создание новых постов
// POST /thread/{slug_or_id}/create
func CreateNewPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

// GetThreadDetails - получение информации о ветке обсуждения
// GET /thread/{slug_or_id}/details
func GetThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

// UpdateThreadDetails - обновление ветки
// POST /thread/{slug_or_id}/details
func UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

// GetThreadPosts - сообщения данной ветки
// GET /thread/{slug_or_id}/posts
func GetThreadPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

// VoteThread - проголосовать за ветвь обсуждения
// POST /thread/{slug_or_id}/vote
func VoteThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

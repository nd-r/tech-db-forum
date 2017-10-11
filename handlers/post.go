package handlers

import (
	"github.com/valyala/fasthttp"
)

// GetPostDetails получение информации о ветке обсуждения
// GET /post/{id}/details
func GetPostDetails (ctx *fasthttp.RequestCtx){
	ctx.SetContentType("application/json")
	
}

// ChangePostDetails изменение сообщения
// POST /post/{id}/details
func ChangePostDetails (ctx *fasthttp.RequestCtx){
	ctx.SetContentType("application/json")
	
}


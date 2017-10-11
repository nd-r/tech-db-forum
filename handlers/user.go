package handlers

import (
	"github.com/valyala/fasthttp"
)

// CreateUser - создание нового пользователя
// POST /user/{nickname}/create
func CreateUser(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

// GetUserProfile - получение информации о пользователе
// GET /user/{nickname}/profile
func GetUserProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

// UpdateUserProfile - изменение информации о пользователе
// POST /user/{nickname}/profile
func UpdateUserProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

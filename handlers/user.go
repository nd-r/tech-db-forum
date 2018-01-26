package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	"github.com/nd-r/tech-db-forum/cache"
	"strings"
	"log"
)

// CreateUser - создание нового пользователя
// POST /user/{nickname}/create
func CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User
	user.UnmarshalJSON(ctx.PostBody())

	nickname := ctx.UserValue("nickname")

	existingUsers, err := database.CreateUser(&user, nickname)

	var resp []byte

	switch err {
	case nil:
		ctx.SetStatusCode(201)
		user.Nickname = ctx.UserValue("nickname").(string)
		resp, _ = user.MarshalJSON()

	case dberrors.ErrUserExists:
		ctx.SetStatusCode(409)
		resp, _ = existingUsers.MarshalJSON()
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

// GetUserProfile - получение информации о пользователе
// GET /user/{nickname}/profile
func GetUserProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	nickname := ctx.UserValue("nickname")
	user := cache.TheGreatUserCache.Get(strings.ToLower(nickname.(string)))

	if user != nil {
		ctx.Write(user)
		return
	}

	var resp []byte
	var err error

	userFromDb, err := database.GetUserProfile(nickname)

	switch err {
	case nil:
		resp, _ = userFromDb.MarshalJSON()
		cache.TheGreatUserCache.Push(strings.ToLower(userFromDb.Nickname), &resp)
	case dberrors.ErrUserNotFound:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg
	}

	ctx.Write(resp)
}

// UpdateUserProfile - изменение информации о пользователе
// POST /user/{nickname}/profile
func UpdateUserProfile(ctx *fasthttp.RequestCtx) {
	user := models.UserUpd{}
	user.UnmarshalJSON(ctx.PostBody())

	nickname := ctx.UserValue("nickname")

	userUpdated, error := database.UpdateUserProfile(&user, &nickname)
	var resp []byte

	switch error {
	case nil:
		resp, _ = userUpdated.MarshalJSON()


	case dberrors.ErrUserConflict:
		ctx.SetStatusCode(409)
		resp = models.ErrorMsg
	case dberrors.ErrUserNotFound:
		ctx.SetStatusCode(404)
		resp = models.ErrorMsg
	}

	ctx.SetContentType("application/json")
	ctx.Write(resp)
}

package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
)

// CreateUser - создание нового пользователя
// POST /user/{nickname}/create
func CreateUser(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	var user models.User
	user.UnmarshalJSON(ctx.PostBody())

	user.Nickname = ctx.UserValue("nickname").(string)

	userArr, statusCode := database.CreateUser(&user)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 201:
		resp, _ := user.MarshalJSON()
		ctx.Write(resp)
	case 409:
		resp, _ := userArr.MarshalJSON()
		ctx.Write(resp)
	}
}

// GetUserProfile - получение информации о пользователе
// GET /user/{nickname}/profile
func GetUserProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	nickname := []byte(ctx.UserValue("nickname").(string))

	user, err := database.GetUserProfile(nickname)

	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(models.ErrorMsg)
		return
	}

	var resp []byte
	resp, err = user.MarshalJSON()
	ctx.Write(resp)
}

// UpdateUserProfile - изменение информации о пользователе
// POST /user/{nickname}/profile
func UpdateUserProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := models.UserUpd{}
	user.UnmarshalJSON(ctx.PostBody())
	
	nickname := ctx.UserValue("nickname").(string)
	user.Nickname = &nickname

	userUpdated, statusCode := database.UpdateUserProfile(&user)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		resp, _ := userUpdated.MarshalJSON()
		ctx.Write(resp)
	default:
		ctx.Write(models.ErrorMsg)
	}
}

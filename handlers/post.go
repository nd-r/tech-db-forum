package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
)

// GetPostDetails получение информации о ветке обсуждения
// GET /post/{id}/details
func GetPostDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id := ctx.UserValue("id").(string)
	related := ctx.QueryArgs().Peek("related")

	postDetails, statusCode := database.GetPostDetails(&id, related)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		resp, _ := postDetails.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}

}

// ChangePostDetails изменение сообщения
// POST /post/{id}/details
func ChangePostDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	id := ctx.UserValue("id").(string)
	postUpd := models.PostUpdate{}

	postUpd.UnmarshalJSON(ctx.PostBody())

	post, statusCode := database.UpdatePostDetails(&id, &postUpd)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		resp, _ := post.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}
}

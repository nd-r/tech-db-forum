package handlers

import (
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/dberrors"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	"log"
)

// CreateForum - handler создания форума
// POST /forum/create
func CreateForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	forum := models.Forum{}

	forum.UnmarshalJSON(ctx.PostBody())

	pqErr := database.CreateForum(&forum)
	if pqErr != nil {
		switch pqErr.Code {
		case dberrors.UniqueConstraint:
			ctx.SetStatusCode(409)

			forumExisting, _ := database.GetForumDetails(forum.Slug)
			resp, _ := forumExisting.MarshalJSON()
			ctx.Write(resp)

			return

		case dberrors.NotNullConstraint:
			ctx.SetStatusCode(404)
			ctx.Write(models.ErrorMsg)

			return
		}
	}

	ctx.SetStatusCode(201)
	resp, _ := forum.MarshalJSON()
	ctx.Write(resp)
}

// CreateThread - handler создания ветки
// POST /forum/{slug}/create
func CreateThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	threadDetails := models.Thread{}
	threadDetails.UnmarshalJSON(ctx.PostBody())

	slug := ctx.UserValue("slug")

	pqErr := database.CreateThread(&slug, &threadDetails)
	if pqErr != nil {
		switch pqErr.Code {
		case dberrors.NotNullConstraint:
			ctx.SetStatusCode(404)
			ctx.Write(models.ErrorMsg)

			return

		case dberrors.UniqueConstraint:
			ctx.SetStatusCode(409)
			threadExisting, err := database.GetThread(*threadDetails.Slug)
			if err != nil{
				log.Fatalln(err)
			}
			
			resp, err := threadExisting.MarshalJSON()
			if err != nil{
				log.Fatalln(err)
			}
			ctx.Write(resp)

			return
		}
	}

	ctx.SetStatusCode(201)
	resp, err := threadDetails.MarshalJSON()
	if err != nil {
		log.Fatalln(err)
	}
	ctx.Write(resp)
}

// GetForumDetails - handler информации ветки
// GET /forum/{slug}/details
func GetForumDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slug")

	forum, err := database.GetForumDetails(slug)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write(models.ErrorMsg)

		return
	}

	ctx.SetStatusCode(200)
	resp, _ := forum.MarshalJSON()
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

	threadArr, statusCode := database.GetForumThreads(&slug, limit, since, desc)
	ctx.SetStatusCode(statusCode)

	switch statusCode {
	case 200:
		if threadArr == nil {
			ctx.Write([]byte("[]"))
			return
		}
		resp, _ := threadArr.MarshalJSON()
		ctx.Write(resp)
	case 404:
		ctx.Write(models.ErrorMsg)
	}
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

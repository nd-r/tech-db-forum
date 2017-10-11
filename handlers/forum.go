package handlers

import (
	"fmt"
	"github.com/nd-r/tech-db-forum/database"
	"github.com/valyala/fasthttp"
)

// CreateForum - handler создания форума
// POST /forum/create
func CreateForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.Write([]byte("ashdfjfa"))
}

// CreateThread - handler создания ветки
// POST /forum/{slug}/create
func CreateThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	// slug := new string(ctx.UserValue("params"))

}

// GetForumDetails - handler информации ветки
// GET /forum/{slug}/details
func GetForumDetails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	database.DB.MustExec("INSERT into users VALUES ('asdasdf', 'asdasdf', 'asdasdf','asdasdf')")
	ctx.Write([]byte("ashdfjfa"))
}

// GetForumThreads - handler получения списка ветвей данного форума
// GET /forum/{slug}/threads
func GetForumThreads(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	database.DB.MustExec("INSERT into vote VALUES (DEFAULT, 'asdasdf', 1)")
}

// GetForumUsers - handler получение списка пользователей
// GET /forum/{slug}/users
func GetForumUsers(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	fmt.Fprintf(ctx, "%s", "lalala")
}

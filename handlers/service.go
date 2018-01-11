package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/cache"
)


// DBClear - очистка всех данных в базе
// POST /service/clear
func DBClear(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	database.DeleteDB()
	cache.TheGreatForumCache.Clear()
	cache.TheGreatUserCache.Clear()
}

// DBStatus - получение информации о базе данных
// GET /service/status
func DBStatus(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	status := database.GetDBStatus()
	resp, _ := status.MarshalJSON()
	ctx.Write(resp)
}

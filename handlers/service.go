package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/nd-r/tech-db-forum/database"
)


// DBClear - очистка всех данных в базе
// POST /service/clear
func DBClear(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	database.DeleteDB()
}

// DBStatus - получение информации о базе данных
// GET /service/status
func DBStatus(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	status := database.GetDBStatus()
	resp, _ := status.MarshalJSON()
	ctx.Write(resp)
}

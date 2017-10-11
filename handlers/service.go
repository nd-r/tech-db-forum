package handlers

import (
	"github.com/valyala/fasthttp"
)

// DBClear - очистка всех данных в базе
// POST /service/clear
func DBClear(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

// DBStatus - получение информации о базе данных
// GET /service/status
func DBStatus(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

}

package main

import (
	_ "github.com/lib/pq"
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/services"
	"github.com/valyala/fasthttp"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(8)

	database.DBPoolInit()
	database.InitDBSchema()

	router := services.RouterInit()
	log.Println(fasthttp.ListenAndServe(":8000", router.Handler))
}

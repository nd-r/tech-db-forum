package main

import (
	_ "github.com/lib/pq"
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/services"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	"log"
	// "runtime"
)

func main() {
	// runtime.GOMAXPROCS(2)

	database.DBPoolInit()
	database.InitDBSchema()

	router := services.RouterInit()

	error := models.ErrorStr{Message: "error occured"}
	models.ErrorMsg, _ = error.MarshalJSON()
	log.SetFlags(log.Llongfile)
	log.Println("started")
	log.Println(fasthttp.ListenAndServe(":5000", router.Handler))
}

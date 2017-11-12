package main

import (
	// "sync"
	_ "github.com/lib/pq"
	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/services"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/valyala/fasthttp"
	// "net/http"
	"log"
	// _ "net/http/pprof"
	// _ "runtime"
)

func main() {
	// runtime.GOMAXPROCS(8)

	database.DBPoolInit()
	database.InitDBSchema()

	router := services.RouterInit()

	error := models.ErrorStr{Message: "error occured"}
	models.ErrorMsg, _ = error.MarshalJSON()
	log.SetFlags(log.Llongfile)
	log.Println("started")
	// go http.ListenAndServe(":1111",nil)
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go fasthttp.ListenAndServe(":5000", router.Handler)
	//  wg.Wait()
	fasthttp.ListenAndServe(":5000", router.Handler)
}

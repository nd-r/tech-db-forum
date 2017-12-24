package main

import (
	"runtime"
	//"sync"

	"github.com/nd-r/tech-db-forum/database"
	"github.com/nd-r/tech-db-forum/models"
	"github.com/nd-r/tech-db-forum/services"
	"github.com/valyala/fasthttp"

	"log"
	//"net/http"
	//_ "net/http/pprof"
)

func main() {
	runtime.GOMAXPROCS(8)
	log.SetFlags(log.Llongfile)

	database.DBPoolInit()

	router := services.RouterInit()

	err := models.ErrorStr{Message: "error occured"}
	models.ErrorMsg, _ = err.MarshalJSON()
	log.Println("started")
	//go http.ListenAndServe(":1111", nil)
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go fasthttp.ListenAndServe(":5000", router.Handler)
	//wg.Wait()
	fasthttp.ListenAndServe(":5000", router.Handler)
}

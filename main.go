package main

import (
	"runtime"
	services "./services"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	runtime.GOMAXPROCS(2)

	router := services.RouterInit()
	log.Println(http.ListenAndServe(":8000", router))
}

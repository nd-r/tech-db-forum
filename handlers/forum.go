package handlers

import (
	// "github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"net/http"
	// "../services"
)

// var db *sqlx.DB

// CreateForum - handler создания форума
// /forum/create
func CreateForum(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("kjasdjfa"))
	// fmt.Fprintf(w, "%s", "/forum/create")
}

// createThread - handler создания ветки
// /forum/{slug}/create
func createThread(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// fmt.Fprintf(w, "%s", "/forum/{slug}/create")
}

// getThreadDetails - handler информации ветки
// /forum/{slug}/details
func getThreadDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("kjasdjfa"))
}

// getForumThreads - handler получения списка ветвей данного форума
// /forum/{slug}/threads
func getForumThreads(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// fmt.Fprintf(w, "%s", "/forum/{slug}/threads")
}

// getForumUsers - handler получение списка пользователей
// /forum/{slug}/users
func getForumUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// fmt.Fprintf(w, "%s", "/forum/{slug}/users")
}

// GETForumHandler - handler функций форума
// /forum/*
func GETForumHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("param") {
	case "details":
		getThreadDetails(w, r, ps)
	case "threads":
		getForumThreads(w, r, ps)
	case "users":
		getForumUsers(w, r, ps)
	}
}

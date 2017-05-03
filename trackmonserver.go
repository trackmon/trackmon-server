package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	log.Println("TRACKMON SERVER UPSTART")
	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler)

	// Make new server
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:80",
	}
	srv.ListenAndServe()
}

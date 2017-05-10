package main

import (
	"database/sql"
	"fmt"
	//"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"serverversion\":\"%s\",\"apiversion\":\"%s\"}", ServerVersion, APIVersion)
}

func UserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	HashedPassword, err := HashPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} // NOTE: EVERYTHING BELOW HERE HAS AUTHENTICATION, BUT IS NOT ACCESS CHECKED!
	// TODO: Check if user exist in database
	// TODO: IF NOT Create new user and write to database
	log.Println(username, HashedPassword, test)
}

func AllAccountHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	HashedPassword, err := HashPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} // NOTE: EVERYTHING BELOW HERE HAS AUTHENTICATION, BUT IS NOT ACCESS CHECKED!
	log.Println(username, HashedPassword, test)
}

func HistoryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	HashedPassword, err := HashPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} // NOTE: EVERYTHING BELOW HERE HAS AUTHENTICATION, BUT IS NOT ACCESS CHECKED!
	log.Println(username, HashedPassword, test)
}

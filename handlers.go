package main

import (
	"database/sql"
	"fmt"
	//"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
	var DoesExist int // sigh... why no bool?
	PrepDoesUserExistQuery.QueryRow(username).Scan(&DoesExist)
	if DoesExist == 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var _ = password
	// TODO: Check if user exist in database
	// TODO: IF NOT Create new user and write to database
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
	}
	if AuthCheck(username, HashedPassword, db) != true {
		w.WriteHeader(http.StatusForbidden)
		return
	} // ALL USERS BELOW HERE ARE AUTHENTICATED
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
	}
	if AuthCheck(username, HashedPassword, db) != true {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}

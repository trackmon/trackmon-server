package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"serverversion\":\"%s\",\"apiversion\":\"%s\"}", ServerVersion, APIVersion)
	w.WriteHeader(http.StatusOK)
}

func NewUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	log.Println(username, HashedPassword)
	// TODO: Check if user exist in database
	// TODO: IF NOT Create new user and write to database
}

func UserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	variables := mux.Vars(r)
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} // All requests below here have given basic auth
	log.Printf("User %s with pw %s wants info about all accounts of %s\n", username, password, string(variables["user_id"]))
	if username != string(variables["user_id"]) {
		w.WriteHeader(http.StatusForbidden)
	}
	// TODO: Check if user exists, if and if password correct, give him info
}

func AccountHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	variables := mux.Vars(r)
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} // All requests below here have given basic auth
	log.Printf("User %s with pw %s wants info about account %s from %s\n", username, password, string(variables["account"]), string(variables["user_id"]))
}

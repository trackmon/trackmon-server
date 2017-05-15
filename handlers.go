package main

import (
	"database/sql"
	"fmt"
	//"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	return
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
	switch r.Method {
	case "POST":
		var DoesExist int // sigh... why no bool?
		PrepDoesUserExistQuery.QueryRow(username).Scan(&DoesExist)
		if DoesExist == 1 { // 1 == true
			w.WriteHeader(http.StatusForbidden)
			return
		}

		HashedPassword, err := HashPassword(password)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		SignupTime := time.Now().Format(time.RFC3339)
		_, err = PrepAddNewUser.Exec(username, HashedPassword, SignupTime)
		if err != nil {
			log.Println("UserHandler:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "DELETE":
		if AuthCheck(username, password, db) != true {
			w.WriteHeader(http.StatusForbidden)
			return
		} // ALL USERS BELOW HERE ARE AUTHENTICATED

		_, err := PrepDeleteExistingUser.Exec(username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = PrepDeleteAccountsFromExistingUser.Exec(username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = PrepDeleteHistoryFromExistingUser.Exec(username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func AllAccountHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if AuthCheck(username, password, db) != true {
		w.WriteHeader(http.StatusForbidden)
		return
	} // ALL USERS BELOW HERE ARE AUTHENTICATED
	switch r.Method {
	case "POST":
		var NewAccount NewAccount
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = fromjson(string(body), &NewAccount)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var NewAccountID int
		err = PrepAddNewAccount.QueryRow(username, NewAccount.Accountname, NewAccount.Currency, NewAccount.Initialamount).Scan(&NewAccountID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		HistoryTime := time.Now().Format(time.RFC3339)
		_, err = PrepAddNewHistoryObject.Exec(NewAccountID, string("Initial amount"), HistoryTime, NewAccount.Initialamount)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "PUT":
		// Update all accounts (even if only one is updated)
	}
}

func HistoryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if AuthCheck(username, password, db) != true {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}

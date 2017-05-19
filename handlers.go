package main

import (
	"database/sql"
	"fmt"
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
		fmt.Fprintf(w, "{\"newaccountid\":%d}", NewAccountID)

	case "PUT":
		// Update all accounts (even if only one is updated)
		// TODO: Think about implementation
		w.WriteHeader(http.StatusNotImplemented)
		return
	case "GET":
		// TODO: Get all accounts from DB
		// TODO: Make json out of it
		// TODO: send it
		w.WriteHeader(http.StatusNotImplemented)
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
	switch r.Method {
	case "PUT":
		w.WriteHeader(http.StatusNotImplemented)
		// TODO: TASK1: Load all database elements WHERE username == username, id == id ASYNC, parse all of that into a hash map with k = id; v = HistoryElement
		// TODO: TASK1: Parse all json into another hash map with above specs
		// TODO: TASK2: Compare HM1 and HM2, changed elements get into HM3 with above specs
		// TODO: TASK2: Compare HM1 and HM2, deleted elements get into HM4 with above specs
		// TODO: TASK3: Delete all elements on the db from HM4, change everything on db from HM3
	case "POST":
		w.WriteHeader(http.StatusNotImplemented)
		// TODO: Parse json
		// TODO: Write Json to DB
	case "GET":
		w.WriteHeader(http.StatusNotImplemented)
		// TODO: get all history from database
		// TODO: make json array out of it
		// TODO: send it
	}
}

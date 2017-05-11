package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func AuthCheck(username string, passwordhash string, db *sql.DB) bool {
	var dbpasswordreturn string
	err := PrepGetUserQuery.QueryRow(username).Scan(&dbpasswordreturn)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Println(err)
			return false
		}
	}
	return CheckPasswordHash(passwordhash, dbpasswordreturn)
}

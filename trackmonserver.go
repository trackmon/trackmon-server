package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	test string = "Hello World!"

	ServerVersion string = "pre version"
	APIVersion    string = "pre api"

	DatabaseSetupUsersTable    string = "CREATE TABLE IF NOT EXISTS users (username varchar(255) PRIMARY KEY NOT NULL, passwordhash varchar(64) NOT NULL, joineddate TIMESTAMP, userid SERIAL)"
	DatabaseSetupAccountsTable string = "CREATE TABLE IF NOT EXISTS accounts (accountid SERIAL PRIMARY KEY NOT NULL, username varchar(255) REFERENCES users(username), currency varchar(3) NOT NULL, balance INT)"
	DatabaseSetupHistoryTable  string = "CREATE TABLE IF NOT EXISTS history (accountid SERIAL REFERENCES accounts(accountid), name varchar(255) NOT NULL, time TIMESTAMP NOT NULL, amount INT NOT NULL, historyid SERIAL NOT NULL PRIMARY KEY)"

	GetUserQuery               string = "SELECT passwordhash FROM users WHERE username = $1"
	DoesUserExistQuery         string = "SELECT count(1) FROM users WHERE username = $1"
	AddNewUser                 string = "INSERT INTO users (username, passwordhash, joineddate) VALUES ($1, $2, $3)"
	DeleteExistingUser string = "DELETE FROM users WHERE username = $1"
	DeleteAccountsFromExistingUser string = "DELETE FROM accounts WHERE username = $1"
	DeleteHistoryFromExistingUser string = "DELETE FROM users WHERE username = $1"
)

var (
	PrepGetUserQuery       *sql.Stmt
	PrepDoesUserExistQuery *sql.Stmt
	PrepAddNewUser         *sql.Stmt
	PrepDeleteExistingUser *sql.Stmt
	PrepDeleteAccountsFromExistingUser *sql.Stmt
	PrepDeleteHistoryFromExistingUser *sql.Stmt
)

func main() {
	fmt.Println("TRACKMON SERVER\nCopyright (c) 2017, Paul Kramme under BSD 2-Clause")
	fmt.Println("Please report bugs to https://github.com/trackmon/trackmon-server")

	// Configure flags
	CreateConfigFlag := flag.Bool("createconfig", false, "Creates a standard configuration and exits")
	ConfigLocation := flag.String("config", "./trackmonserv.conf", "Location of config file. Standard is ./trackmonserv")
	ShowLicenses := flag.Bool("licenses", false, "Shows licenses and exits")
	ShowVersion := flag.Bool("version", false, "Shows version and exits")
	ShowJsonVersion := flag.Bool("versionjson", false, "Shows version in json and exits")

	// Check flags
	flag.Parse()
	if *CreateConfigFlag == true {
		CreateConfig()
		return
	}
	if *ShowLicenses == true {
		fmt.Println("trackmon servers license\n")
		fmt.Print(trackmonlicense)
		fmt.Println("\n")

		fmt.Println("This project uses github.com/gorilla/mux\n")
		fmt.Print(muxlicense)
		fmt.Println("\n")

		return
	}
	if *ShowVersion == true {
		fmt.Printf("Server Version: %s\nAPI Version: %s\n", ServerVersion, APIVersion)
		return
	}
	if *ShowJsonVersion == true {
		fmt.Printf("{\"serverversion\":\"%s\",\"apiversion\":\"%s\"}", ServerVersion, APIVersion)
	}

	// Load config
	var Config Configuration
	Configfile, err := ioutil.ReadFile(*ConfigLocation)
	if err != nil {
		fmt.Println("Couldn't find or open config file. Create one with -createconfig")
		panic(err)
	}
	err = fromjson(string(Configfile), &Config)
	if err != nil {
		panic(err)
	}

	// Setup database connection
	DatabaseConnectionString := fmt.Sprintf("dbname=trackmon_server_production user=trackmon host=%s password=%s", Config.DatabaseAddress, Config.DatabasePassword)
	db, err := sql.Open("postgres", DatabaseConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare database statements and setup database
	DatabaseSetup(db)

	PrepGetUserQuery, err = db.Prepare(GetUserQuery)
	if err != nil {
		panic(err)
	}
	defer PrepGetUserQuery.Close()

	PrepDoesUserExistQuery, err = db.Prepare(DoesUserExistQuery)
	if err != nil {
		panic(err)
	}
	defer PrepDoesUserExistQuery.Close()

	PrepAddNewUser, err = db.Prepare(AddNewUser)
	if err != nil {
		panic(err)
	}
	defer PrepAddNewUser.Close()

	PrepDeleteExistingUser, err = db.Prepare(DeleteExistingUser)
	if err != nil {
		panic(err)
	}
	defer PrepDeleteExistingUser.Close()

	PrepDeleteAccountsFromExistingUser, err = db.Prepare(DeleteAccountsFromExistingUser)
	if err != nil {
		panic(err)
	}
	defer PrepDeleteAccountsFromExistingUser.Close()

	PrepDeleteHistoryFromExistingUser, err = db.Prepare(DeleteHistoryFromExistingUser)
	if err != nil {
		panic(err)
	}
	defer PrepDeleteHistoryFromExistingUser.Close()

	// Configure router and server
	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler) // Returnes 200 OK, can be used for health checks
	r.HandleFunc("/version", VersionHandler)

	r.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		AllAccountHandler(w, r, db)
	})

	r.HandleFunc("/account/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		HistoryHandler(w, r, db)
	})

	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		UserHandler(w, r, db)
	})

	srv := &http.Server{
		Handler: r,
		Addr:    Config.ListeningAddress,
	}

	// Check update
	if Config.AutoUpdateChecker != false {
		go checkupdate("https://api.github.com/repo/trackmon/trackmon-server/releases/latest", ServerVersion)
	}

	// Start the server
	log.Println("Initialization Complete")
	log.Fatal(srv.ListenAndServe())
}

func DatabaseSetup(db *sql.DB) {
	_, err := db.Exec(DatabaseSetupUsersTable)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(DatabaseSetupAccountsTable)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(DatabaseSetupHistoryTable)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"database/sql"
	"io/ioutil"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"log"
	"net/http"
)

var (
	ServerVersion string = "pre version"
	APIVersion    string = "pre api"
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

	// Configure router and server
	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler) // Returnes 200 OK, can be used for health checks
	r.HandleFunc("/version", VersionHandler)

	r.HandleFunc("/user/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		UserHandler(w, r, db)
	})
	r.HandleFunc("/user/{user_id}/{account}", func(w http.ResponseWriter, r *http.Request) {
		AccountHandler(w, r, db)
	})
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		NewUserHandler(w, r, db)
	})

	srv := &http.Server{
		Handler: r,
		Addr:    Config.ListeningAddress,
	}

	// Start the server
	log.Println("Initialization complete")
	log.Fatal(srv.ListenAndServe())
}

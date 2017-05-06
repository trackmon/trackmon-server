package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
	_ "github.com/lib/pq"
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
	db_connection_string := fmt.Sprintf("dbname=trackmon_server_production user=trackmon host=%s password=%s", Config.DatabaseAddress, Config.DatabasePassword)
	db, err := sql.Open("postgres", db_connection_string)
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
	r.HandleFunc("/user/{user_id}", UserHandler)
	r.HandleFunc("/user/{user_id}/{account}", AccountHandler)
	r.HandleFunc("/signup", NewUserHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    Config.ListeningAddress,
	}

	// Start the server
	log.Println("Initialization complete")
	srv.ListenAndServe()
}

/*
 ██████  ██████  ███    ██ ███████ ██  ██████
██      ██    ██ ████   ██ ██      ██ ██
██      ██    ██ ██ ██  ██ █████   ██ ██   ███
██      ██    ██ ██  ██ ██ ██      ██ ██    ██
 ██████  ██████  ██   ████ ██      ██  ██████
*/

type Configuration struct {
	ListeningAddress string
	DatabaseAddress  string
	DatabasePassword string
}

func CreateConfig() {
	var Config Configuration

	// Standard config
	Config.ListeningAddress = ":80"
	Config.DatabaseAddress = "localhost"
	Config.DatabasePassword = ""

	ByteJsonConfig, err := toprettyjson(Config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Writing configuration to ./trackmonserv.conf")
	err = ioutil.WriteFile("./trackmonserv.conf", ByteJsonConfig, 0644)
	if err != nil {
		panic(err)
	}
}

/*
██   ██  █████  ███    ██ ██████  ██      ███████ ██████
██   ██ ██   ██ ████   ██ ██   ██ ██      ██      ██   ██
███████ ███████ ██ ██  ██ ██   ██ ██      █████   ██████
██   ██ ██   ██ ██  ██ ██ ██   ██ ██      ██      ██   ██
██   ██ ██   ██ ██   ████ ██████  ███████ ███████ ██   ██
*/

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"serverversion\":\"%s\",\"apiversion\":\"%s\"}", ServerVersion, APIVersion)
	w.WriteHeader(http.StatusOK)
}

func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	hashedpw, err := HashPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(username, hashedpw)
	// TODO: Check if user exist in database
	// TODO: IF NOT Create new user and write to database
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
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

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	username, password, ok := r.BasicAuth()
	if ok != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} // All requests below here have given basic auth
	log.Printf("User %s with pw %s wants info about account %s from %s\n", username, password, string(variables["account"]), string(variables["user_id"]))
}

/*
 ██████ ██████  ██    ██ ██████  ████████
██      ██   ██  ██  ██  ██   ██    ██
██      ██████    ████   ██████     ██
██      ██   ██    ██    ██         ██
 ██████ ██   ██    ██    ██         ██
*/

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
     ██ ███████  ██████  ███    ██
     ██ ██      ██    ██ ████   ██
     ██ ███████ ██    ██ ██ ██  ██
██   ██      ██ ██    ██ ██  ██ ██
 █████  ███████  ██████  ██   ████
*/

func fromjson(src string, v interface{}) error {
	return json.Unmarshal([]byte(src), v)
}

func tojson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func toprettyjson(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}

/*
██      ██  ██████ ███████ ███    ██ ███████ ███████
██      ██ ██      ██      ████   ██ ██      ██
██      ██ ██      █████   ██ ██  ██ ███████ █████
██      ██ ██      ██      ██  ██ ██      ██ ██
███████ ██  ██████ ███████ ██   ████ ███████ ███████
*/

const (
	muxlicense string = `Copyright (c) 2012 Rodrigo Moraes. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

	 * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
	 * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
	 * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`
	trackmonlicense string = `BSD 2-Clause License

	Copyright (c) 2017, Paul Kramme
	All rights reserved.

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are met:

	* Redistributions of source code must retain the above copyright notice, this
	  list of conditions and the following disclaimer.

	* Redistributions in binary form must reproduce the above copyright notice,
	  this list of conditions and the following disclaimer in the documentation
	  and/or other materials provided with the distribution.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
)

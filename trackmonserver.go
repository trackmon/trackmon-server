package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ServerVersion string = "pre version"
	APIVersion string = "pre api"
)

type Configuration struct {
	ListeningAddress string
}

func CreateConfig() {
	var Config Configuration

	// Standard config
	Config.ListeningAddress = ":80"

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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"serverversion\":\"%s\",\"apiversion\":\"%s\"}", ServerVersion, APIVersion)
	w.WriteHeader(http.StatusOK)
}

func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// TODO: Write function which creates new users and stores them on the DB
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	// log.Printf("User %s wants info", string(variables["user_id"]))
	// TODO: Check if user exists, if and if password correct, give him info
}

func main() {
	fmt.Println("TRACKMON SERVER licensed under BSD 2-Clause")
	fmt.Println("Please report bugs to https://github.com/trackmon/trackmon-server")

	// Configure flags
	CreateConfigFlag := flag.Bool("createconfig", false, "Creates a standard configuration and exits")
	ConfigLocation := flag.String("config", "./trackmonserv.conf", "Location of config file. Standard is ./trackmonserv")
	ShowLicenses := flag.Bool("licenses", false, "Shows licenses and exits")
	flag.Parse()

	// Check flags
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

	// Configure router and server
	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler) // Returnes 200 OK, can be used for health checks
	r.HandleFunc("/version", VersionHandler)
	r.HandleFunc("/user", NewUserHandler)
	r.HandleFunc("/user/{user_id}", UserHandler)
	srv := &http.Server{
		Handler: r,
		Addr:    Config.ListeningAddress,
	}

	// Start the server
	log.Println("Initialization complete")
	srv.ListenAndServe()
}

func fromjson(src string, v interface{}) error {
	return json.Unmarshal([]byte(src), v)
}

func tojson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func toprettyjson(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}

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

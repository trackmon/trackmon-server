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

type Configuration struct {
	ListeningAddress string
}

func CreateConfig() {
	var Config Configuration
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

func main() {
	fmt.Println("trackmon server by Paul Kramme")
	fmt.Println("Please report bugs to https://github.com/trackmon/trackmon-server")

	// Configure flags
	CreateConfigFlag := flag.Bool("createconfig", false, "Creates a standard configuration and exits")
	ConfigLocation := flag.String("config", "./trackmonserv.conf", "Location of config file. Standard is ./trackmonserv")
	flag.Parse()

	// Check flags
	if *CreateConfigFlag == true {
		CreateConfig()
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
	r.HandleFunc("/", RootHandler) // Returnes 200 OK, can be used for health checks.

	srv := &http.Server{
		Handler: r,
		Addr:    Config.ListeningAddress,
	}

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
	return json.MarshalIndent(v, "", "    ")
}

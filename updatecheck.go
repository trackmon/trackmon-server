package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GithubReleasesAssetsApiResponse struct {
	Browser_download_url string
}

type GithubReleasesApiResponse struct {
	Message  string
	Tag_name string
	Html_url string
	Assets   []GithubReleasesAssetsApiResponse
}

func checkupdate(url string, version string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Autoupdate:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var ghapiresp GithubReleasesApiResponse
	err = fromjson(string(body), &ghapiresp)
	if ghapiresp.Message == "Not Found" {
		return
	}
	if ghapiresp.Tag_name != version {
		fmt.Printf("New version %s available at %s\n", ghapiresp.Tag_name, ghapiresp.Html_url)
	}
}

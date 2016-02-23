package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xtraclabs/signup/auth"
	"github.com/xtraclabs/signup/handlers"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func validateEnvironment() {
	if os.Getenv("ROLL_ENDPOINT") == "" {
		log.Println("ROLL_ENDPOINT environment variable not set.")
		os.Exit(1)
	}
}

func loginAsRollApp() {
	var rollClientId = os.Getenv("ROLL_CLIENTID")
	if rollClientId == "" {
		log.Println("ROLL_CLIENTID environment variable not set.")
		os.Exit(1)
	}

	var rollClientSecret = os.Getenv("ROLL_SECRET")
	if rollClientSecret == "" {
		log.Println("ROLL_SECRET environment variable not set.")
		os.Exit(1)
	}

	resp, err := http.PostForm(os.Getenv("ROLL_ENDPOINT")+"/oauth2/token",
		url.Values{"grant_type": {"password"},
			"client_id":     {rollClientId},
			"client_secret": {rollClientSecret},
			"username":      {"abc"},
			"password":      {"xxxxxxxx"}})

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

	var jsonResponse accessTokenResponse
	err = json.Unmarshal([]byte(body), &jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Use this token", jsonResponse.AccessToken)

	auth.AccessToken = jsonResponse.AccessToken

}

func main() {

	validateEnvironment()
	loginAsRollApp()

	var port = flag.Int("port", -1, "Port to listen on")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	log.Println("Listening on port ", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), handlers.Handler())
}

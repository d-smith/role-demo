package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xtraclabs/roll/roll"
	"html/template"
	"log"
	"net/http"
	"os"
)

var rollEndpoint = os.Getenv("ROLL_ENDPOINT")

var templates = template.Must(template.ParseFiles("../html/signup.html", "../html/thanks.html", "../html/createapp.html",
	"../html/appdetails.html", "../html/error.html"))

func handleSignupGet(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "signup.html", nil); err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
}

func handleSignupPost(w http.ResponseWriter, r *http.Request) {

	log.Println("extract signup form values")
	dev := roll.Developer{
		Email:     r.FormValue("email"),
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
	}

	log.Println("encode developer json")
	bodyReader := new(bytes.Buffer)
	enc := json.NewEncoder(bodyReader)
	err := enc.Encode(dev)

	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	log.Println("form put request")
	requestUrl := fmt.Sprintf("%s/v1/developers/%s", rollEndpoint, dev.Email)
	log.Println("put to", requestUrl)
	req, err := http.NewRequest("PUT", requestUrl, bodyReader)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	log.Println("submit put request")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		renderErrorPage(w, resp.StatusCode, responseAsString(resp))
		return
	}

	log.Println("store of developer info ok")
	ctx := devContext{Email: dev.Email}
	if err := templates.ExecuteTemplate(w, "thanks.html", ctx); err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
}

func handleSignup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSignupGet(w, r)
		case "POST":
			handleSignupPost(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

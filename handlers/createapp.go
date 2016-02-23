package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/signup/auth"
	"io/ioutil"
	"log"
	"net/http"
)

type devContext struct {
	Email string
}

type clientId struct {
	ClientId string `json:"client_id"`
}

func handleCreateGet(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	ctx := devContext{Email: email}
	if err := templates.ExecuteTemplate(w, "createapp.html", ctx); err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {

	//Define the app def using the form fields
	app := roll.Application{
		DeveloperEmail:  r.FormValue("email"),
		ApplicationName: r.FormValue("appname"),
		RedirectURI:     r.FormValue("redirecturi"),
		LoginProvider:   r.FormValue("loginprovider"),
	}

	//Create the json payload
	log.Println("encode application json")
	bodyReader := new(bytes.Buffer)
	enc := json.NewEncoder(bodyReader)
	err := enc.Encode(app)

	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Create the post request
	log.Println("form post request")
	requestUrl := fmt.Sprintf("%s/v1/applications", rollEndpoint)
	log.Println("post to", requestUrl)
	req, err := http.NewRequest("POST", requestUrl, bodyReader)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Add the auth header to the request
	req.Header.Add("Authorization", "Bearer "+auth.AccessToken)

	//Post the request
	log.Println("submit post request")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error on put", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//If there's an error load the error page
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		renderErrorPage(w, resp.StatusCode, responseAsString(resp))
		return
	}

	//Read the body and grab the returned client id
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		renderErrorPage(w, resp.StatusCode, err.Error())
		return
	}

	var jsonResponse clientId
	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		renderErrorPage(w, resp.StatusCode, err.Error())
		return
	}

	//Still here? Grab the app details
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/v1/applications/%s", rollEndpoint, jsonResponse.ClientId), nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+auth.AccessToken)

	log.Println("retrieve full definition to include client secret")
	retResp, err := client.Do(req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	defer retResp.Body.Close()

	if retResp.StatusCode != http.StatusOK {
		renderErrorPage(w, retResp.StatusCode, responseAsString(resp))
		return
	}

	log.Println("decode response into json")
	var decodedApp roll.Application
	defer retResp.Body.Close()

	dec := json.NewDecoder(retResp.Body)
	if err := dec.Decode(&decodedApp); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if err := templates.ExecuteTemplate(w, "appdetails.html", decodedApp); err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
}

func handleCreateApp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleCreateGet(w, r)
		case "POST":
			handleCreatePost(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

type devContext struct {
	Email string
}

func handleCreateGet(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	ctx := devContext{Email: email}
	if err := templates.ExecuteTemplate(w, "createapp.html", ctx); err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {

	log.Println("extract form values and form application detail")
	u, err := uuid.NewV4()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	app := roll.Application{
		DeveloperEmail:  r.FormValue("email"),
		ClientID:        u.String(),
		ApplicationName: r.FormValue("appname"),
		RedirectURI:     r.FormValue("redirecturi"),
		LoginProvider:   r.FormValue("loginprovider"),
	}

	log.Println("encode application json")
	bodyReader := new(bytes.Buffer)
	enc := json.NewEncoder(bodyReader)
	err = enc.Encode(app)

	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	log.Println("form put request")
	requestUrl := fmt.Sprintf("%s/v1/applications/%s", rollEndpoint, app.ClientID)
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
		log.Println("error on put", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		renderErrorPage(w, resp.StatusCode, responseAsString(resp))
		return
	}

	log.Println("retrieve full definition to include client secret")
	retResp, err := http.Get(fmt.Sprintf("%s/v1/applications/%s", rollEndpoint, app.ClientID))
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

package handlers

import (
	"errors"
	"net/http"
)

type appDetails struct {
	ClientID         string
	ApplicationName  string
	ClientSecret     string
	RedirectURI      string
	LoginProvider    string
	JWTFlowPublicKey string
}

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

	details := appDetails{
		ApplicationName: "foo app",
		ClientID:        "id",
		ClientSecret:    "don't tell",
		RedirectURI:     "redirect",
		LoginProvider:   "login provider",
	}

	if err := templates.ExecuteTemplate(w, "appdetails.html", details); err != nil {
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

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

func handleCreateGet(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "createapp.html", nil); err != nil {
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

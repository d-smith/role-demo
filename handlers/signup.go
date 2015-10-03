package handlers

import (
	"errors"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("../html/signup.html", "../html/thanks.html"))

func handleSignupGet(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "signup.html", nil); err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
}

func handleSignupPost(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "thanks.html", nil); err != nil {
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

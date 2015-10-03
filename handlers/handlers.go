package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func respondError(w http.ResponseWriter, status int, err error) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/signup", handleSignup())
	return mux
}

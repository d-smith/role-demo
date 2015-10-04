package handlers

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
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
	mux.Handle("/createapp", handleCreateApp())
	return mux
}


func responseAsString(r *http.Response) string {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "error reading body"
	}

	return string(body)
}
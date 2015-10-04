package handlers

import "net/http"

type ErrorCtx struct {
	StatusCode   int
	ErrorMessage string
}

func renderErrorPage(w http.ResponseWriter, statusCode int, errorMessage string) {
	errorCtx := ErrorCtx{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}

	err := templates.ExecuteTemplate(w, "error.html", errorCtx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
}

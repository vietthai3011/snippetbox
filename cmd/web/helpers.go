package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL
		trade  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "Method", method, "URL", url, "Trace", trade)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(
	w http.ResponseWriter, r *http.Request, status int, page string, data templateData,
) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("The template %s does exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", page)
	if err != nil {
		app.serverError(w, r, err)
	}
}

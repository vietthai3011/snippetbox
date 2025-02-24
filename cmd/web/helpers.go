package main

import (
	"bytes"
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
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(
	w http.ResponseWriter, r *http.Request, status int, page string, data templateData,
) {
	// page dược truyền vào để lấy *template.Template
	ts, ok := app.templateCache[page]

	// nếu page không có trong map thì trả về lỗi sai
	if !ok {
		err := fmt.Errorf("the template %s does exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

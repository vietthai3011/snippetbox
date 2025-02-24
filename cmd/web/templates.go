package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/vietthai3011/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var function = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// key(name) : value(page)
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(function).ParseFiles("./ui/html/base.html")

		if err != nil {
			return nil, err
		}

		// Tự động thêm tất cả template con trong partials/, không cần viết tay.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

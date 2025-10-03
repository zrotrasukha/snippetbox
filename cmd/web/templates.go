package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/zrotrasukha/snippetbox/internal/models"
)

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
	Form        any
	Flash       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// parsing base file
		ts, err := template.New(name).Funcs(functions).ParseFiles("ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// parsing partial file
		ts, err = ts.ParseGlob("ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}
	return cache, nil
}

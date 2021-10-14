package controller

import (
	"embed"
	"html/template"
	"net/http"
)

// NewController ...
func NewController(homeHTML embed.FS) Controller {
	return Controller{
		homeHTML: homeHTML,
	}
}

// Controller ...
type Controller struct {
	homeHTML embed.FS
}

// HandleHome ...
func (c Controller) HandleHome(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(c.homeHTML, "index.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

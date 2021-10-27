package controller

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	//"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting/impl"
)

// NewController ...
func NewController(homeHTML embed.FS, homepage embed.FS, vs impl.VotingSystem) Controller {
	return Controller{
		homeHTML: homeHTML,
		homepage: homepage,
		vs:       vs,
	}
}

// Controller ...
type Controller struct {
	homeHTML embed.FS
	homepage embed.FS
	vs       impl.VotingSystem
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

	//req.Form.Get(election//)

	name := req.FormValue("username")
	description := req.FormValue("description")
	roomID := req.FormValue("roomID")

	//Only print. Have to be stored on database
	fmt.Fprintln(w, "Username = \n", name)
	fmt.Fprintln(w, "Description = \n", description)
	fmt.Fprintln(w, "RoomID = \n", roomID)
}

func (c Controller) HandleHomePage(w http.ResponseWriter, req *http.Request) {
	t2, err := template.ParseFS(c.homepage, "homepage.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = t2.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

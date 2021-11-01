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
func NewController(homeHTML embed.FS, homepage embed.FS, views embed.FS, vs impl.VotingSystem) Controller {
	return Controller{
		homeHTML: homeHTML,
		homepage: homepage,
		views:    views,
		vs:       vs,
	}
}

// Controller ...
type Controller struct {
	homeHTML embed.FS
	homepage embed.FS
	views    embed.FS
	vs       impl.VotingSystem
}

// HandleHome ...
func (c Controller) HandleHome(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		http.Error(w, "Not found.", http.StatusNotFound)
		return
	}

	t, err := template.ParseFS(c.homeHTML, "web/index.html")
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

	t2, err := template.ParseFS(c.homepage, "web/homepage.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title             string
		VotingInstanceTab map[string]impl.VotingInstance
	}{Title: "TestTitle", VotingInstanceTab: c.vs.VotingInstancesList}

	err = t2.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// id := req.URL.Query().Get("id")
	// fmt.Fprintln(w, "URL :::::::: ", req.URL)
	// if id == "" {
	// 	http.Error(w, "The id query parameter is missing", http.StatusBadRequest)
	// 	return
	// }

	//creating a button for all the differents voting instances created

	for _, v := range c.vs.VotingInstancesList {
		var s string = "<input type=\"button\" name=\"RoomID\" value=" + "\"" + v.Id + "\"" + " onclick=\"self.location.href='/homepage/" + v.Id + "'\" >"
		w.Write([]byte(s))
	}
	//if err != nil {
	// Handle error here via logging and then return
	//} else if req.Method != "POST" {
	//	http.Redirect(w, req, "/homepage/"+ids, 301)
	//}

}

func (c Controller) HandleShowElection(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFS(c.views, "web/views/election.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = req.ParseForm()
	if err != nil {
		//error
		return
	}

	id := req.Form.Get("id")
	if id == "" {
		http.Error(w, "failed to get id: ", http.StatusInternalServerError)
		return
	}

	election, found := c.vs.VotingInstancesList[id]
	if !found {
		http.Error(w, "Election not found: "+id, http.StatusInternalServerError)
		return
	}

	data := struct {
		Election impl.VotingInstance
	}{
		Election: election,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

	status := c.vs.VotingInstancesList["001"].Status
	title := c.vs.VotingInstancesList["001"].Config.Title
	description := c.vs.VotingInstancesList["001"].Config.Description
	voters := c.vs.VotingInstancesList["001"].Config.Voters
	w.Write([]byte("Current status : " + status))
	w.Write([]byte("<br>Title : " + title))
	w.Write([]byte("<br>Description : " + description))
	w.Write([]byte("<br>List of voters : "))
	for _, v := range voters {
		w.Write([]byte(v))
	}
}

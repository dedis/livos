package controller

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	//"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting/impl"
	//"honnef.co/go/js/dom"
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
		http.Error(w, "Not found the path.", http.StatusNotFound)
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
	}{Title: "HomePage", VotingInstanceTab: c.vs.VotingInstancesList}

	err = t2.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// d := dom.GetWindow().Document()

	// button := d.GetElementByID("create")

	// button.AddEventListener("click", false, func(event dom.Event) {
	// 	fmt.Println("CLICKCLICK")
	// })

	if req.Method == "POST" {
		title := req.PostFormValue("title")
		if title == "" {
			http.Error(w, "failed to get title: ", http.StatusInternalServerError)
			return
		}
		id := req.FormValue("id")
		if id == "" {
			http.Error(w, "failed to get id: ", http.StatusInternalServerError)
			return
		}
		status := req.FormValue("status")
		if status == "" {
			http.Error(w, "failed to get status: ", http.StatusInternalServerError)
			return
		}
		description := req.FormValue("desc")
		if description == "" {
			http.Error(w, "failed to get description: ", http.StatusInternalServerError)
			return
		}
		voterList := req.FormValue("votersList")
		if voterList == "" {
			http.Error(w, "failed to get list of voters: ", http.StatusInternalServerError)
			return
		}
		voterListParsed := strings.Split(voterList, ",")
		candidats := req.FormValue("candidates")
		if candidats == "" {
			http.Error(w, "failed to get list of candidates: ", http.StatusInternalServerError)
			return
		}
		candidatesParsed := strings.Split(candidats, ",")

		fmt.Fprintln(w, "TEST DE PRINT POUR VOIR SI RECUP VALUE FONCTIONNE")
		fmt.Fprintln(w, "Title = \n", title)
		fmt.Fprintln(w, "Description = \n", description)
		fmt.Fprintln(w, "Status = \n", status)
		fmt.Fprintln(w, "id = \n", id)
		fmt.Fprintln(w, "List of voters = \n", voterListParsed[0], voterListParsed[1], voterListParsed[2])
		fmt.Fprintln(w, "List of candidates = \n", candidatesParsed[0])
	}

	//recup les donnes form.get
	//listener .. Create() creer la votinginstance
}

func (c Controller) HandleShowElection(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFS(c.views, "web/views/election.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
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
		id       string
	}{
		Election: election,
		id:       id,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

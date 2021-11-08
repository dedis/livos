package controller

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	//"strings"

	//"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting"

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
	// name := req.FormValue("username")
	// description := req.FormValue("description")
	// roomID := req.FormValue("roomID")
}

func (c Controller) HandleHomePage(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t2, err := template.ParseFS(c.homepage, "web/homepage.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	VotingInstanceTabOpen := make(map[string]*impl.VotingInstance)
	VotingInstanceTabClose := make(map[string]*impl.VotingInstance)

	for key, value := range c.vs.VotingInstancesList {
		if value.Status == "open" {
			VotingInstanceTabOpen[key] = value
		} else {
			VotingInstanceTabClose[key] = value
		}

	}

	data := struct {
		Title                  string
		VotingInstanceTab      map[string]*impl.VotingInstance
		VotingInstanceTabOpen  map[string]*impl.VotingInstance
		VotingInstanceTabClose map[string]*impl.VotingInstance
	}{Title: "HomePage",
		VotingInstanceTab:      c.vs.VotingInstancesList,
		VotingInstanceTabOpen:  VotingInstanceTabOpen,
		VotingInstanceTabClose: VotingInstanceTabClose}

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

		// fmt.Fprintln(w, "TEST DE PRINT POUR VOIR SI RECUP VALUE FONCTIONNE")
		fmt.Println("Title = \n", title)
		// fmt.Fprintln(w, "Description = \n", description)
		// fmt.Fprintln(w, "Status = \n", status)
		// fmt.Fprintln(w, "id = \n", id)
		//fmt.Println("List of voters = \n", voterListParsed[0], voterListParsed[1], voterListParsed[2])
		//fmt.Println("List of candidates = \n", candidatesParsed[0])

		votingConfig, err := impl.NewVotingConfig(voterListParsed, title, description, candidatesParsed)
		if err != nil {
			http.Error(w, "NewVotingConfig is incorrect", http.StatusInternalServerError)
		}

		votes := make(map[string]*voting.Choice)
		c.vs.CreateAndAdd(id, votingConfig, status, votes)

		http.Redirect(w, req, "/homepage", http.StatusSeeOther)
	}

	err = t2.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
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

	electionAdd, found := c.vs.VotingInstancesList[id]
	if !found {
		http.Error(w, "Election not found: "+id, http.StatusInternalServerError)
		return
	}

	deleg := make(map[string]voting.Liquid)
	yesChoice := make(map[string]voting.Liquid)
	liq100, err := impl.NewLiquid(100)
	if err != nil {
		http.Error(w, "Liquid creation incorrect", http.StatusInternalServerError)
	}

	yesChoice["yes"] = liq100
	choiceGuillaume, err := impl.NewChoice(deleg, yesChoice, 0, 100)
	if err != nil {
		http.Error(w, "Choice creation incorrect", http.StatusInternalServerError)
	}

	data := struct {
		Election impl.VotingInstance
		id       string
		Choice   voting.Choice
	}{
		Election: *electionAdd,
		id:       id,
		Choice:   choiceGuillaume,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c Controller) HandleShowResults(w http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFS(c.views, "web/views/results.html")
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

	electionAdd, found := c.vs.VotingInstancesList[id]
	if !found {
		http.Error(w, "Election not found: "+id, http.StatusInternalServerError)
		return
	}

	fmt.Println("ELECTION ADD object: ", *electionAdd)

	data := struct {
		Election impl.VotingInstance
		id       string
		Results  map[string]float64
	}{
		Election: *electionAdd,
		id:       id,
		Results:  electionAdd.GetResults(),
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c Controller) HandleManageVoting(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFS(c.views, "web/views/manage.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id := req.Form.Get("id")
	if id == "" {
		http.Error(w, "failed to get id: ", http.StatusInternalServerError)
		return
	}

	electionAdd, found := c.vs.VotingInstancesList[id]
	if !found {
		http.Error(w, "Election not found: "+id, http.StatusInternalServerError)
		return
	}

	data := struct {
		Election impl.VotingInstance
		id       string
	}{
		Election: *electionAdd,
		id:       id,
	}

	if req.Method == "POST" {

		title := req.PostFormValue("title")
		if title != "" {
			c.vs.VotingInstancesList[id].Config.Title = title
		}
		status := req.FormValue("status")
		c.vs.VotingInstancesList[id].SetStatus(status)

		description := req.FormValue("desc")
		if description != "" {
			c.vs.VotingInstancesList[id].Config.Description = description
		}

		voterList := req.FormValue("votersList")
		voterListParsed := strings.Split(voterList, ",")
		if voterList != "" {
			c.vs.VotingInstancesList[id].Config.Voters = voterListParsed
		}

		candidats := req.FormValue("candidates")
		candidatesParsed := strings.Split(candidats, ",")
		if candidats != "" {
			c.vs.VotingInstancesList[id].Config.Candidates = candidatesParsed

		}

		http.Redirect(w, req, "/homepage", http.StatusSeeOther)

	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

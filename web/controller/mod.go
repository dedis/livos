package controller

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
)

// NewController ...
func NewController(homeHTML embed.FS, homepage embed.FS, views embed.FS, vs voting.VotingSystem) Controller {
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
	vs       voting.VotingSystem
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

	VotingInstanceTabOpen := make(map[string]voting.VotingInstance)
	VotingInstanceTabClose := make(map[string]voting.VotingInstance)

	for key, value := range c.vs.GetVotingInstanceList() {
		if value.GetStatus() == "open" {
			VotingInstanceTabOpen[key] = value
		} else {
			VotingInstanceTabClose[key] = value
		}
	}

	data := struct {
		Title                  string
		VotingInstanceTab      map[string]voting.VotingInstance
		VotingInstanceTabOpen  map[string]voting.VotingInstance
		VotingInstanceTabClose map[string]voting.VotingInstance
	}{Title: "HomePage",
		VotingInstanceTab:      c.vs.GetVotingInstanceList(),
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
		//removing all whitespace
		voterList = strings.ReplaceAll(voterList, " ", "")

		//parsing the list to get usernames
		voterListParsed := strings.Split(voterList, ",")

		fmt.Println("List of voters : ", voterListParsed)
		candidats := req.FormValue("candidates")
		// if candidats == "" {
		// 	http.Error(w, "failed to get list of candidates: ", http.StatusInternalServerError)
		// 	return
		// }
		candidats = strings.ReplaceAll(candidats, " ", "")
		candidatesParsed := strings.Split(candidats, ",")

		delegTo := make(map[string]voting.Liquid)
		delegFrom := make(map[string]voting.Liquid)
		//userListParsed := make([]voting.User, 0)
		voterListParsedintoUser := make([]*voting.User, len(voterListParsed))
		histoChoice := make([]voting.Choice, 0)
		for idx, name := range voterListParsed {
			u, err := c.vs.NewUser(name, delegTo, delegFrom, histoChoice)
			//userListParsed = append(userListParsed, u)
			if err != nil {
				http.Error(w, "User creation is incorrect"+err.Error(), http.StatusInternalServerError)
			}
			fmt.Println("The user created is : ", u)
			voterListParsedintoUser[idx] = &u //&userListParsed[idx]
		}

		votingConfig, err := impl.NewVotingConfig(voterListParsedintoUser, title, description, candidatesParsed)
		if err != nil {
			http.Error(w, "NewVotingConfig is incorrect"+err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("The voting config is : ", votingConfig)

		//votes := make(map[string]voting.Choice)
		_, err = c.vs.CreateAndAdd(id, votingConfig, status)
		if err != nil {
			http.Error(w, "CreateAndAdd is incorrect"+err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, req, "/homepage", http.StatusSeeOther)
	}

	err = t2.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c Controller) HandleShowElection(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFS(c.views, "web/views/election.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id := req.Form.Get("id")
	if id == "" {
		http.Error(w, "failed to get id (id is null) ", http.StatusInternalServerError)
		return
	}

	electionAdd, found := c.vs.GetVotingInstanceList()[id]
	if !found {
		http.Error(w, "Election not found: "+id, http.StatusInternalServerError)
		return
	}

	if req.Method == "POST" {

		// VOTING CHOICE : VOTER, YESCHOICE, NOCHOICE => CASTVOTE -----------------

		//get the voter name for the vote to be cast
		voter := req.PostFormValue("voter")
		if voter == "" {
			http.Error(w, "failed to get voter: ", http.StatusInternalServerError)
			return
		}

		//get the YES value of the the vote to be cast
		liquidYes := voting.Liquid{}
		YesChoice := req.PostFormValue("yesPercent")
		if YesChoice == "" {
			liquidYes, err = impl.NewLiquid(0)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
		} else {
			temp, err := strconv.ParseFloat(YesChoice, 64)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
			liquidYes, err = impl.NewLiquid(temp)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
		}

		//get the NO value of the the vote to be cast
		liquidNo := voting.Liquid{}
		NoChoice := req.PostFormValue("noPercent")
		if NoChoice == "" {
			liquidNo, err = impl.NewLiquid(0)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
		} else {
			temp, err := strconv.ParseFloat(NoChoice, 64)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
			liquidNo, err = impl.NewLiquid(temp)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
		}

		//get the user (object) from the retrieved name
		//fmt.Println("VOTER IS ::: ", voter)
		//fmt.Println("VOTER LIST is ::: ", electionAdd.Config.Voters)
		userVoter, err := electionAdd.GetUser(voter)
		if err != nil {
			http.Error(w, "User cannot be found.", http.StatusInternalServerError)
		}
		//fmt.Println("error is :::", err.Error())
		//fmt.Println("User is : ", userVoter)

		//construct the choice from the YES and NO values above
		choice := make(map[string]voting.Liquid)
		choice["yes"] = liquidYes
		choice["no"] = liquidNo
		choiceUser, err := impl.NewChoice(choice)
		if err != nil {
			http.Error(w, "Choice creation incorrect", http.StatusInternalServerError)
		}

		//set the choice to the user
		//fmt.Println("::::::00 Result of the setchoice of guillaume", userVoter.MyChoice)
		//fmt.Println(":::::: CHOICE USER choice of guillaume", choiceUser)

		//fmt.Println("address de uservoter", &userVoter)

		//fmt.Println("AVANT LE SET CHOICE : user = ", userVoter, "  choice = ", choiceUser)
		if liquidNo.Percentage != 0. || liquidYes.Percentage != 0. {
			err = electionAdd.SetVote(userVoter, choiceUser)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Println("ET APRES l'erreur : ")
			//fmt.Println("::::::11 Result of the setchoice of guillaume", userVoter.MyChoice)

			//cast the vote of the user
			// err = electionAdd.CastVote(userVoter)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// }
		}

		// DELEGATION : VOTER1, VOTER2, QUANTITY => DELEG_TO -----------------

		//get the voterSender name for the delegation
		voterSender := req.PostFormValue("voterSender")
		if voterSender == "" {
			http.Error(w, "failed to get voter: ", http.StatusInternalServerError)
			return
		}

		//get the voterReceiver name for the delegation
		voterReceiver := req.PostFormValue("voterReceiver")
		if voterReceiver == "" {
			http.Error(w, "failed to get voter: ", http.StatusInternalServerError)
			return
		}

		//get the QUANTITY value of the the vote to be cast
		liquidQuantity := voting.Liquid{}
		Quantity := req.PostFormValue("quantity")
		if Quantity == "" {
			liquidNo, err = impl.NewLiquid(0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			temp, err := strconv.ParseFloat(Quantity, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			liquidQuantity, err = impl.NewLiquid(temp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		//get the userSender (object) from the retrieved name
		userSender, err := electionAdd.GetUser(voterSender)
		if err != nil {
			http.Error(w, "User cannot be found."+err.Error(), http.StatusInternalServerError)
		}

		//get the userReceiver (object) from the retrieved name
		userReceiver, err := electionAdd.GetUser(voterReceiver)
		if err != nil {
			http.Error(w, "User cannot be found."+err.Error(), http.StatusInternalServerError)
		}

		if liquidQuantity.Percentage != 0. {
			fmt.Println("JUST BEFORE THE DELEG_TO")
			err = electionAdd.DelegTo(userSender, userReceiver, liquidQuantity)
			if err != nil {
				http.Error(w, "DelegTo incorrect"+err.Error(), http.StatusInternalServerError)
			}
			fmt.Println("JUST AFTER THE DELEG_TO")
		}

		http.Redirect(w, req, "/election?id="+id, http.StatusSeeOther)
	}

	data := struct {
		Election voting.VotingInstance
		id       string
	}{
		Election: electionAdd,
		id:       id,
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

	electionAdd, found := c.vs.GetVotingInstanceList()[id]
	if !found {
		http.Error(w, "Election not found: "+id, http.StatusInternalServerError)
		return
	}

	//fmt.Println("ELECTION ADD object: ", *electionAdd)
	results := electionAdd.GetResults()

	blanks := 100. - results["yes"] - results["no"]

	data := struct {
		Election voting.VotingInstance
		id       string
		Results  map[string]float64
		Blanks   float64
	}{
		Election: electionAdd,
		id:       id,
		Results:  results,
		Blanks:   blanks,
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

	electionAdd, found := c.vs.GetVotingInstanceList()[id]
	if !found {
		http.Error(w, "Election not found: "+id, http.StatusInternalServerError)
		return
	}

	data := struct {
		Election voting.VotingInstance
		id       string
	}{
		Election: electionAdd,
		id:       id,
	}

	if req.Method == "POST" {

		title := req.PostFormValue("title")
		if title != "" {
			c.vs.GetVotingInstanceList()[id].SetTitle(title)
		}
		status := req.FormValue("status")
		c.vs.GetVotingInstanceList()[id].SetStatus(status)

		description := req.FormValue("desc")
		if description != "" {
			c.vs.GetVotingInstanceList()[id].SetDescription(description)
		}

		voterList := req.FormValue("votersList")
		voterListParsed := strings.Split(voterList, ",")
		if voterList != "" {
			delegTo := make(map[string]voting.Liquid)
			delegFrom := make(map[string]voting.Liquid)
			voterListParsedintoUser := make([]*voting.User, len(voterListParsed))
			histoChoice := make([]voting.Choice, 0)
			for idx, name := range voterListParsed {
				u, err := c.vs.NewUser(name, delegTo, delegFrom, histoChoice)
				if err != nil {
					http.Error(w, "User creation is incorrect", http.StatusInternalServerError)
				}
				voterListParsedintoUser[idx] = &u
			}
			c.vs.GetVotingInstanceList()[id].SetVoters(voterListParsedintoUser)
		}

		candidats := req.FormValue("candidates")
		candidatesParsed := strings.Split(candidats, ",")
		if candidats != "" {
			c.vs.GetVotingInstanceList()[id].SetCandidates(candidatesParsed)

		}

		http.Redirect(w, req, "/homepage", http.StatusSeeOther)

	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

package controller

import (
	"embed"
	"fmt"
	"html/template"
	"math"
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

	length_open := len(VotingInstanceTabOpen)
	length_close := len(VotingInstanceTabClose)

	data := struct {
		Title                  string
		VotingInstanceTab      map[string]voting.VotingInstance
		VotingInstanceTabOpen  map[string]voting.VotingInstance
		VotingInstanceTabClose map[string]voting.VotingInstance
		Length_open            int
		Length_close           int
	}{Title: "HomePage",
		VotingInstanceTab:      c.vs.GetVotingInstanceList(),
		VotingInstanceTabOpen:  VotingInstanceTabOpen,
		VotingInstanceTabClose: VotingInstanceTabClose,
		Length_open:            length_open,
		Length_close:           length_close,
	}

	if req.Method == "POST" {
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
		if status == "open" {
			c.vs.GetVotingInstanceList()[id].SetStatus("close")
		} else {
			c.vs.GetVotingInstanceList()[id].SetStatus("open")
		}

		http.Redirect(w, req, "/homepage", http.StatusSeeOther)
	}
	err = t2.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c Controller) HandleShowElectionYesNo(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFS(c.views, "web/views/electionYesOrNoQuestion.html")
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
			liquidYes, err = impl.NewLiquid(int(temp))
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
			liquidNo, err = impl.NewLiquid(int(temp))
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
			//fmt.Println("ET APRES l'erreur : ")
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
			liquidQuantity, err = impl.NewLiquid(0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			temp, err := strconv.ParseFloat(Quantity, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			liquidQuantity, err = impl.NewLiquid(int(temp))
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
			//fmt.Println("JUST BEFORE THE DELEG_TO")
			err = electionAdd.DelegTo(userSender, userReceiver, liquidQuantity)
			if err != nil {
				http.Error(w, "DelegTo incorrect"+err.Error(), http.StatusInternalServerError)
			}
			//fmt.Println("JUST AFTER THE DELEG_TO")
		}

		http.Redirect(w, req, "/electionYesOrNoQuestion?id="+id, http.StatusSeeOther)
	}

	data := struct {
		Election voting.VotingInstance
		Id       string
	}{
		Election: electionAdd,
		Id:       id,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
func (c Controller) HandleGraphYesNo(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFS(c.views, "web/views/graphYesNo.html")
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

	stringConstruction := ""
	liquid_0, _ := impl.NewLiquid(0.)
	stringConstruction += "node [style=filled];yes[color=green];no[color=red];"

	for _, user := range electionAdd.GetConfig().Voters {
		cumulativeHistoryOfChoice := make([]voting.Choice, 0)
		new_vote_value := make(map[string]voting.Liquid)
		for _, choice := range user.HistoryOfChoice {
			for name, value := range choice.VoteValue {
				new_vote_value[name], err = impl.AddLiquid(new_vote_value[name], value)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		new_choice, err := impl.NewChoice(new_vote_value)
		if err != nil {
			fmt.Println(err.Error())
		}
		cumulativeHistoryOfChoice = append(cumulativeHistoryOfChoice, new_choice)

		for _, choice2 := range cumulativeHistoryOfChoice {
			for name2, valueToVote := range choice2.VoteValue {
				if valueToVote.Percentage > liquid_0.Percentage {
					if name2 == "yes" {
						stringConstruction += user.UserID + " ->" + name2 + "[color=green, label=" + strconv.FormatInt(int64(valueToVote.Percentage), 10) + "]" + ";"
					} else {
						stringConstruction += user.UserID + " ->" + name2 + "[color=red, label=" + strconv.FormatInt(int64(valueToVote.Percentage), 10) + "]" + ";"
					}

				}
			}
		}

		for nametodeleg, valueToDeleg := range user.DelegatedTo {
			temp := false
			for nameFromDeleg := range user.DelegatedFrom {
				if user.UserID == nameFromDeleg {
					temp = true
				}
			}
			if temp {
				stringConstruction += user.UserID + " ->" + nametodeleg + "[color=purple, label=" + strconv.FormatInt(int64(valueToDeleg.Percentage), 10) + "]" + ";"
			}
		}
	}

	data := struct {
		Election  voting.VotingInstance
		Id        string
		Infograph string
	}{
		Election:  electionAdd,
		Id:        id,
		Infograph: stringConstruction,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("this is the construction :", stringConstruction)

}

func (c Controller) HandleGraphCandidates(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFS(c.views, "web/views/graphCandidates.html")
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

	stringConstruction := ""
	liquid_0, _ := impl.NewLiquid(0.)
	stringConstruction += "node [style=filled];"
	listOfCOlorCand := []string{"", "firebrick", "magenta", "darkseagreen", "darkolivegreen"}
	colorOfCand := ""

	for l, cand := range electionAdd.GetConfig().Candidates {
		colorOfCand = listOfCOlorCand[int(math.Ceil(float64(l+1)/2))] + strconv.FormatInt(int64(2*l%4)+1, 10)
		stringConstruction += cand.CandidateID + "[color=" + colorOfCand + "];"
	}

	for _, user := range electionAdd.GetConfig().Voters {
		cumulativeHistoryOfChoice := make([]voting.Choice, 0)
		new_vote_value := make(map[string]voting.Liquid)
		for _, choice := range user.HistoryOfChoice {
			for name, value := range choice.VoteValue {
				new_vote_value[name], err = impl.AddLiquid(new_vote_value[name], value)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		new_choice, err := impl.NewChoice(new_vote_value)
		if err != nil {
			fmt.Println(err.Error())
		}
		cumulativeHistoryOfChoice = append(cumulativeHistoryOfChoice, new_choice)

		for _, choice2 := range cumulativeHistoryOfChoice {
			for name2, valueToVote := range choice2.VoteValue {
				if valueToVote.Percentage > liquid_0.Percentage {
					stringConstruction += user.UserID + " ->" + name2 + "[color=" + colorOfCand + ", label=" + strconv.FormatInt(int64(valueToVote.Percentage), 10) + "]" + ";"
				}
			}
		}

		for nametodeleg, valueToDeleg := range user.DelegatedTo {
			temp := false
			for nameFromDeleg := range user.DelegatedFrom {
				if user.UserID == nameFromDeleg {
					temp = true
				}
			}
			if temp {
				stringConstruction += user.UserID + " ->" + nametodeleg + "[color=purple, label=" + strconv.FormatInt(int64(valueToDeleg.Percentage), 10) + "]" + ";"
			}
		}
	}

	data := struct {
		Election  voting.VotingInstance
		Id        string
		Infograph string
	}{
		Election:  electionAdd,
		Id:        id,
		Infograph: stringConstruction,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("this is the construction :", stringConstruction)

}

func (c Controller) HandleShowElectionCandidate(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFS(c.views, "web/views/electionCandidateQuestion.html")
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

		candidate := req.PostFormValue("candidate")
		if candidate == "" {
			http.Error(w, "failed to get candidate: ", http.StatusInternalServerError)
			return
		}

		//get the value of the the vote to be cast for candidate
		liquidVote := voting.Liquid{}
		VoteChoice := req.PostFormValue("quantityPercent")
		if VoteChoice == "" {
			liquidVote, err = impl.NewLiquid(0)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
		} else {
			temp, err := strconv.ParseFloat(VoteChoice, 64)
			if err != nil {
				http.Error(w, "Creation of liquid is incorrect.", http.StatusInternalServerError)
			}
			liquidVote, err = impl.NewLiquid(int(temp))
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

		//construct the choice from the Vote value above for candidate
		choice := make(map[string]voting.Liquid)
		choice[candidate] = liquidVote
		choiceUser, err := impl.NewChoice(choice)
		if err != nil {
			http.Error(w, "Choice creation incorrect", http.StatusInternalServerError)
		}

		//set the choice to the user
		//fmt.Println("::::::00 Result of the setchoice of guillaume", userVoter.MyChoice)
		//fmt.Println(":::::: CHOICE USER choice of guillaume", choiceUser)

		//fmt.Println("address de uservoter", &userVoter)

		//fmt.Println("AVANT LE SET CHOICE : user = ", userVoter, "  choice = ", choiceUser)
		if liquidVote.Percentage != 0. {
			err = electionAdd.SetVote(userVoter, choiceUser)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			//fmt.Println("ET APRES l'erreur : ")
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
			liquidQuantity, err = impl.NewLiquid(0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			temp, err := strconv.ParseFloat(Quantity, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			liquidQuantity, err = impl.NewLiquid(int(temp))
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
			//fmt.Println("JUST BEFORE THE DELEG_TO")
			err = electionAdd.DelegTo(userSender, userReceiver, liquidQuantity)
			if err != nil {
				http.Error(w, "DelegTo incorrect"+err.Error(), http.StatusInternalServerError)
			}
			//fmt.Println("JUST AFTER THE DELEG_TO")
		}

		http.Redirect(w, req, "/electionCandidateQuestion?id="+id, http.StatusSeeOther)
	}

	data := struct {
		Election voting.VotingInstance
		Id       string
	}{
		Election: electionAdd,
		Id:       id,
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

	round := func(num float64) int {
		return int(num + math.Copysign(0.5, num))
	}

	toFixed := func(num float64, precision int) float64 {
		output := math.Pow(10, float64(precision))
		return float64(round(num*output)) / output
	}

	for s, res := range results {
		results[s] = toFixed(res, 2)
	}

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

	var listOfVoters = ""
	for _, user := range electionAdd.GetConfig().Voters {
		if listOfVoters == "" {
			listOfVoters = listOfVoters + user.UserID
		} else {
			listOfVoters = listOfVoters + "," + user.UserID
		}
	}

	var listOfCandidats = ""

	if len(electionAdd.GetConfig().Candidates) > 1 {
		for _, cand := range electionAdd.GetConfig().Candidates {
			if listOfCandidats == "" {
				listOfCandidats = listOfCandidats + cand.CandidateID
			} else {
				listOfCandidats = listOfCandidats + "," + cand.CandidateID
			}
		}
	}
	fmt.Println("CAndidattttttttttttts:", listOfCandidats)

	var description = electionAdd.GetConfig().Description

	data := struct {
		Election        voting.VotingInstance
		id              string
		ListOfVoters    string
		ListOfCandidats string
		Description     string
	}{
		Election:        electionAdd,
		id:              id,
		ListOfVoters:    listOfVoters,
		ListOfCandidats: listOfCandidats,
		Description:     description,
	}

	if req.Method == "POST" {

		title := req.PostFormValue("title")
		if title != "" {
			c.vs.GetVotingInstanceList()[id].SetTitle(title)
		}

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
			preferenceDelegationList := make([]*voting.User, 0)
			for idx, name := range voterListParsed {
				u, err := c.vs.NewUser(name, delegTo, delegFrom, histoChoice, voting.None, preferenceDelegationList)
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
			voterListParsedintoCandidate := make([]*voting.Candidate, len(candidatesParsed))
			fmt.Println("List of candidats", voterListParsedintoCandidate)
			for idx, name := range candidatesParsed {
				u, err := c.vs.NewCandidate(name)
				if err != nil {
					http.Error(w, "Candidate creation is incorrect", http.StatusInternalServerError)
				}
				voterListParsedintoCandidate[idx] = &u
			}
			fmt.Println("List of candidats", voterListParsedintoCandidate)
			c.vs.GetVotingInstanceList()[id].SetCandidates(voterListParsedintoCandidate)

		}

		http.Redirect(w, req, "/homepage", http.StatusSeeOther)

	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "failed to execute: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c Controller) HandleCreateVotingRoom(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse the form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t2, err := template.ParseFS(c.views, "web/views/createVotes.html")
	if err != nil {
		http.Error(w, "failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title             string
		VotingInstanceTab map[string]voting.VotingInstance
	}{Title: "Creation of Voting Room",
		VotingInstanceTab: c.vs.GetVotingInstanceList(),
	}

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
		typeofConfig := req.FormValue("typeOfConfig")
		if typeofConfig == "" {
			http.Error(w, "failed to get typeOfConfig: ", http.StatusInternalServerError)
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
		candidats = strings.ReplaceAll(candidats, " ", "")
		candidatesParsed := strings.Split(candidats, ",")
		voterListParsedintoCandidate := make([]*voting.Candidate, len(candidatesParsed))
		if candidats != "" {
			for idx, name := range candidatesParsed {
				u, err := c.vs.NewCandidate(name)
				if err != nil {
					http.Error(w, "Candidate creation is incorrect", http.StatusInternalServerError)
				}
				voterListParsedintoCandidate[idx] = &u
				fmt.Println("List of candidates : ", &voterListParsedintoCandidate)
			}
		}
		fmt.Println("CandidatParsed passed")

		delegTo := make(map[string]voting.Liquid)
		delegFrom := make(map[string]voting.Liquid)
		//userListParsed := make([]voting.User, 0)
		voterListParsedintoUser := make([]*voting.User, len(voterListParsed))
		histoChoice := make([]voting.Choice, 0)
		preferenceDelegationList := make([]*voting.User, 0)
		for idx, name := range voterListParsed {
			u, err := c.vs.NewUser(name, delegTo, delegFrom, histoChoice, voting.None, preferenceDelegationList)
			//userListParsed = append(userListParsed, u)
			if err != nil {
				http.Error(w, "User creation is incorrect"+err.Error(), http.StatusInternalServerError)
			}
			fmt.Println("The user created is : ", u)
			voterListParsedintoUser[idx] = &u //&userListParsed[idx]
		}
		fmt.Println("Creation of votingconfig is comming")

		votingConfig, err := impl.NewVotingConfig(voterListParsedintoUser, title, description, voterListParsedintoCandidate, voting.TypeOfVotingConfig(typeofConfig))
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

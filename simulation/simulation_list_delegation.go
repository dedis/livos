package simulation

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"github.com/mazen160/go-random"
	"github.com/yourbasic/graph"
	"golang.org/x/xerrors"
)

// GenerateItemsGraphviz creates a graphviz representation of the items. One can
// generate a graphical representation with `dot -Tpdf graph.dot -o graph.pdf`
func Simulation_list_delegation(out io.Writer) {

	var VoteList = make(map[string]voting.VotingInstance)
	var VoteSystem = impl.NewVotingSystem(nil, VoteList)
	var histoChoice = make([]voting.Choice, 0)
	var voters = make([]*voting.User, 0)

	// var randomNumOfUser, err = random.IntRange(15, 20)
	// if err != nil {
	// 	xerrors.Errorf(err.Error())
	// }

	// //Random creating of a user and adds it to the list of voters
	// for i := 0; i < randomNumOfUser; i++ {
	// 	var chooseType, err1 = random.IntRange(1, 101)
	// 	if err1 != nil {
	// 		xerrors.Errorf(err.Error())
	// 	}
	// 	switch {
	// 	case chooseType < 20:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.YesVoter, make([]*voting.User, 0))
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	case chooseType < 40:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NoVoter, make([]*voting.User, 0))
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)

	// 	case chooseType < 60:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.IndeciseVoter, make([]*voting.User, 0))
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	case chooseType < 90:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ThresholdVoter, make([]*voting.User, 0))
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	default:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, make([]*voting.User, 0))
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	}
	// }

	//Manually entering the number of each categories

	YesNumber := 10
	NoNumber := 10
	IndecisiveNumber := 10
	ThresholdNumber := 10
	NonResponsibleNumber := 10
	TotalNumber := NonResponsibleNumber + YesNumber + NoNumber + IndecisiveNumber + ThresholdNumber

	i := 0
	for i = 0; i < YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.YesVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for i = i; i < NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NoVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for i = i; i < IndecisiveNumber+NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.IndecisiveVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for i = i; i < ThresholdNumber+IndecisiveNumber+NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ThresholdVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for i = i; i < NonResponsibleNumber+ThresholdNumber+IndecisiveNumber+NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NonResponsibleVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}

	//filling the preference list for delegation
	for i, user := range voters {
		randomNumberOfPreference, err := random.IntRange(0, 4)
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomAction")
		}
		for j := 0; j < randomNumberOfPreference; j++ {
			//random index creation (must NOT be == to index of current user)
			randomDelegateToIndex, err := random.IntRange(0, len(voters))
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
			}
			checkAlreadyHaveInPreferenceList := func(randIndex int) bool {
				for _, u := range user.PreferenceDelegationList {
					if voters[randIndex] == u {
						return true
					}
				}
				return false
			}

			for randomDelegateToIndex == i || checkAlreadyHaveInPreferenceList(randomDelegateToIndex) {
				randomDelegateToIndex, err = random.IntRange(0, len(voters))
				if err != nil {
					fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
				}
				//randomDelegateToIndex = randIndex
			}

			user.PreferenceDelegationList = append(user.PreferenceDelegationList, voters[randomDelegateToIndex])
		}

		fmt.Print("The preference list of "+user.UserID, " is : ")
		for _, user := range user.PreferenceDelegationList {
			fmt.Print(user.UserID + ", ")
		}
		fmt.Println(" ")
	}

	findIndexInListOfUser := func(list []*voting.User, user *voting.User) int {
		for i, u := range list {
			if u == user {
				return i
			}
		}
		return -1
	}

	//construct the graph of delegation, run acyclic detection algorithm to detect cycles
	g := graph.New(TotalNumber)
	for i, u := range voters {
		if len(u.PreferenceDelegationList) > 0 {
			g.Add(i, findIndexInListOfUser(voters, u.PreferenceDelegationList[0]))
		}
	}
	IsThereCycle := !graph.Acyclic(g)
	if IsThereCycle {
		fmt.Println("There is a cycle !!!")
	} else {
		fmt.Println("There is no cycle :) ")
	}

	// var list = make([]int, 0)
	// list = append(list, 1, 2, 1)
	// fmt.Println("LIST :::::::  ", list)
	// fmt.Println("FIRST ELEMENT IS  ", list[0])
	// list = list[1:len(list)]
	// fmt.Println("NEWLIST ::::::: ->>> ", list)
	// fmt.Println("FIRST ELEMENT IS  ", list[0])

	var stateDelegation = make([](struct {
		curr     int
		test     int
		max_size int
	}), len(voters))
	for j, user := range voters {
		stateDelegation[j] = struct {
			curr     int
			test     int
			max_size int
		}{0, 0, len(user.PreferenceDelegationList) - 1}
	}

	fmt.Print(stateDelegation)

	for !graph.Acyclic(g) {
		//changer les preferences listes des personnes jusqu'a ne plus avoir de cycle
		for i, u := range voters {
			if stateDelegation[i].test < stateDelegation[i].max_size {
				g.Delete(i, findIndexInListOfUser(voters, u.PreferenceDelegationList[stateDelegation[i].curr]))
				g.Add(i, findIndexInListOfUser(voters, u.PreferenceDelegationList[stateDelegation[i].test+1]))
				stateDelegation[i].test = stateDelegation[i].test + 1
				if graph.Acyclic(g) {
					stateDelegation[i].curr = stateDelegation[i].test
					break
				} else {
					g.Delete(i, findIndexInListOfUser(voters, u.PreferenceDelegationList[stateDelegation[i].test]))
					g.Add(i, findIndexInListOfUser(voters, u.PreferenceDelegationList[stateDelegation[i].curr]))
				}
			}
		}
	}

	IsThereCycle2 := !graph.Acyclic(g)
	if IsThereCycle2 {
		fmt.Println("There is still a cycle !!!")
	} else {
		fmt.Println("There is now no cycle :) ")
		fmt.Print(stateDelegation)
	}

	//candidats
	var candidats = make([]*voting.Candidate, 3)

	//empty list of votes
	//var votes = make(map[string]voting.Choice)

	//creation of votingConfig
	voteConfig, err := impl.NewVotingConfig(voters, "Simulation 1", "Sunny day everyday ?", candidats, "YesOrNoQuestion")
	if err != nil {
		fmt.Println(err.Error())
	}

	//creation of the voting instance
	VoteInstance, err := VoteSystem.CreateAndAdd("Simulation01", voteConfig, "open")
	if err != nil {
		fmt.Println(err.Error())
	}

	randomVote := func(user *voting.User, i int) {
		randomAction, err := random.IntRange(1, 4)
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomAction")
		}

		if randomAction == 1 {
			//Delegation action
			quantity_to_deleg, err := impl.NewLiquid(float64(user.VotingPower))
			if err != nil {
				fmt.Println(err.Error(), "fail to do quantity to deleg")
			}
			if len(user.PreferenceDelegationList) > 0 {
				err = VoteInstance.DelegTo(user, user.PreferenceDelegationList[0], quantity_to_deleg)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", user.PreferenceDelegationList[0].UserID, "il était", user.TypeOfUser)
			} else {
				//random index creation (must NOT be == to index of current user)
				randomDelegateToIndex, err := random.IntRange(0, len(voters))
				if err != nil {
					fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
				}
				for ok := true; ok; ok = (randomDelegateToIndex == i) {
					randomDelegateToIndex, err = random.IntRange(0, len(voters))
					if err != nil {
						fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
					}
				}
				err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
				if err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", voters[randomDelegateToIndex].UserID, "il était ", user.TypeOfUser)
			}
		} else if randomAction == 2 {
			//Vote action

			quantity := user.VotingPower
			quantity_to_Vote, err := impl.NewLiquid(float64(quantity))
			if err != nil {
				fmt.Println(err.Error())
			}
			liquid_0, err := impl.NewLiquid(0)
			if err != nil {
				fmt.Println(err.Error())
			}

			choiceTab := make(map[string]voting.Liquid)

			if len(user.HistoryOfChoice) == 0 {
				yesOrNo, err := random.IntRange(1, 3)
				if err != nil {
					fmt.Println(err.Error(), "fail to do yesOrNo ")
				}

				if yesOrNo == 1 {
					choiceTab["yes"] = quantity_to_Vote
					choiceTab["no"] = liquid_0
				} else {
					choiceTab["yes"] = liquid_0
					choiceTab["no"] = quantity_to_Vote
				}
			} else if user.HistoryOfChoice[0].VoteValue["no"].Percentage != 0. {
				choiceTab["yes"] = liquid_0
				choiceTab["no"] = quantity_to_Vote
			} else {
				choiceTab["yes"] = quantity_to_Vote
				choiceTab["no"] = liquid_0
			}

			//quantity to vote

			//create choice
			choice, err := impl.NewChoice(choiceTab)
			if err != nil {
				fmt.Println(err.Error())
			}

			//set the choice
			err = VoteInstance.SetVote(user, choice)
			if err != nil {
				fmt.Println(err.Error())
			}

			//cast the vote
			// err = VoteInstance.CastVote(user)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }

			fmt.Println(user.UserID, " a voté pour ", quantity, "%", "il était", user.TypeOfUser)
		}
	}

	yesVote := func(user *voting.User, votingPower float64) {
		quantity := votingPower
		quantity_to_Vote, err := impl.NewLiquid(float64(quantity))
		if err != nil {
			fmt.Println(err.Error())
		}
		liquid_0, err := impl.NewLiquid(0)
		if err != nil {
			fmt.Println(err.Error())
		}

		choiceTab := make(map[string]voting.Liquid)

		choiceTab["yes"] = quantity_to_Vote
		choiceTab["no"] = liquid_0
		//create choice
		choice, err := impl.NewChoice(choiceTab)
		if err != nil {
			fmt.Println(err.Error())
		}

		//set the choice
		err = VoteInstance.SetVote(user, choice)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(user.UserID, " a voté pour ", quantity, "%", "il était", user.TypeOfUser)
	}
	noVote := func(user *voting.User, votingPower float64) {
		quantity := votingPower
		quantity_to_Vote, err := impl.NewLiquid(float64(quantity))
		if err != nil {
			fmt.Println(err.Error())
		}
		liquid_0, err := impl.NewLiquid(0)
		if err != nil {
			fmt.Println(err.Error())
		}

		choiceTab := make(map[string]voting.Liquid)

		choiceTab["no"] = quantity_to_Vote
		choiceTab["yes"] = liquid_0
		//create choice
		choice, err := impl.NewChoice(choiceTab)
		if err != nil {
			fmt.Println(err.Error())
		}

		//set the choice
		err = VoteInstance.SetVote(user, choice)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(user.UserID, " a voté pour ", quantity, "%", "il était", user.TypeOfUser)
	}

	IndecisiveVote := func(user *voting.User, i int) {
		//Delegation action
		quantity_to_deleg, err := impl.NewLiquid(float64(user.VotingPower))
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to deleg")
		}
		if len(user.PreferenceDelegationList) > 0 {
			err = VoteInstance.DelegTo(user, user.PreferenceDelegationList[0], quantity_to_deleg)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", user.PreferenceDelegationList[0].UserID, "il était", user.TypeOfUser)
		} else {
			//random index creation (must NOT be == to index of current user)
			randomDelegateToIndex, err := random.IntRange(0, len(voters))
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
			}
			for ok := true; ok; ok = (randomDelegateToIndex == i) {
				randomDelegateToIndex, err = random.IntRange(0, len(voters))
				if err != nil {
					fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
				}
			}
			err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", voters[randomDelegateToIndex].UserID, "il était ", user.TypeOfUser)
		}
	}

	ThresholdVote := func(user *voting.User, i int, threshold int) {

		var thresholdCOmparator = user.VotingPower

		if thresholdCOmparator > float64(threshold) {
			//Delegation action
			quantity_to_deleg, err := impl.NewLiquid(float64(user.VotingPower))
			if err != nil {
				fmt.Println(err.Error(), "fail to do quantity to deleg")
			}
			if len(user.PreferenceDelegationList) > 0 {
				err = VoteInstance.DelegTo(user, user.PreferenceDelegationList[0], quantity_to_deleg)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", user.PreferenceDelegationList[0].UserID, "il était", user.TypeOfUser)
			} else {
				//random index creation (must NOT be == to index of current user)
				randomDelegateToIndex, err := random.IntRange(0, len(voters))
				if err != nil {
					fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
				}
				for ok := true; ok; ok = (randomDelegateToIndex == i) {
					randomDelegateToIndex, err = random.IntRange(0, len(voters))
					if err != nil {
						fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
					}
				}
				err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
				if err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", voters[randomDelegateToIndex].UserID, "il était ", user.TypeOfUser)
			}

		} else {
			//Vote action

			quantity := user.VotingPower
			quantity_to_Vote, err := impl.NewLiquid(float64(quantity))
			if err != nil {
				fmt.Println(err.Error())
			}
			liquid_0, err := impl.NewLiquid(0)
			if err != nil {
				fmt.Println(err.Error())
			}

			choiceTab := make(map[string]voting.Liquid)

			if len(user.HistoryOfChoice) == 0 {
				yesOrNo, err := random.IntRange(1, 3)
				if err != nil {
					fmt.Println(err.Error(), "fail to do yesOrNo ")
				}

				if yesOrNo == 1 {
					choiceTab["yes"] = quantity_to_Vote
					choiceTab["no"] = liquid_0
				} else {
					choiceTab["yes"] = liquid_0
					choiceTab["no"] = quantity_to_Vote
				}
			} else if user.HistoryOfChoice[0].VoteValue["no"].Percentage != 0. {
				choiceTab["yes"] = liquid_0
				choiceTab["no"] = quantity_to_Vote
			} else {
				choiceTab["yes"] = quantity_to_Vote
				choiceTab["no"] = liquid_0
			}

			//quantity to vote

			//create choice
			choice, err := impl.NewChoice(choiceTab)
			if err != nil {
				fmt.Println(err.Error())
			}

			//set the choice
			err = VoteInstance.SetVote(user, choice)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(user.UserID, " a voté pour ", quantity, "%", "il était", user.TypeOfUser)
		}
	}

	NonResponsibleVoter := func(user *voting.User, i int) {
		if len(user.HistoryOfChoice) == 0 {
			var randomNumberToChooseYesOrNo, err = random.IntRange(0, 2)
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
			}
			if randomNumberToChooseYesOrNo == 0 {
				yesVote(user, 100.)
			} else {
				noVote(user, 100.)
			}
		} else {
			//Delegation action
			IndecisiveVote(user, i)
		}
	}

	for ok := true; ok; ok = VoteInstance.CheckVotingPowerOfVoters() {
		for i, user := range VoteInstance.GetConfig().Voters {

			if user.VotingPower > 0 {
				switch user.TypeOfUser {
				case voting.YesVoter:
					yesVote(user, user.VotingPower)
				case voting.NoVoter:
					noVote(user, user.VotingPower)
				case voting.IndecisiveVoter:
					IndecisiveVote(user, i)
				case voting.ThresholdVoter:
					var threshold = 600
					ThresholdVote(user, i, threshold)
				case voting.NonResponsibleVoter:
					NonResponsibleVoter(user, i)
				case voting.None:
					randomVote(user, i)
				}
			}
		}
	}

	counterYesVoter := 0
	counterNoVoter := 0
	counterIndecisiveVoter := 0
	counterThresholdVoter := 0
	counterNormalVoter := 0
	for _, user := range VoteInstance.GetConfig().Voters {
		//fmt.Println("Voting power of ", user.UserID, " = ", user.VotingPower, "il était de type", user.TypeOfUser)
		if user.TypeOfUser == "YesVoter" {
			counterYesVoter++
		} else if user.TypeOfUser == "NoVoter" {
			counterNoVoter++
		} else if user.TypeOfUser == "IndeciseVoter" {
			counterIndecisiveVoter++
		} else if user.TypeOfUser == "ThresholdVoter" {
			counterThresholdVoter++
		} else {
			counterNormalVoter++
		}
	}
	fmt.Println("There is", counterYesVoter, "yesVoter,", counterNoVoter, "noVoter,", counterThresholdVoter, "Threshold Voter", counterIndecisiveVoter, "IndecisiveVoter and ", counterNormalVoter, "normalVoter")

	for _, user := range VoteInstance.GetConfig().Voters {
		fmt.Println("Voting power of ", user.UserID, " = ", user.VotingPower)
	}

	results := VoteInstance.GetResults()
	s := "%"

	fmt.Fprintf(out, "digraph network_activity {\n")
	fmt.Fprintf(out, "labelloc=\"t\";")
	fmt.Fprintf(out, "label = <Votation Diagram of %d nodes.    Results are Yes = %.4v %s, No = %.4v %s<font point-size='10'><br/>(generated %s)</font>>;", len(voters)+2, results["yes"], s, results["no"], s, time.Now().Format("2 Jan 06 - 15:04:05"))
	fmt.Fprintf(out, "graph [fontname = \"helvetica\"];")
	fmt.Fprintf(out, "node [fontname = \"helvetica\" area = 10 style= filled];")
	for j, user := range VoteInstance.GetConfig().Voters {
		colorOfUser := "#FFFFFF"
		if user.TypeOfUser == "YesVoter" { //YesVoter
			colorOfUser = "#42D03F"
		} else if user.TypeOfUser == "NoVoter" { //NoVoter
			colorOfUser = "#FC5A5A"
		} else if user.TypeOfUser == "IndecisiveVoter" { //IndecisiveVoter
			colorOfUser = "#B7FCFF"
		} else if user.TypeOfUser == "ThresholdVoter" { //ThresholdVoter
			colorOfUser = "#6BA7E8"
		} else if user.TypeOfUser == "NonResponsibleVoter" { //NonResponsibleVoter
			colorOfUser = "#6066D3"
		} else if user.TypeOfUser == "ResponsibleVoter" { //ResponsibleVoter
			colorOfUser = "#111111"
		} else { //NormalVoter
			colorOfUser = "#FFFFFF"
		}
		s := strconv.FormatInt(int64(j), 10)
		fmt.Fprintf(out, "user%s [filledcolor=%s];\n", s, colorOfUser)
	}

	fmt.Fprintf(out, "edge [fontname = \"helvetica\"];\n")

	for _, user := range VoteInstance.GetConfig().Voters {

		colorVoteYes := "#22bd27"
		colorVoteNo := "#cf1111"
		colorDeleg := "#8A2BE2"

		//creation d'un tableau qui a les cumulative values (plus simple pour le graph)
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

		//creation of the arrows for the votes
		for _, choice := range cumulativeHistoryOfChoice {
			if choice.VoteValue["yes"].Percentage != 0. {
				fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
					"[ label = < <font color='#22bd27'><b>%v</b></font><br/>> color=\"%s\" penwidth=%v];\n",
					user.UserID, "YES", choice.VoteValue["yes"].Percentage, colorVoteYes, choice.VoteValue["yes"].Percentage/40)
			}

			if choice.VoteValue["no"].Percentage != 0. {
				fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
					"[ label = < <font color='#cf1111'><b>%v</b></font><br/>> color=\"%s\" penwidth=%v];\n",
					user.UserID, "NO", choice.VoteValue["no"].Percentage, colorVoteNo, choice.VoteValue["no"].Percentage/40)
			}
		}

		for other, quantity := range user.DelegatedTo {
			fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
				"[ label = < <font color='#8A2BE2'><b>%v</b></font><br/>> color=\"%s\" penwidth=%v];\n",
				user.UserID, other, quantity.Percentage, colorDeleg, quantity.Percentage/40)
		}
	}

	fmt.Fprintf(out, "}\n")

}

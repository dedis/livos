package simulation

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"github.com/mazen160/go-random"
	"golang.org/x/xerrors"
)

// GenerateItemsGraphviz creates a graphviz representation of the items. One can
// generate a graphical representation with `dot -Tpdf graph.dot -o graph.pdf`
func Simulation2(out io.Writer) {

	const InitialVotingPower = 100.

	var VoteList = make(map[string]voting.VotingInstance)
	var VoteSystem = impl.NewVotingSystem(nil, VoteList)
	var histoChoice = make([]voting.Choice, 0)

	var randomNumOfUser, err = random.IntRange(98, 101)
	if err != nil {
		xerrors.Errorf(err.Error())
	}

	//Random creating of a user and adds it to the list of voters
	var voters = make([]*voting.User, 0)
	for i := 0; i < randomNumOfUser; i++ {
		var chooseType, err1 = random.IntRange(1, 101)
		if err1 != nil {
			xerrors.Errorf(err.Error())
		}
		switch {
		case chooseType < 5:
			var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.YesVoter, nil)
			if err != nil {
				xerrors.Errorf(err.Error())
			}
			voters = append(voters, &user)
		case chooseType < 10:
			var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NoVoter, nil)
			if err != nil {
				xerrors.Errorf(err.Error())
			}
			voters = append(voters, &user)

		case chooseType < 60:
			var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.IndecisiveVoter, nil)
			if err != nil {
				xerrors.Errorf(err.Error())
			}
			voters = append(voters, &user)
		case chooseType < 80:
			var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ThresholdVoter, nil)
			if err != nil {
				xerrors.Errorf(err.Error())
			}
			voters = append(voters, &user)
		case chooseType < 90:
			var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NonResponsibleVoter, nil)
			if err != nil {
				xerrors.Errorf(err.Error())
			}
			voters = append(voters, &user)
		default:
			var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
			if err != nil {
				xerrors.Errorf(err.Error())
			}
			voters = append(voters, &user)
		}

	}

	//candidats
	var candidats = make([]string, 3)

	//empty list of votes
	//var votes = make(map[string]voting.Choice)

	//creation of votingConfig
	voteConfig, err := impl.NewVotingConfig(voters, "Simulation 1", "Sunny day everyday ?", candidats)
	if err != nil {
		fmt.Println(err.Error())
	}

	//creation of the voting instance
	VoteInstance, err := VoteSystem.CreateAndAdd("Simulation01", voteConfig, "open")
	if err != nil {
		fmt.Println(err.Error())
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
		quantity_to_deleg, err := impl.NewLiquid(float64(user.VotingPower))
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to deleg")
		}
		err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", voters[randomDelegateToIndex].UserID, "il était", user.TypeOfUser)
	}
	randomVote := func(user *voting.User, i int) {
		randomAction, err := random.IntRange(1, 4)
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomAction")
		}

		if randomAction == 1 {
			//Delegation action

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
			randomQuantityToDelegate, err := random.IntRange(1, int(user.VotingPower/10)+1)
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomQuantityToDelegate")
			}
			randomQuantityToDelegate *= 10
			quantity_to_deleg, err := impl.NewLiquid(float64(randomQuantityToDelegate))
			if err != nil {
				fmt.Println(err.Error(), "fail to do quantity to deleg")
			}
			err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", voters[randomDelegateToIndex].UserID, "il était", user.TypeOfUser)

		} else if randomAction == 2 {
			//Vote action

			quantity := user.VotingPower
			var randomNumberToChooseYesOrNo, err = random.IntRange(0, 2)
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
			}
			if randomNumberToChooseYesOrNo == 0 {
				yesVote(user, quantity)
			} else {
				noVote(user, quantity)
			}

		}
	}

	ThresholdVote := func(user *voting.User, i int, threshold int) {

		var thresholdComparator = 0.
		for i := range user.HistoryOfChoice {
			thresholdComparator += user.HistoryOfChoice[i].VoteValue["yes"].Percentage
			thresholdComparator += user.HistoryOfChoice[i].VoteValue["no"].Percentage
		}

		if thresholdComparator > float64(threshold) {
			//Delegation action
			IndecisiveVote(user, i)

		} else {
			//Vote action

			quantity := user.VotingPower

			if len(user.HistoryOfChoice) == 0 {
				yesOrNo, err := random.IntRange(1, 3)
				if err != nil {
					fmt.Println(err.Error(), "fail to do yesOrNo ")
				}

				if yesOrNo == 1 {
					yesVote(user, quantity)
				} else {
					noVote(user, quantity)
				}
			} else if user.HistoryOfChoice[0].VoteValue["no"].Percentage != 0. {
				noVote(user, quantity)
			} else {
				yesVote(user, quantity)
			}
		}
	}
	NonResponsibleVoter := func(user *voting.User, i int) {
		if len(user.HistoryOfChoice) == 0 {
			var randomNumberToChooseYesOrNo, err = random.IntRange(0, 2)
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
			}
			if randomNumberToChooseYesOrNo == 0 {
				yesVote(user, InitialVotingPower)
			} else {
				noVote(user, InitialVotingPower)
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
	counterNonResponsibleVoter := 0
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
		} else if user.TypeOfUser == "NonResponsibleVoter" {
			counterNonResponsibleVoter++
		} else {
			counterNormalVoter++
		}
	}
	fmt.Println("There is", counterYesVoter, "yesVoter,", counterNoVoter, "noVoter,", counterThresholdVoter, "Threshold Voter", counterIndecisiveVoter, "IndecisiveVoter and ", counterNormalVoter, "normalVoter")

	results := VoteInstance.GetResults()
	s := "%"
	fmt.Fprintf(out, "digraph network_activity {\n")
	fmt.Fprintf(out, "labelloc=\"t\";")
	fmt.Fprintf(out, "label = <Votation Diagram of %d nodes.    Results are Yes = %.4v %s, No = %.4v %s<font point-size='10'><br/>(generated: %s)<br/> Il y a %v YesVoter, %v NoVoter, %v Threshold Voters, %v Non responsibleVoter, %v IndecisiveVoter and %v NormalVoter</font>>; ", len(voters)+2, results["yes"], s, results["no"], s, time.Now(), counterYesVoter, counterNoVoter, counterThresholdVoter, counterNonResponsibleVoter, counterIndecisiveVoter, counterNormalVoter)
	fmt.Fprintf(out, "graph [fontname = \"helvetica\"];")
	fmt.Fprintf(out, "node [fontname = \"helvetica\" area = 10 fillcolor=gold];")
	fmt.Fprintf(out, "edge [fontname = \"helvetica\"];\n")

	for _, user := range VoteInstance.GetConfig().Voters {

		colorVoteYes := "#22bd27"
		colorVoteNo := "#cf1111"
		colorDeleg := "#8A2BE2"

		/* colorOfUser := "#FFFFFF"
		if user.TypeOfUser == 0 { //YesVoter
			colorOfUser = "#42D03F"
		} else if user.TypeOfUser == 1 { //NoVoter
			colorOfUser = "#FC5A5A"
		} else if user.TypeOfUser == 2 { //IndecisiveVoter
			colorOfUser = "#B7FCFF"
		} else if user.TypeOfUser == 3 { //ThresholdVoter
			colorOfUser = "#6BA7E8"
		} else if user.TypeOfUser == 4 { //NonResponsibleVoter
			colorOfUser = "#6066D3"
		} else { //NormalVoter
			colorOfUser = "#FFFFFF"
		} */

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
					user.UserID, "YES", choice.VoteValue["yes"].Percentage, colorVoteYes, choice.VoteValue["yes"].Percentage/60)
			}

			if choice.VoteValue["no"].Percentage != 0. {
				fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
					"[ label = < <font color='#cf1111'><b>%v</b></font><br/>> color=\"%s\" penwidth=%v];\n",
					user.UserID, "NO", choice.VoteValue["no"].Percentage, colorVoteNo, choice.VoteValue["no"].Percentage/60)
			}
		}

		for other, quantity := range user.DelegatedTo {
			fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
				"[ label = < <font color='#8A2BE2'><b>%v</b></font><br/>> color=\"%s\" penwidth=%v];\n",
				user.UserID, other, quantity.Percentage, colorDeleg, quantity.Percentage/60)
		}
	}

	fmt.Fprintf(out, "}\n")

}

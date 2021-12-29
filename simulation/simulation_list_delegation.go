package simulation

import (
	"fmt"
	"io"
	"strconv"

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

	YesNumber := 5
	NoNumber := 5
	IndecisiveNumber := 5
	ThresholdNumber := 5
	NonResponsibleNumber := 5
	ResponsibleNumber := 5
	TotalNumber := ResponsibleNumber + NonResponsibleNumber + YesNumber + NoNumber + IndecisiveNumber + ThresholdNumber

	i := 0
	for i = 0; i < YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.YesVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NoVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < IndecisiveNumber+NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.IndecisiveVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < ThresholdNumber+IndecisiveNumber+NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ThresholdVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < NonResponsibleNumber+ThresholdNumber+IndecisiveNumber+NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NonResponsibleVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < ResponsibleNumber+NonResponsibleNumber+ThresholdNumber+IndecisiveNumber+NoNumber+YesNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ResponsibleVoter, make([]*voting.User, 0))
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

	for ok := true; ok; ok = VoteInstance.CheckVotingPowerOfVoters() {
		for i, user := range VoteInstance.GetConfig().Voters {

			if user.VotingPower > 0 {
				switch user.TypeOfUser {
				case voting.YesVoter:
					VoteInstance.YesVote(user, user.VotingPower)
				case voting.NoVoter:
					VoteInstance.NoVote(user, user.VotingPower)
				case voting.IndecisiveVoter:
					VoteInstance.IndecisiveVote(user, i, user.VotingPower)
				case voting.ThresholdVoter:
					var threshold = 600
					VoteInstance.ThresholdVote(user, i, threshold)
				case voting.NonResponsibleVoter:
					VoteInstance.NonResponsibleVote(user, i)
				case voting.ResponsibleVoter:
					VoteInstance.ResponsibleVote(user, i)
				case voting.None:
					VoteInstance.RandomVote(user, i)
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

	VoteInstance.ConstructTextForGraph(out)

}

package simulation

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"golang.org/x/xerrors"
)

// GenerateItemsGraphviz creates a graphviz representation of the items. One can
// generate a graphical representation with `dot -Tpdf graph.dot -o graph.pdf`
func Simulation_RealData_Candidats(out_liquid io.Writer, out_normal io.Writer) float64 {

	const InitialVotingPower = 100

	var VoteList = make(map[string]voting.VotingInstance)
	var VoteSystem = impl.NewVotingSystem(nil, VoteList)
	var histoChoice = make([]voting.Choice, 0)

	var voters = make([]*voting.User, 0)

	// var randomNumOfUser, err = random.IntRange(10, 11)
	// if err != nil {
	// 	xerrors.Errorf(err.Error())
	// }

	//Random creating of a user and adds it to the list of voters

	// for i := 0; i < randomNumOfUser; i++ {
	// 	var chooseType, err1 = random.IntRange(1, 101)
	// 	if err1 != nil {
	// 		xerrors.Errorf(err.Error())
	// 	}
	// 	switch {
	// 	case chooseType < 5:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.YesVoter, nil)
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	case chooseType < 50:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.IndecisiveVoter, nil)
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	case chooseType < 70:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ThresholdVoter, nil)
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	case chooseType < 80:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NonResponsibleVoter, nil)
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	case chooseType < 90:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ResponsibleVoter, nil)
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	default:
	// 		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
	// 		if err != nil {
	// 			xerrors.Errorf(err.Error())
	// 		}
	// 		voters = append(voters, &user)
	// 	}

	// }

	//Manually entering the number of each categories

	const MULTIPLICATOR = 50

	//Realist data
	// CandVoterNumber := 23 * MULTIPLICATOR
	// IndecisiveNumber := 20 * MULTIPLICATOR
	// ThresholdNumber := 16 * MULTIPLICATOR
	// NonResponsibleNumber := 1 * MULTIPLICATOR
	// ResponsibleNumber := 40 * MULTIPLICATOR

	//Realist Data (without indecisive)
	// CandVoterNumber := 29 * MULTIPLICATOR
	// IndecisiveNumber := 0
	// ThresholdNumber := 20 * MULTIPLICATOR
	// NonResponsibleNumber := 1 * MULTIPLICATOR
	// ResponsibleNumber := 50 * MULTIPLICATOR

	//data for graphic simulation
	CandVoterNumber := 3
	IndecisiveNumber := 4
	ThresholdNumber := 3
	NonResponsibleNumber := 1
	ResponsibleNumber := 4

	//TotalNumber := NonResponsibleNumber + YesNumber + NoNumber + IndecisiveNumber + ThresholdNumber

	i := 0
	for i = 0; i < CandVoterNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.CandVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < IndecisiveNumber+CandVoterNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.IndecisiveVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < ThresholdNumber+IndecisiveNumber+CandVoterNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ThresholdVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < NonResponsibleNumber+ThresholdNumber+IndecisiveNumber+CandVoterNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.NonResponsibleVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}
	for ; i < ResponsibleNumber+NonResponsibleNumber+ThresholdNumber+IndecisiveNumber+CandVoterNumber; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.ResponsibleVoter, make([]*voting.User, 0))
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}

	//fmt.Println("voters list is : ", voters)

	//candidats
	var candidatIC, _ = VoteSystem.NewCandidate("IC")
	var candidatSV, _ = VoteSystem.NewCandidate("SV")
	var candidatMECA, _ = VoteSystem.NewCandidate("MECA")
	var candidatMICRO, _ = VoteSystem.NewCandidate("MICRO")

	var candidats = []*voting.Candidate{&candidatIC, &candidatSV, &candidatMECA, &candidatMICRO}

	//empty list of votes
	//var votes = make(map[string]voting.Choice)

	//creation of votingConfig
	voteConfig, err := impl.NewVotingConfig(voters, "Simulation 1", "Which EPFL faculty is the best ?", candidats, "CandidateQuestion")
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
				case voting.CandVoter:
					VoteInstance.CandidateVote(user, i, user.VotingPower)
				case voting.IndecisiveVoter:
					VoteInstance.IndecisiveVote(user, i, user.VotingPower)
				case voting.ThresholdVoter:
					var thresholdPourcentage, err = impl.GenerateRandomThreshold()
					if err != nil {
						fmt.Println(err.Error())
					}
					var threshold int = thresholdPourcentage * len(voters)
					VoteInstance.ThresholdVoteCandidate(user, i, threshold, user.VotingPower)
				case voting.NonResponsibleVoter:
					VoteInstance.NonResponsibleVoteCandidate(user, i, user.VotingPower)
				case voting.ResponsibleVoter:
					VoteInstance.ResponsibleVoteCandidate(user, i, user.VotingPower)
				case voting.None:
					VoteInstance.DefaultVoteCandidate(user, i)
				}
			}
		}
	}

	//Counters for the repartition information
	counterCandidateVoter := 0
	counterIndecisiveVoter := 0
	counterThresholdVoter := 0
	counterNormalVoter := 0
	counterNonResponsibleVoter := 0
	counterResponsibleVoter := 0
	for _, user := range voters {
		//fmt.Println("Voting power of ", user.UserID, " = ", user.VotingPower, "il était de type", user.TypeOfUser)
		if user.TypeOfUser == "CandVoter" {
			counterCandidateVoter++
		} else if user.TypeOfUser == "IndecisiveVoter" {
			counterIndecisiveVoter++
		} else if user.TypeOfUser == "ThresholdVoter" {
			counterThresholdVoter++
		} else if user.TypeOfUser == "NonResponsibleVoter" {
			counterNonResponsibleVoter++
		} else if user.TypeOfUser == "ResponsibleVoter" {
			counterResponsibleVoter++
		} else {
			counterNormalVoter++
		}
	}

	fmt.Println("There is ", counterCandidateVoter, "CandidateVoter,", counterThresholdVoter, "Threshold Voter,", counterNonResponsibleVoter, "NonresponsibleVoter,", counterResponsibleVoter, "ResponsibleVoter,", counterIndecisiveVoter, "IndecisiveVoter and", counterNormalVoter, "normalVoter")

	VoteInstance.ConstructTextForGraphCandidates(out_liquid, VoteInstance.GetResults())

	LiquidResults := VoteInstance.GetResults()

	//modify the intern choices of the users (history of choice) to replace them with non liquid decisions.

	//adding the "new" choice that is BLANK CHOICE
	var blankCandidate, _ = VoteSystem.NewCandidate("Blank")
	new_candidates := append(candidats, &blankCandidate)

	VoteInstance.SetConfig(VoteInstance.GetConfig().SetCandidates(new_candidates))

	for _, user := range voters {
		//creation fo the liquid 100 that will be needed always
		liquid_100, err := impl.NewLiquid(InitialVotingPower)
		if err != nil {
			fmt.Println(err.Error())
		}

		if len(user.HistoryOfChoice) != 0 {
			//he voted for at least 1 person

			//build the cumulative of the user to help find the max
			//creation d'un tableau qui a les cumulative values (plus simple pour le graph)
			cumulativeHistoryOfChoice := make([]voting.Choice, 0)
			new_vote_value := make(map[string]voting.Liquid)
			for _, choice := range user.HistoryOfChoice {
				for name, value := range choice.VoteValue {
					var err error
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

			//find the max amongs the candidats (in max name)
			var max_name string
			var max_value int = 0
			for _, choiceCumul := range cumulativeHistoryOfChoice {
				for name, value := range choiceCumul.VoteValue {
					if value.Percentage > max_value {
						max_value = value.Percentage
						max_name = name
					}
				}
			}

			//we construct new historyOfChoice with 100 given to the max_name candidate (the prefered one)
			new_vote_value = make(map[string]voting.Liquid)
			new_vote_value[max_name] = liquid_100

			new_choice, err = impl.NewChoice(new_vote_value)
			if err != nil {
				fmt.Println(err.Error())
			}
			new_HistoryOfChoice := make([]voting.Choice, 0)
			new_HistoryOfChoice = append(new_HistoryOfChoice, new_choice)

			user.HistoryOfChoice = new_HistoryOfChoice
		} else {
			//the case where the user only delegated : we construct new historyOfChoice with 100 given to the BLANK_candidate

			new_vote_value := make(map[string]voting.Liquid)
			new_vote_value[blankCandidate.CandidateID] = liquid_100

			new_choice, err := impl.NewChoice(new_vote_value)
			if err != nil {
				fmt.Println(err.Error())
			}
			new_HistoryOfChoice := make([]voting.Choice, 0)
			new_HistoryOfChoice = append(new_HistoryOfChoice, new_choice)

			user.HistoryOfChoice = new_HistoryOfChoice
		}

		//need to erase the delegatedTo and delegatedFrom
		EmptyDelegatadMap := make(map[string]voting.Liquid, 0)
		user.DelegatedTo = EmptyDelegatadMap
		user.DelegatedFrom = EmptyDelegatadMap
	}

	//at this state all the historyOfChoice of the users are modified (simplified) for the normal version.
	//When calling GetResults we obtain results that are traditionnal
	VoteInstance.ConstructTextForGraphCandidates(out_normal, VoteInstance.GetResults())

	NormalResults := VoteInstance.GetResults()

	totalSumOfDifference := 0.
	fmt.Println("==========================================")
	fmt.Println("DIFFERENCES PRECISION : ")
	for _, cand := range candidats {
		difference := math.Abs(LiquidResults[cand.CandidateID] - NormalResults[cand.CandidateID])
		totalSumOfDifference += difference
		fmt.Println(difference, " => ", cand.CandidateID)
	}
	fmt.Println("Total difference : ", totalSumOfDifference)
	fmt.Println("==========================================")

	return totalSumOfDifference
}

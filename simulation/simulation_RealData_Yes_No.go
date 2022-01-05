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
func Simulation_RealData_Yes_No(out_liquid io.Writer, out_normal io.Writer) float64 {

	const InitialVotingPower = 100

	var VoteList = make(map[string]voting.VotingInstance)
	var VoteSystem = impl.NewVotingSystem(nil, VoteList)
	var histoChoice = make([]voting.Choice, 0)
	var voters = make([]*voting.User, 0)

	const MULTIPLICATOR = 100

	// YesNumber := 11 * MULTIPLICATOR
	// NoNumber := 11 * MULTIPLICATOR
	// IndecisiveNumber := 21 * MULTIPLICATOR
	// ThresholdNumber := 16 * MULTIPLICATOR
	// NonResponsibleNumber := 1 * MULTIPLICATOR
	// ResponsibleNumber := 40 * MULTIPLICATOR

	YesNumber := 15 * MULTIPLICATOR
	NoNumber := 15 * MULTIPLICATOR
	IndecisiveNumber := 0 * MULTIPLICATOR
	ThresholdNumber := 19 * MULTIPLICATOR
	NonResponsibleNumber := 1 * MULTIPLICATOR
	ResponsibleNumber := 50 * MULTIPLICATOR

	//TotalNumber := NonResponsibleNumber + YesNumber + NoNumber + IndecisiveNumber + ThresholdNumber

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

	//candidats
	var candidats = make([]*voting.Candidate, 3)

	//empty list of votes
	//var votes = make(map[string]voting.Choice)

	//creation of votingConfig
	voteConfig, err := impl.NewVotingConfig(voters, "Simulation01", "Sunny day everyday ?", candidats, "YesOrNoQuestion")
	if err != nil {
		fmt.Println(err.Error())
	}

	//creation of the voting instance
	VoteInstance, err := VoteSystem.CreateAndAdd("Simulation01", voteConfig, "open")
	if err != nil {
		fmt.Println(err.Error())
	}

	//Proceed to the votation rounds until everyone has spend all of their voting power
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
					var thresholdPourcentage, err = impl.GenerateRandomThreshold()
					if err != nil {
						fmt.Println(err.Error())
					}
					var threshold int = thresholdPourcentage * len(voters)
					VoteInstance.ThresholdVote(user, i, threshold, user.VotingPower)
				case voting.NonResponsibleVoter:
					VoteInstance.NonResponsibleVote(user, i, user.VotingPower)
				case voting.ResponsibleVoter:
					VoteInstance.ResponsibleVote(user, i, user.VotingPower)
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
	counterNonResponsibleVoter := 0
	counterResponsibleVoter := 0
	for _, user := range VoteInstance.GetConfig().Voters {
		//fmt.Println("Voting power of ", user.UserID, " = ", user.VotingPower, "il était de type", user.TypeOfUser)
		if user.TypeOfUser == "YesVoter" {
			counterYesVoter++
		} else if user.TypeOfUser == "NoVoter" {
			counterNoVoter++
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
	fmt.Println("There is ", counterYesVoter, "yesVoter,", counterNoVoter, "noVoter,", counterThresholdVoter, "Threshold Voter,", counterNonResponsibleVoter, "NonresponsibleVoter,", counterResponsibleVoter, "ResponsibleVoter,", counterIndecisiveVoter, "IndecisiveVoter and", counterNormalVoter, "normalVoter")

	//call de la fonction qui ecrit toutes les infos dans le fichier texte.
	VoteInstance.ConstructTextForGraph(out_liquid)
	LiquidResults := VoteInstance.GetResults()

	///////////////////////////////////////////////////////////////////////////////////////////////////

	VoteInstance.SetConfig(VoteInstance.GetConfig())

	for _, user := range VoteInstance.GetConfig().Voters {
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

			//we construct new historyOfChoice with 100 given to the max_name yes/no (the prefered one)
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
			//the case where the user only delegated : we construct new historyOfChoice with 100 given to the blank voting

			new_vote_value := make(map[string]voting.Liquid)
			new_vote_value["blank"] = liquid_100

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
	VoteInstance.ConstructTextForGraph(out_normal)

	NormalResults := VoteInstance.GetResults()

	totalSumOfDifference := 0.
	fmt.Println("==========================================")
	fmt.Println("DIFFERENCES PRECISION : ")
	difference_Yes := math.Abs(LiquidResults["yes"] - NormalResults["yes"])
	difference_No := math.Abs(LiquidResults["no"] - NormalResults["no"])
	fmt.Println(difference_Yes, " => ", "yes")
	fmt.Println(difference_No, " => ", "no")
	totalSumOfDifference = difference_No + difference_Yes
	fmt.Println("Total difference : ", totalSumOfDifference)
	fmt.Println("==========================================")

	return totalSumOfDifference
}

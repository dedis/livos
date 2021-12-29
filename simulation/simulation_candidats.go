package simulation

import (
	"fmt"
	"io"
	"strconv"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"golang.org/x/xerrors"
)

// GenerateItemsGraphviz creates a graphviz representation of the items. One can
// generate a graphical representation with `dot -Tpdf graph.dot -o graph.pdf`
func Simulation_candidats(out io.Writer) {

	const InitialVotingPower = 100.

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

	YesNumber := 2
	NoNumber := 2
	IndecisiveNumber := 5
	ThresholdNumber := 3
	NonResponsibleNumber := 2
	ResponsibleNumber := 6
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

	//candidats inputs
	var candidatTrump, _ = VoteSystem.NewCandidate("Trump")
	var candidatObama, _ = VoteSystem.NewCandidate("Obama")
	var candidatJeanMi, _ = VoteSystem.NewCandidate("JeanMi")
	var candidatMacron, _ = VoteSystem.NewCandidate("Macron")

	var candidats = []*voting.Candidate{&candidatObama, &candidatTrump, &candidatJeanMi, &candidatMacron}

	//empty list of votes
	//var votes = make(map[string]voting.Choice)

	//creation of votingConfig
	voteConfig, err := impl.NewVotingConfig(voters, "Simulation 1", "Who are you gonna elect as a President ?", candidats, "CandidateQuestion")
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
					VoteInstance.CandidateVote(user, i, user.VotingPower)
				case voting.IndecisiveVoter:
					VoteInstance.IndecisiveVoteCandidate(user, i, user.VotingPower)
				case voting.ThresholdVoter:
					var threshold = 600
					VoteInstance.ThresholdVoteCandidate(user, i, threshold, user.VotingPower)
				case voting.NonResponsibleVoter:
					VoteInstance.NonResponsibleVoteCandidate(user, i, user.VotingPower)
				case voting.ResponsibleVoter:
					VoteInstance.ResponsibleVoteCandidate(user, i, user.VotingPower)
				case voting.None:
					VoteInstance.RandomVoteCandidate(user, i)
				}
			}
		}
	}

	counterYesVoter := 0
	counterIndecisiveVoter := 0
	counterThresholdVoter := 0
	counterNormalVoter := 0
	counterNonResponsibleVoter := 0
	counterResponsibleVoter := 0
	for _, user := range voters {
		//fmt.Println("Voting power of ", user.UserID, " = ", user.VotingPower, "il était de type", user.TypeOfUser)
		if user.TypeOfUser == "YesVoter" {
			counterYesVoter++
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
	fmt.Println("There is ", counterYesVoter, "yesVoter,", counterThresholdVoter, "Threshold Voter,", counterNonResponsibleVoter, "NonresponsibleVoter,", counterResponsibleVoter, "ResponsibleVoter,", counterIndecisiveVoter, "IndecisiveVoter and", counterNormalVoter, "normalVoter")

	VoteInstance.ConstructTextForGraphCandidates(out, VoteInstance.GetResults())

}

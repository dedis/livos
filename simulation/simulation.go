package simulation

import (
	"fmt"
	"io"
	"strconv"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"github.com/mazen160/go-random"
	"golang.org/x/xerrors"
)

// GenerateItemsGraphviz creates a graphviz representation of the items. One can
// generate a graphical representation with `dot -Tpdf graph.dot -o graph.pdf`
func Simulation(out io.Writer) {

	var VoteList = make(map[string]voting.VotingInstance)
	var VoteSystem = impl.NewVotingSystem(nil, VoteList)
	var histoChoice = make([]voting.Choice, 0)

	var randomNumOfUser, err = random.IntRange(20, 22)
	if err != nil {
		xerrors.Errorf(err.Error())
	}

	//Random creating of a user and adds it to the list of voters
	var voters = make([]*voting.User, 0)
	for i := 0; i < randomNumOfUser; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
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
			//fmt.Println("USER :::::::", *user)
			if user.VotingPower > 0 {
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
					quantity_to_deleg, err := impl.NewLiquid(randomQuantityToDelegate)
					if err != nil {
						fmt.Println(err.Error(), "fail to do quantity to deleg")
					}
					err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
					if err != nil {
						fmt.Println(err.Error())
					}

					fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", voters[randomDelegateToIndex].UserID)

				} else if randomAction == 2 {
					//Vote action

					quantity := user.VotingPower
					quantity_to_Vote, err := impl.NewLiquid(quantity)
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

					fmt.Println(user.UserID, " a voté pour ", quantity, "%")

				}

			}
		}
	}

	for _, user := range VoteInstance.GetConfig().Voters {
		fmt.Println("Voting power of ", user.UserID, " = ", user.VotingPower)
	}

	VoteInstance.ConstructTextForGraph(out)
}

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
func Simulation(out io.Writer) {

	var VoteList = make(map[string]*impl.VotingInstance)
	var VoteSystem = impl.NewVotingSystem(nil, VoteList)
	var histoChoice = make([]voting.Choice, 0)

	var randomNumOfUser, err = random.IntRange(20, 25)
	if err != nil {
		xerrors.Errorf(err.Error())
	}

	//Random creating of a user and adds it to the list of voters
	var voters = make([]*voting.User, 0)
	for i := 0; i < randomNumOfUser; i++ {
		var user, err = VoteSystem.NewUser("user"+strconv.FormatInt(int64(i), 10), make(map[string]voting.Liquid), make(map[string]voting.Liquid), voting.Choice{}, histoChoice)
		if err != nil {
			xerrors.Errorf(err.Error())
		}
		voters = append(voters, &user)
	}

	//candidats
	var candidats = make([]string, 3)

	//empty list of votes
	var votes = make(map[string]voting.Choice)

	//creation of votingConfig
	voteConfig, err := impl.NewVotingConfig(voters, "Simulation 1", "Sunny day everyday ?", candidats)
	if err != nil {
		fmt.Println(err.Error())
	}

	//creation of the voting instance
	VoteInstance, err := VoteSystem.CreateAndAdd("Simulation01", voteConfig, "open", votes)
	if err != nil {
		fmt.Println(err.Error())
	}

	for ok := true; ok; ok = VoteInstance.CheckVotingPowerOfVoters() {
		for i, user := range VoteInstance.Config.Voters {
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
					quantity_to_deleg, err := impl.NewLiquid(float64(randomQuantityToDelegate))
					if err != nil {
						fmt.Println(err.Error(), "fail to do quantity to deleg")
					}
					err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
					if err != nil {
						fmt.Println(err.Error())
					}

					fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", voters[randomDelegateToIndex].UserID)

				} else if randomAction == 2 {
					//Vote YES action

					//quantity to yes vote
					randomQuantityToYesVote, err := random.IntRange(1, int(user.VotingPower)/10)
					if err != nil {
						fmt.Println(err.Error(), "fail to do randomQuantityToYesVote ")
					}
					randomQuantityToYesVote *= 10
					quantity_to_yesVote, err := impl.NewLiquid(float64(randomQuantityToYesVote))
					if err != nil {
						fmt.Println(err.Error())
					}
					liquid_0, err := impl.NewLiquid(0)
					if err != nil {
						fmt.Println(err.Error())
					}

					choiceTab := make(map[string]voting.Liquid)
					choiceTab["yes"] = quantity_to_yesVote
					choiceTab["no"] = liquid_0

					//create choice
					choice, err := impl.NewChoice(choiceTab)
					if err != nil {
						fmt.Println(err.Error())
					}

					//set the choice
					err = VoteInstance.SetChoice(user, choice)
					if err != nil {
						fmt.Println(err.Error())
					}

					//cast the vote
					err = VoteInstance.CastVote(user)
					if err != nil {
						fmt.Println(err.Error())
					}

					fmt.Println(user.UserID, " a voté YES pour ", quantity_to_yesVote, "%")

				} else if randomAction == 3 {
					//Vote NO action

					//quantity to yes vote
					randomQuantityToNoVote, err := random.IntRange(1, int(user.VotingPower)/10)
					if err != nil {
						fmt.Println(err.Error(), "fail to get randomQuantityToNoVote ")
					}
					randomQuantityToNoVote *= 10
					quantity_to_noVote, err := impl.NewLiquid(float64(randomQuantityToNoVote))
					if err != nil {
						fmt.Println(err.Error())
					}
					liquid_0, err := impl.NewLiquid(0)
					if err != nil {
						fmt.Println(err.Error())
					}

					choiceTab := make(map[string]voting.Liquid)
					choiceTab["no"] = quantity_to_noVote
					choiceTab["yes"] = liquid_0

					//create choice
					choice, err := impl.NewChoice(choiceTab)
					if err != nil {
						fmt.Println(err.Error())
					}

					//set the choice
					err = VoteInstance.SetChoice(user, choice)
					if err != nil {
						fmt.Println(err.Error())
					}

					//cast the vote
					err = VoteInstance.CastVote(user)
					if err != nil {
						fmt.Println(err.Error())
					}

					fmt.Println(user.UserID, " a voté NON pour ", quantity_to_noVote, "%")

				}
			}
		}
	}

	for _, user := range VoteInstance.Config.Voters {
		fmt.Println("Voting power of ", user.UserID, " = ", user.VotingPower)
	}

	results := VoteInstance.GetResults()

	fmt.Fprintf(out, "digraph network_activity {\n")
	fmt.Fprintf(out, "labelloc=\"t\";")
	fmt.Fprintf(out, "label = <Votation Diagram of %d nodes.    Results are Yes = %v, No = %v <font point-size='10'><br/>(generated %s)</font>>;", len(voters)+2, results["yes"], results["no"], time.Now().Format("2 Jan 06 - 15:04:05"))
	fmt.Fprintf(out, "graph [fontname = \"helvetica\"];")
	fmt.Fprintf(out, "node [fontname = \"helvetica\" area = 10 fillcolor=gold];")
	fmt.Fprintf(out, "edge [fontname = \"helvetica\"];\n")

	for _, user := range VoteInstance.Config.Voters {

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

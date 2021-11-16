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

	var randomNumOfUser, err = random.IntRange(3, 4)
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

	// userNoemien, err := VoteSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), voting.Choice{}, histoChoice)
	// if err != nil {
	// 	xerrors.Errorf(err.Error())
	// }

	// userBastien, err := VoteSystem.NewUser("Bastien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), voting.Choice{}, histoChoice)
	// if err != nil {
	// 	xerrors.Errorf(err.Error())
	// }

	// //creation of list of voters
	// var voters = []*voting.User{&userNoemien, &userBastien}
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

	fmt.Println("LIST OF VOTERS :::::::", VoteInstance.Config.Voters)

	for ok := true; ok; ok = VoteInstance.CheckVotingPowerOfVoters() {
		for i, user := range VoteInstance.Config.Voters {
			fmt.Println("USER :::::::", *user)
			if user.VotingPower > 0 {
				randomAction, err := random.IntRange(1, 4)
				if err != nil {
					fmt.Println(err.Error())
				}

				if randomAction == 1 {
					//Delegation action

					//random index creation (must NOT be == to index of current user)
					randomDelegateToIndex, err := random.IntRange(0, len(voters))
					if err != nil {
						fmt.Println(err.Error())
					}
					for ok := true; ok; ok = (randomDelegateToIndex == i) {
						randomDelegateToIndex, err = random.IntRange(0, len(voters))
						if err != nil {
							fmt.Println(err.Error())
						}
					}
					randomQuantityToDelegate, err := random.IntRange(1, int(user.VotingPower/10))
					if err != nil {
						fmt.Println(err.Error())
					}
					randomQuantityToDelegate *= 10
					quantity_to_deleg, err := impl.NewLiquid(float64(randomQuantityToDelegate))
					if err != nil {
						fmt.Println(err.Error())
					}
					err = VoteInstance.DelegTo(user, voters[randomDelegateToIndex], quantity_to_deleg)
					if err != nil {
						fmt.Println(err.Error())
					}

				} else if randomAction == 2 {
					//Vote YES action

					//quantity to yes vote
					randomQuantityToYesVote, err := random.IntRange(1, int(user.VotingPower/10))
					if err != nil {
						fmt.Println(err.Error())
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

				} else if randomAction == 3 {
					//Vote NO action

					//quantity to yes vote
					randomQuantityToNoVote, err := random.IntRange(1, int(user.VotingPower/10))
					if err != nil {
						fmt.Println(err.Error())
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
				}
			}
		}
	}

	/* //creation of the yesChoice map (containing vote 100% for yes)
	yesChoice := make(map[string]voting.Liquid)
	//creation of the mid map choice (50%, 50%)
	midChoice := make(map[string]voting.Liquid)

	//liquid 100%
	liquid_100, err := impl.NewLiquid(100)
	if err != nil {
		fmt.Println(err.Error())
	}
	//liquid 0%
	liquid_0, err := impl.NewLiquid(0)
	if err != nil {
		fmt.Println(err.Error())
	}
	//liquid 50%
	liquid_50, err := impl.NewLiquid(50)
	if err != nil {
		fmt.Println(err.Error())
	}

	//fill the yesChoice map with the liquids
	yesChoice["yes"] = liquid_100
	yesChoice["no"] = liquid_0
	midChoice["yes"] = liquid_50
	midChoice["no"] = liquid_50

	//create noemien's choice (100%)
	choiceNoemien, err := impl.NewChoice(yesChoice)
	if err != nil {
		fmt.Println(err.Error())
	}
	//set the yesChoice for Noemien
	err = VoteInstance.SetChoice(&userNoemien, choiceNoemien)
	if err != nil {
		fmt.Println(err.Error())
	}
	//cast and register the vote for noemien
	err = VoteInstance.CastVote(&userNoemien)
	if err != nil {
		fmt.Println(err.Error())
	}

	//delegation of the 100% of the voting power of bastien to noemien
	err = VoteInstance.DelegTo(&userBastien, &userNoemien, liquid_100)
	if err != nil {
		fmt.Println(err.Error())
	}

	//fmt.Println("ICI MAP(bastien, 100)", userNoemien.DelegatedFrom)

	//update noemien's choice (50%, 50%)
	choiceNoemien, err = impl.NewChoice(midChoice)
	if err != nil {
		fmt.Println(err.Error())
	}

	//set the midChoice for Noemien
	err = VoteInstance.SetChoice(&userNoemien, choiceNoemien)
	if err != nil {
		fmt.Println(err.Error())
	}
	//cast and register the second vote for noemien
	err = VoteInstance.CastVote(&userNoemien)
	if err != nil {
		fmt.Println(err.Error())
	} */

	fmt.Fprintf(out, "digraph network_activity {\n")
	fmt.Fprintf(out, "labelloc=\"t\";")
	fmt.Fprintf(out, "label = <Network Diagram of %d nodes <font point-size='10'><br/>(generated %s)</font>>;", len(voters)+2, time.Now().Format("2 Jan 06 - 15:04:05"))
	fmt.Fprintf(out, "graph [fontname = \"helvetica\"];")
	fmt.Fprintf(out, "node [fontname = \"helvetica\"];")
	fmt.Fprintf(out, "edge [fontname = \"helvetica\"];\n")

	for _, user := range VoteInstance.Config.Voters {

		color := "#4AB2FF"

		for _, choice := range user.HistoryOfChoice {
			if choice.VoteValue["yes"].Percentage != 0. {
				fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
					"[ label = < <font color='#303030'><b>%v</b></font><br/>> color=\"%s\" ];\n",
					user.UserID, "YES", choice.VoteValue["yes"].Percentage, color)
			}

			if choice.VoteValue["no"].Percentage != 0. {
				fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
					"[ label = < <font color='#303030'><b>%v</b></font><br/>> color=\"%s\" ];\n",
					user.UserID, "NO", choice.VoteValue["no"].Percentage, color)
			}
		}

		for other, quantity := range user.DelegatedTo {
			fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
				"[ label = < <font color='#303030'><b>%v</b></font><br/>> color=\"%s\" ];\n",
				user.UserID, other, quantity.Percentage, color)
		}
	}

	fmt.Fprintf(out, "}\n")
}

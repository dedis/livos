package impl

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting"
	"github.com/mazen160/go-random"
	"github.com/yourbasic/graph"
	"golang.org/x/xerrors"
)

//VOTING IMPLEMENTATION

const PERCENTAGE = 100

const InitialVotingPower = 100.

type VotingInstance struct {
	//voting instance's id
	Id string

	//parameters of the Voting
	Config voting.VotingConfig

	// open / closed
	Status string

	// Votes contains the choice of each voter, references by a userID
	//Votes map[string]voting.Choice

	//db.socket personnalisé pour chacun ?
}

//VOTING INSTANCE FUNCTIONS :::::

//Close the voting instance session (open -> close)
func (vi *VotingInstance) CloseVoting() {
	vi.SetStatus("close")
}

//Set the status of the voting instance (open or close)
func (vi *VotingInstance) SetStatus(status string) error {
	if !(status == "open") && !(status == "close") {
		return xerrors.Errorf("The status is incorrect, should be either 'open' or 'close'.")
	}
	vi.Status = status
	return nil
}

//Set the title of the voting instance
func (vi *VotingInstance) SetTitle(title string) error {
	vi.Config.Title = title
	return nil
}

//Set the voters of the voting instance
func (vi *VotingInstance) SetVoters(users []*voting.User) error {
	vi.Config.Voters = users
	return nil
}

//Set the candidates of the voting instance
func (vi *VotingInstance) SetCandidates(candidates []*voting.Candidate) error {
	vi.Config.Candidates = candidates
	return nil
}

//Set the description of the voting instance
func (vi *VotingInstance) SetDescription(description string) error {
	vi.Config.Description = description
	return nil
}

//Set the type of Voting config candidate or yes/no question
func (vi *VotingInstance) SetTypeOfVotingConfig(typeOfVotingConfig string) error {
	vi.Config.TypeOfVotingConfig = voting.TypeOfVotingConfig(typeOfVotingConfig)
	return nil
}

//Return the status of the voting instance
func (vi *VotingInstance) GetStatus() string {
	return vi.Status
}

//Return the configuration of the voting instance
func (vi *VotingInstance) GetConfig() voting.VotingConfig {
	return vi.Config
}

func (vi *VotingInstance) RandomWithProbabilities(user *voting.User) int {
	randomAction, err := random.IntRange(1, 3)
	if err != nil {
		fmt.Println(err.Error())
	}
	nbOfCand := 0
	if randomAction == 1 {
		nbOfCand = 1
	} else {
		randomAction2, err := random.IntRange(1, 11)
		if err != nil {
			fmt.Println(err.Error())
		}
		if randomAction2 < 7 {
			nbOfCand = 2
		} else {
			nbOfCand = 3
		}
	}
	return nbOfCand
}

//Give the result of the choices of the voting instance in the form: map[no:50 yes:50]
func (vi *VotingInstance) GetResults() map[string]float64 {
	if vi.Config.TypeOfVotingConfig == "YesOrNoQuestion" {
		results := make(map[string]float64, len(vi.Config.Voters))
		var yesPower int = 0
		var noPower int = 0
		for _, user := range vi.Config.Voters {
			for _, choice := range user.HistoryOfChoice {
				yesPower += choice.VoteValue["yes"].Percentage
				noPower += choice.VoteValue["no"].Percentage
			}
		}
		//in order to get 4 and not 4.6666666... for example
		// var temp1 = float64(int(yesPower/float64(counter))*100) / 100
		// var temp2 = float64(int(noPower/float64(counter))*100) / 100
		results["yes"] = float64(yesPower / (len(vi.Config.Voters)))
		results["no"] = float64(noPower / (len(vi.Config.Voters)))

		return results
	} else {
		resultsCandidate := make(map[string]float64, len(vi.Config.Candidates))

		for _, candidate := range vi.Config.Candidates {
			var candidateResult int = 0
			for _, user := range vi.Config.Voters {
				for _, choice := range user.HistoryOfChoice {
					candidateResult += choice.VoteValue[candidate.CandidateID].Percentage
				}
			}
			resultsCandidate[candidate.CandidateID] = float64(candidateResult / len(vi.Config.Voters))
		}
		return resultsCandidate
	}
}

//Give the result of the choices of the voting instance in the form: map[no:50 yes:50]
func (vi *VotingInstance) GetResultsQuadraticVoting() map[string]float64 {
	if vi.Config.TypeOfVotingConfig == "YesOrNoQuestion" {
		results := make(map[string]float64, len(vi.Config.Voters))
		var yesPower int = 0
		var noPower int = 0
		for _, user := range vi.Config.Voters {
			for _, choice := range user.HistoryOfChoice {
				yesPower += choice.VoteValue["yes"].Percentage
				noPower += choice.VoteValue["no"].Percentage
			}
		}
		//in order to get 4 and not 4.6666666... for example
		// var temp1 = float64(int(yesPower/float64(counter))*100) / 100
		// var temp2 = float64(int(noPower/float64(counter))*100) / 100
		results["yes"] = float64(yesPower / len(vi.Config.Voters))
		results["no"] = float64(noPower / len(vi.Config.Voters))

		return results
	} else {
		resultsCandidate := make(map[string]float64, len(vi.Config.Candidates))

		for _, candidate := range vi.Config.Candidates {
			var candidateResult float64 = 0
			for _, user := range vi.Config.Voters {
				for _, choice := range user.HistoryOfChoice {
					intermediateUserResult := float64(choice.VoteValue[candidate.CandidateID].Percentage)
					rootedUserResult := math.Sqrt(intermediateUserResult/20) * 20
					candidateResult += rootedUserResult
				}
			}
			resultsCandidate[candidate.CandidateID] = candidateResult / float64(len(vi.Config.Voters))
		}
		var totalSumOfAllCandidates float64 = 0
		for _, res := range resultsCandidate {
			totalSumOfAllCandidates += res
		}

		for cand, res := range resultsCandidate {
			resultsCandidate[cand] = (res / totalSumOfAllCandidates) * 100
		}

		return resultsCandidate
	}
}

//Return the user object from the string passed in argument
func (vi *VotingInstance) GetUser(userID string) (*voting.User, error) {
	for _, value := range vi.Config.Voters {
		if value.UserID == userID {
			return value, nil
		}
	}
	return nil, xerrors.Errorf("Cannot find the user. UserId is incorrect.")
}

//Return the candidate object from the string passed in argument
func (vi *VotingInstance) GetCandidate(candidateID string) (*voting.Candidate, error) {
	for _, value := range vi.Config.Candidates {
		if value.CandidateID == candidateID {
			return value, nil
		}
	}
	return nil, xerrors.Errorf("Cannot find the Candidate. CandidateID is incorrect.")
}

//Check that the voting power of a user is always positive
func (vi *VotingInstance) CheckVotingPower(user *voting.User) error {
	if user.VotingPower < 0 {
		return xerrors.Errorf("Value of voting power is negative: Was %d, must be more or equal than 0", int(user.VotingPower))
	}
	return nil
}

//Check if all the voters have used all their voting power (usefull for simulation)
func (vi *VotingInstance) CheckVotingPowerOfVoters() bool {
	for _, user := range vi.Config.Voters {
		if user.VotingPower != 0. {
			return true
		}
	}
	return false
}

//Set (link) a given choice to a given user
func (vi *VotingInstance) SetVote(user *voting.User, choice voting.Choice) error {

	var sumOfVotingPower int = 0
	for _, v := range choice.VoteValue {
		sumOfVotingPower += v.Percentage
	}

	if user.VotingPower-sumOfVotingPower >= 0 {
		user.VotingPower -= sumOfVotingPower

		//update history of choice with the current choice
		histVoteValue := make(map[string]voting.Liquid)
		for key, value := range choice.VoteValue {
			histVoteValue[key] = value
		}
		histChoice, err := NewChoice(histVoteValue)
		if err != nil {
			return xerrors.Errorf(err.Error())
		}
		user.HistoryOfChoice = append(user.HistoryOfChoice, histChoice)
	} else {
		return xerrors.Errorf("Voting power can't be negative.")
	}

	return nil
}

//DIFFERENT FUNCTIONS FOR THE VOTE SIMULATIONS

//YES OR NO VOTE FUNCTIONS ----------------------------------------------------------------------------------------

func (vi *VotingInstance) YesVote(user *voting.User, votingPower int) {
	quantity := votingPower
	quantity_to_Vote, err := NewLiquid(quantity)
	if err != nil {
		fmt.Println(err.Error())
	}
	liquid_0, err := NewLiquid(0)
	if err != nil {
		fmt.Println(err.Error())
	}

	choiceTab := make(map[string]voting.Liquid)

	choiceTab["yes"] = quantity_to_Vote
	choiceTab["no"] = liquid_0
	//create choice
	choice, err := NewChoice(choiceTab)
	if err != nil {
		fmt.Println(err.Error())
	}

	//set the choice
	err = vi.SetVote(user, choice)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(user.UserID, " a voté pour ", quantity, "%", "il était", user.TypeOfUser)
}

func (vi *VotingInstance) NoVote(user *voting.User, votingPower int) {
	quantity := votingPower
	quantity_to_Vote, err := NewLiquid(quantity)
	if err != nil {
		fmt.Println(err.Error())
	}
	liquid_0, err := NewLiquid(0)
	if err != nil {
		fmt.Println(err.Error())
	}

	choiceTab := make(map[string]voting.Liquid)

	choiceTab["no"] = quantity_to_Vote
	choiceTab["yes"] = liquid_0
	//create choice
	choice, err := NewChoice(choiceTab)
	if err != nil {
		fmt.Println(err.Error())
	}

	//set the choice
	err = vi.SetVote(user, choice)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(user.UserID, " a voté pour ", quantity, "%", "il était", user.TypeOfUser)
}

func (vi *VotingInstance) IndecisiveVote(user *voting.User, i int, quantityToDeleg int) {
	//Delegation action
	//random index creation (must NOT be == to index of current user)

	randomDelegateToIndex, err := random.IntRange(0, len(vi.GetConfig().Voters))
	if err != nil {
		fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
	}
	if len(user.DelegatedTo) != 0 {
		for nameDelegat := range user.DelegatedTo {
			for index, nameVoter := range vi.Config.Voters {
				if nameDelegat == nameVoter.UserID {
					randomDelegateToIndex = index
				}
			}
		}
	} else {
		for ok := true; ok; ok = (randomDelegateToIndex == i) {
			randomDelegateToIndex, err = random.IntRange(0, len(vi.GetConfig().Voters))
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
			}
		}
	}

	quantity_to_deleg, err := NewLiquid(quantityToDeleg)
	if err != nil {
		fmt.Println(err.Error(), "fail to do quantity to deleg")
	}
	err = vi.DelegTo(user, vi.GetConfig().Voters[randomDelegateToIndex], quantity_to_deleg)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", vi.GetConfig().Voters[randomDelegateToIndex].UserID, "il était", user.TypeOfUser)
}

func (vi *VotingInstance) BreakTheCycle(user *voting.User) {
	RandomAction, err := random.IntRange(0, 100)
	if err != nil {
		fmt.Println(err.Error(), "fail to do RandomAction first time")
	}
	switch {
	case RandomAction < 52:
		RandomActionToChooseWhatToVote, err := random.IntRange(0, 3)
		if err != nil {
			fmt.Println(err.Error(), "fail to do RandomActionToChooseWhatToVote first time")
		}
		if RandomActionToChooseWhatToVote == 0 {
			vi.YesVote(user, user.VotingPower)
		} else {
			vi.NoVote(user, user.VotingPower)
		}
	case RandomAction < 84:
		//delegate to another person
	case RandomAction < 98:
		//vote blanc
	default:
		//do nothing as don't vote
	}
}

func (vi *VotingInstance) RandomVote(user *voting.User, i int) {
	randomAction, err := random.IntRange(1, 3)
	if err != nil {
		fmt.Println(err.Error(), "fail to do randomAction")
	}

	if randomAction == 1 {
		//Delegation action

		//random index creation (must NOT be == to index of current user)
		randomDelegateToIndex, err := random.IntRange(0, len(vi.GetConfig().Voters))
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
		}
		for ok := true; ok; ok = (randomDelegateToIndex == i) {
			randomDelegateToIndex, err = random.IntRange(0, len(vi.GetConfig().Voters))
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
			}
		}
		randomQuantityToDelegate, err := random.IntRange(1, int(user.VotingPower/10)+1)
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomQuantityToDelegate")
		}
		randomQuantityToDelegate *= 10
		quantity_to_deleg, err := NewLiquid(randomQuantityToDelegate)
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to deleg")
		}
		err = vi.DelegTo(user, vi.GetConfig().Voters[randomDelegateToIndex], quantity_to_deleg)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", vi.GetConfig().Voters[randomDelegateToIndex].UserID, "il était", user.TypeOfUser)

	} else if randomAction == 2 {
		//Vote action

		quantity := user.VotingPower
		if len(user.HistoryOfChoice) == 0 {
			yesOrNo, err := random.IntRange(1, 3)
			if err != nil {
				fmt.Println(err.Error(), "fail to do yesOrNo ")
			}

			if yesOrNo == 1 {
				vi.YesVote(user, quantity)
			} else {
				vi.NoVote(user, quantity)
			}
		} else if user.HistoryOfChoice[0].VoteValue["no"].Percentage != 0. {
			vi.NoVote(user, quantity)
		} else {
			vi.YesVote(user, quantity)
		}

	}
}

func (vi *VotingInstance) ThresholdVote(user *voting.User, i int, threshold int) {

	var thresholdComparator = 0
	for i := range user.HistoryOfChoice {
		thresholdComparator += user.HistoryOfChoice[i].VoteValue["yes"].Percentage
		thresholdComparator += user.HistoryOfChoice[i].VoteValue["no"].Percentage
	}

	if thresholdComparator > threshold {
		//Delegation action
		vi.IndecisiveVote(user, i, user.VotingPower)

	} else {
		//Vote action
		quantity, err := random.IntRange(1, (int(user.VotingPower)/10)+1)
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to voteAction")
		}
		quantity_to_deleg, err := NewLiquid(int(user.VotingPower - quantity*10))
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to deleg")
		}

		if len(user.HistoryOfChoice) == 0 {
			yesOrNo, err := random.IntRange(1, 3)
			if err != nil {
				fmt.Println(err.Error(), "fail to do yesOrNo ")
			}

			if yesOrNo == 1 {
				vi.YesVote(user, quantity*10)
			} else {
				vi.NoVote(user, quantity*10)
			}
			vi.IndecisiveVote(user, i, quantity_to_deleg.Percentage)
		} else if user.HistoryOfChoice[0].VoteValue["no"].Percentage != 0. {
			vi.NoVote(user, quantity*10)
			vi.IndecisiveVote(user, i, quantity_to_deleg.Percentage)
		} else {
			vi.YesVote(user, quantity*10)
			vi.IndecisiveVote(user, i, quantity_to_deleg.Percentage)
		}
	}
}

func (vi *VotingInstance) NonResponsibleVote(user *voting.User, i int) {
	if len(user.HistoryOfChoice) == 0 {
		var randomNumberToChooseYesOrNo, err = random.IntRange(0, 2)
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
		}
		if randomNumberToChooseYesOrNo == 0 {
			vi.YesVote(user, InitialVotingPower)
		} else {
			vi.NoVote(user, InitialVotingPower)
		}
	} else {
		//Delegation action
		vi.IndecisiveVote(user, i, user.VotingPower)
	}
}

func (vi *VotingInstance) ResponsibleVote(user *voting.User, i int) {
	randomAction, err := random.IntRange(1, 4)
	if err != nil {
		fmt.Println(err.Error(), "fail to do randomAction")
	}

	if len(user.HistoryOfChoice) != 0 {
		randomAction = 2
	} else if len(user.DelegatedTo) != 0 {
		randomAction = 1
	}

	if randomAction == 1 {
		//Delegation action
		vi.IndecisiveVote(user, i, user.VotingPower)

	} else if randomAction == 2 || randomAction == 3 {
		//Vote action
		quantity, err := random.IntRange(1, (int(user.VotingPower)/10)+1)
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to voteAction")
		}
		quantity_to_deleg, err := NewLiquid(int(user.VotingPower - quantity*10))
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to deleg")
		}

		if len(user.HistoryOfChoice) == 0 {
			yesOrNo, err := random.IntRange(1, 3)
			if err != nil {
				fmt.Println(err.Error(), "fail to do yesOrNo ")
			}

			if yesOrNo == 1 {
				vi.YesVote(user, quantity*10)
			} else {
				vi.NoVote(user, quantity*10)
			}
			vi.IndecisiveVote(user, i, quantity_to_deleg.Percentage)
		} else if user.HistoryOfChoice[0].VoteValue["no"].Percentage != 0. {
			vi.NoVote(user, quantity*10)
			vi.IndecisiveVote(user, i, quantity_to_deleg.Percentage)
		} else {
			vi.YesVote(user, quantity*10)
			vi.IndecisiveVote(user, i, quantity_to_deleg.Percentage)
		}
	}
}

//CANDIDATS VOTE FUNCTIONS ----------------------------------------------------------------------------------------

func (vi *VotingInstance) CandidateVote(user *voting.User, i int, votingPower int) {

	choiceTab := make(map[string]voting.Liquid)
	CandidatMappedPercentage := make(map[string]float64)

	if len(user.HistoryOfChoice) != 0 {
		//the user already voted for some candidates => must distribute the new voting in the same manner
		total_votingPower_spent := 0
		for _, choice := range user.HistoryOfChoice {
			for _, value := range choice.VoteValue {
				total_votingPower_spent += value.Percentage
			}
		}

		//create a map CANDIDATE_NAME => PERCENTAGE VOTED FOR
		for _, choice := range user.HistoryOfChoice {
			for name, value := range choice.VoteValue {
				CandidatMappedPercentage[name] = CandidatMappedPercentage[name] + (float64(value.Percentage)/float64(total_votingPower_spent))*100
			}
		}

		//the case where user have less that 10 voting power, we just redirect it onto his favorite precedent choice
		if votingPower <= 10 {
			var prefered_candidate string
			var max = 0.
			for name, val := range CandidatMappedPercentage {
				if val > max {
					max = val
					prefered_candidate = name
				}
			}

			var err error
			choiceTab[prefered_candidate], err = NewLiquid(votingPower)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(user.UserID, " a voté pour ", votingPower, "à", prefered_candidate, "car il est ", user.TypeOfUser, "et qu'il as moins de 10 de voting power")

		} else {
			for candID, percentage := range CandidatMappedPercentage {
				//re-vote the right amount of the new voting power given the percentage
				var err error
				choiceTab[candID], err = NewLiquid(int(votingPower * int(percentage) / 100))
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(user.UserID, " a voté pour ", votingPower*int(percentage)/100, "à", candID, "car il est ", user.TypeOfUser, "(revote)")
			}
		}

		//create choice
		choice, err := NewChoice(choiceTab)
		if err != nil {
			fmt.Println(err.Error())
		}

		//set the choice
		err = vi.SetVote(user, choice)
		if err != nil {
			fmt.Println(err.Error())
		}

	} else {
		//the user hasn't voted yet =>
		//while the user still have some voting power, vote amongst different candidates
		quantity1, quantity2, quantity3 := vi.SplitVotingPower(user, i, votingPower)

		vi.RandomVoteCandidate(user, i, quantity1*10)

		if quantity2 != 0 {
			vi.RandomVoteCandidate(user, i, quantity2*10)
		}

		if quantity3 != 0 {
			vi.RandomVoteCandidate(user, i, quantity3*10)
		}
	}
}

func (vi *VotingInstance) BreakTheCycleCandidate(user *voting.User, i int, votingPower int) {
	RandomAction, err := random.IntRange(0, 101)
	if err != nil {
		fmt.Println(err.Error(), "fail to do RandomAction first time")
	}
	switch {
	case RandomAction < 62:
		//forced to vote
		vi.CandidateVote(user, i, user.VotingPower)
	default:
		//delegate to ANOTHER person

		//the liquid struct of the voting power to deleg
		quantity_to_deleg, err := NewLiquid(votingPower)
		if err != nil {
			fmt.Println(err.Error())
		}

		indexBelongToDelegatedToList := func(idx int) bool {
			for name := range user.DelegatedTo {
				if vi.GetConfig().Voters[idx].UserID == name {
					return true
				}
			}
			return false
		}

		//random index creation (must NOT be == to index of current user or to index of already delegated people)
		randomDelegateToIndex, err := random.IntRange(0, len(vi.GetConfig().Voters))
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
		}
		for ok := true; ok; ok = (randomDelegateToIndex == i || indexBelongToDelegatedToList(randomDelegateToIndex)) {
			randomDelegateToIndex, err = random.IntRange(0, len(vi.GetConfig().Voters))
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
			}
		}

		err = vi.DelegTo(user, vi.GetConfig().Voters[randomDelegateToIndex], quantity_to_deleg)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(user.UserID, " a délégué pour ", votingPower, "à", vi.GetConfig().Voters[randomDelegateToIndex].UserID, "car il est ", user.TypeOfUser, "(break the cycle)")
	}
}

func (vi *VotingInstance) IndecisiveVoteCandidate(user *voting.User, i int, votingPower int) {

	UsersMappedPercentage := make(map[string]float64)

	if len(user.HistoryOfChoice) != 0 || len(user.DelegatedTo) != 0 {
		//the user already voted => must distribute the new voting in the same manner

		//construction du graph de delegation
		findIndexInListOfUser := func(list []*voting.User, user *voting.User) int {
			for i, u := range list {
				if u == user {
					return i
				}
			}
			return -1
		}

		//construct the graph of delegation, run acyclic detection algorithm to detect cycles
		g := graph.New(len(vi.GetConfig().Voters))

		contructPartialGraphWithDelegatedToList := func(curr_user *voting.User, i int, g *graph.Mutable) {
			for other := range curr_user.DelegatedTo {
				otherObject, err := vi.GetUser(other)
				if err != nil {
					fmt.Println(err.Error())
				}
				j := findIndexInListOfUser(vi.GetConfig().Voters, otherObject)
				g.Add(i, j)
			}
		}

		for idx, u := range vi.GetConfig().Voters {
			contructPartialGraphWithDelegatedToList(u, idx, g)
		}

		IsThereCycle := !graph.Acyclic(g)

		if IsThereCycle {
			vi.BreakTheCycleCandidate(user, i, votingPower)
		} else {
			//on a oublié un cas. Le cas où un Non responsible voter vote ses 100 à Trump, puis on lui redonne 100.
			//alors il vient dans cette boucle (assuming there is no circle) mais il se construit des truc qui valent 0 car il n'as jamais délégué!
			//if il n'as pas délégué encore
			if len(user.DelegatedTo) == 0 {
				vi.RandomDelegate(user, i, votingPower)
			} else {
				//redelegate in the same manner as before
				total_votingPower_delegated := 0
				for _, value := range user.DelegatedTo {
					total_votingPower_delegated += value.Percentage
				}

				//create a map USERS_NAME => PERCENTAGE VOTED FOR
				for name, value := range user.DelegatedTo {
					UsersMappedPercentage[name] = UsersMappedPercentage[name] + (float64(value.Percentage)/float64(total_votingPower_delegated))*100
				}

				//the case where user have less that 10 voting power, we just redirect it onto his favorite precedent choice
				if votingPower <= 10 {
					var prefered_user string
					var max = 0.
					for name, val := range UsersMappedPercentage {
						if val > max {
							max = val
							prefered_user = name
						}
					}

					prefered_user_object, err := vi.GetUser(prefered_user)
					if err != nil {
						fmt.Println(err.Error())
					}
					quantity_to_deleg, err := NewLiquid(votingPower)
					if err != nil {
						fmt.Println(err.Error())
					}
					err = vi.DelegTo(user, prefered_user_object, quantity_to_deleg)
					if err != nil {
						fmt.Println(err.Error())
					}
					fmt.Println(user.UserID, " a délégué pour ", quantity_to_deleg.Percentage, "à", prefered_user_object.UserID, "car il est ", user.TypeOfUser, "et qu'il as moins de 10 de voting power (redelegate)")

				} else {
					//re-delegate the right amount of the new voting power given the percentage of previous delegation
					for other, percentage := range UsersMappedPercentage {
						quantity_to_deleg, err := NewLiquid(int(votingPower * int(percentage) / 100))
						if err != nil {
							fmt.Println(err.Error(), "fail to do quantity to deleg")
						}

						otherObject, err := vi.GetUser(other)

						err = vi.DelegTo(user, otherObject, quantity_to_deleg)
						if err != nil {
							fmt.Println(err.Error())
						}
						fmt.Println(user.UserID, " a délégué pour ", quantity_to_deleg.Percentage, "à", otherObject.UserID, "car il est ", user.TypeOfUser, "(redelegate)")
					}
				}
			}

		}

	} else {
		quantity1, quantity2, quantity3 := vi.SplitVotingPower(user, i, votingPower)

		vi.RandomDelegate(user, i, quantity1*10)

		if quantity2 != 0 {
			vi.RandomDelegate(user, i, quantity2*10)
		}

		if quantity3 != 0 {
			vi.RandomDelegate(user, i, quantity3*10)
		}
	}
}

func (vi *VotingInstance) DefaultVoteCandidate(user *voting.User, i int) {
	randomAction, err := random.IntRange(1, 3)
	if err != nil {
		fmt.Println(err.Error(), "fail to do randomAction")
	}

	if randomAction == 1 {
		//Delegation action

		//random index creation (must NOT be == to index of current user)
		randomDelegateToIndex, err := random.IntRange(0, len(vi.GetConfig().Voters))
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
		}
		for ok := true; ok; ok = (randomDelegateToIndex == i) {
			randomDelegateToIndex, err = random.IntRange(0, len(vi.GetConfig().Voters))
			if err != nil {
				fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
			}
		}
		randomQuantityToDelegate, err := random.IntRange(1, int(user.VotingPower/10)+1)
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomQuantityToDelegate")
		}
		randomQuantityToDelegate *= 10
		quantity_to_deleg, err := NewLiquid(randomQuantityToDelegate)
		if err != nil {
			fmt.Println(err.Error(), "fail to do quantity to deleg")
		}
		err = vi.DelegTo(user, vi.GetConfig().Voters[randomDelegateToIndex], quantity_to_deleg)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(user.UserID, " a delegué ", quantity_to_deleg, " à : ", vi.GetConfig().Voters[randomDelegateToIndex].UserID, "il était", user.TypeOfUser)

	} else if randomAction == 2 {
		//Vote action

		quantity := user.VotingPower
		vi.CandidateVote(user, i, quantity)

	}
}

func (vi *VotingInstance) ThresholdVoteCandidate(user *voting.User, i int, threshold int, votingPower int) {
	total_votingPower_spent := 0
	for _, choice := range user.HistoryOfChoice {
		for _, value := range choice.VoteValue {
			total_votingPower_spent += value.Percentage
		}
	}

	if total_votingPower_spent > threshold {
		//Delegation action
		vi.IndecisiveVoteCandidate(user, i, votingPower)
	} else {
		//Vote action
		vi.ResponsibleVoteCandidate(user, i, votingPower)
	}
}

func (vi *VotingInstance) NonResponsibleVoteCandidate(user *voting.User, i int, votingPower int) {
	if len(user.HistoryOfChoice) != 0 || len(user.DelegatedTo) != 0 {
		//the user already voted => must delegate all (SHY VOTER)
		vi.IndecisiveVoteCandidate(user, i, votingPower)
	} else {
		//the user hasn't voted yet =>
		//while the user still have some voting power, vote amongst different candidates and delegate
		//split between proportion to delegate and proportion to vote (can be 30/70 or 50/50 etc...)
		vi.SplitVPintoActions(user, i, 100)
	}
}

func (vi *VotingInstance) ResponsibleVoteCandidate(user *voting.User, i int, votingPower int) {

	choiceTab := make(map[string]voting.Liquid)
	CandidatMappedPercentage := make(map[string]float64)

	if len(user.HistoryOfChoice) != 0 || len(user.DelegatedTo) != 0 {

		if votingPower <= 10 {

			CandidatMappedValue := make(map[string]int)

			var prefered string
			var max = 0
			for name, val := range user.DelegatedTo {
				if val.Percentage > max {
					max = val.Percentage
					prefered = name
				}
			}

			//create a map CANDIDATE_NAME =>  CUMUL VALUE VOTED FOR
			for _, choice := range user.HistoryOfChoice {
				for name, value := range choice.VoteValue {
					CandidatMappedValue[name] = CandidatMappedValue[name] + (value.Percentage)
				}
			} //trump => 120, macron => 40

			for name, val := range CandidatMappedValue {
				if val > max {
					max = val
					prefered = name
				}
			}

			var err error
			//at this state 'prefered' contains the name of the candidate or voter that the user gave the most votingpower

			prefered_obj, err := vi.GetUser(prefered)
			if prefered_obj == nil {
				//case where 'prefered' is a candidate
				choiceTab[prefered], err = NewLiquid(votingPower)
				if err != nil {
					fmt.Println(err.Error())
				}
				//create choice
				choice, err := NewChoice(choiceTab)
				if err != nil {
					fmt.Println(err.Error())
				}

				//set the choice
				err = vi.SetVote(user, choice)
				if err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println(user.UserID, " a voté pour ", votingPower, "à", prefered, "car il est ", user.TypeOfUser, "moins de 10 voting power, case where prefered is a candidate (revote)")
			} else {
				//case where 'prefered' is a voter
				quantity_to_deleg, err := NewLiquid(votingPower)
				if err != nil {
					fmt.Println(err.Error())
				}
				err = vi.DelegTo(user, prefered_obj, quantity_to_deleg)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(user.UserID, " a délégué pour ", quantity_to_deleg.Percentage, "à", prefered_obj.UserID, "car il est ", user.TypeOfUser, "moins de 10 voting power, case where prefered is a voter (redelege)")
			}
		} else {
			//the user already voted for some candidates => must distribute the new voting in the same manner
			total_votingPower_spent := 0
			for _, choice := range user.HistoryOfChoice {
				for _, value := range choice.VoteValue {
					total_votingPower_spent += value.Percentage
				}
			}
			for _, value := range user.DelegatedTo {
				total_votingPower_spent += value.Percentage
			}

			//create a map CANDIDATE_NAME => PERCENTAGE VOTED FOR
			for _, choice := range user.HistoryOfChoice {
				for name, value := range choice.VoteValue {
					CandidatMappedPercentage[name] = CandidatMappedPercentage[name] + (float64(value.Percentage)/float64(total_votingPower_spent))*100
				}
			} //trump => 45, Marcon => 35, user0 => 20

			var total_Percentage_Voted int = 0

			//re-vote the right amount of the new voting power given the percentage
			for candID, percentage := range CandidatMappedPercentage {
				var err error
				choiceTab[candID], err = NewLiquid(int(votingPower * int(percentage) / 100))
				if err != nil {
					fmt.Println(err.Error())
				}
				total_Percentage_Voted += int(percentage)
				fmt.Println(user.UserID, " a voté pour ", choiceTab[candID].Percentage, "à", candID, "car il est ", user.TypeOfUser, "(revote)")
			}

			//create choice
			choice, err := NewChoice(choiceTab)
			if err != nil {
				fmt.Println(err.Error())
			}

			//set the choice
			err = vi.SetVote(user, choice)
			if err != nil {
				fmt.Println(err.Error())
			}

			//how many percentage to delegate
			var total_Percentage_Delegated int = 100 - total_Percentage_Voted

			if total_Percentage_Delegated != 0 {
				vi.IndecisiveVoteCandidate(user, i, total_Percentage_Delegated/100*votingPower)
			}
		}

	} else {
		//the user hasn't voted yet =>
		//while the user still have some voting power, vote amongst different candidates and delegate :
		//split between proportion to delegate and proportion to vote (can be 30/70 or 50/50 etc...)
		vi.SplitVPintoActions(user, i, user.VotingPower)
	}
}

func (vi *VotingInstance) SplitVotingPower(user *voting.User, i int, votingPowerToSplit int) (int, int, int) {
	quantity1 := 0
	quantity2 := 0
	quantity3 := 0
	var err error
	nbOfAction := vi.RandomWithProbabilities(user)

	if nbOfAction == 1 {
		//You must vote for a candidate and therefore you must vote for him for 100
		quantity1 = votingPowerToSplit / 10
	} else if nbOfAction == 2 {
		//You must for a candidate
		quantity1, err = random.IntRange(1, (votingPowerToSplit / 10))
		if err != nil {
			fmt.Println(err.Error())
		}

		//Later we will choose what to do with this quantity either deleg or vote for another candidate
		quantity2 = votingPowerToSplit/10 - quantity1
	} else {
		//You must for a candidate
		quantity1, err = random.IntRange(1, votingPowerToSplit/10-1)
		if err != nil {
			fmt.Println(err.Error())
		}

		//Later we will choose what to do with this quantity either deleg or vote for another candidate
		quantity2, err = random.IntRange(1, votingPowerToSplit/10-quantity1)
		if err != nil {
			fmt.Println(err.Error())
		}

		//Later we will choose what to do with this quantity either deleg or vote for another candidate
		quantity3 = votingPowerToSplit/10 - quantity1 - quantity2
	}

	return quantity1, quantity2, quantity3
}

func (vi *VotingInstance) SplitVPintoActions(user *voting.User, i int, votingPowerToSplit int) {

	//at this state we have split in 1, 2 or 3 way the voting power inital (100)
	quantity1, quantity2, quantity3 := vi.SplitVotingPower(user, i, votingPowerToSplit)

	vi.RandomVoteCandidate(user, i, quantity1*10)

	if quantity2 != 0 {
		chooseBetweenVoteAndDeleg, err := random.IntRange(1, 11)
		if err != nil {
			fmt.Println(err.Error())
		}
		//you have 60% chance to vote as your second choice for a candidate
		if chooseBetweenVoteAndDeleg < 7 {
			vi.RandomVoteCandidate(user, i, quantity2*10)
		} else { //you delegate your vote
			vi.RandomDelegate(user, i, quantity2*10)
		}
	}

	if quantity3 != 0 { // same as quantity2
		chooseBetweenVoteAndDeleg, err := random.IntRange(1, 11)
		if err != nil {
			fmt.Println(err.Error())
		}
		if chooseBetweenVoteAndDeleg < 7 {
			vi.RandomVoteCandidate(user, i, quantity3*10)
		} else {
			vi.RandomDelegate(user, i, quantity3*10)
		}
	}
}

func (vi *VotingInstance) RandomDelegate(user *voting.User, i int, votingPower int) {
	//Delegation action
	//random index creation (must NOT be == to index of current user)
	randomDelegateToIndex, err := random.IntRange(0, len(vi.GetConfig().Voters))
	if err != nil {
		fmt.Println(err.Error(), "fail to do randomDelegateToIndex first time")
	}

	for ok := true; ok; ok = (randomDelegateToIndex == i) {
		randomDelegateToIndex, err = random.IntRange(0, len(vi.GetConfig().Voters))
		if err != nil {
			fmt.Println(err.Error(), "fail to do randomDelegateToIndex")
		}
	}

	quantity_to_deleg, err := NewLiquid(votingPower)
	if err != nil {
		fmt.Println(err.Error(), "fail to do quantity to deleg")
	}
	err = vi.DelegTo(user, vi.GetConfig().Voters[randomDelegateToIndex], quantity_to_deleg)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(user.UserID, " a delegué pour ", quantity_to_deleg.Percentage, " à : ", vi.GetConfig().Voters[randomDelegateToIndex].UserID, "il était", user.TypeOfUser, "(deleg to a random)")
}

func (vi *VotingInstance) RandomVoteCandidate(user *voting.User, i int, votingPower int) {

	choiceTab := make(map[string]voting.Liquid)

	//index of which candidate does the voter vote for
	candidateChoice, err := random.IntRange(0, len(vi.GetConfig().Candidates))
	if err != nil {
		fmt.Println(err.Error())
	}

	quantity_to_Vote, err := NewLiquid(votingPower)
	if err != nil {
		fmt.Println(err.Error())
	}

	choiceTab[vi.GetConfig().Candidates[candidateChoice].CandidateID], err = AddLiquid(choiceTab[vi.GetConfig().Candidates[candidateChoice].CandidateID], quantity_to_Vote)
	if err != nil {
		fmt.Println(err.Error())
	}

	//create choice
	choice, err := NewChoice(choiceTab)
	if err != nil {
		fmt.Println(err.Error())
	}

	//set the choice
	err = vi.SetVote(user, choice)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(user.UserID, " a voté ", votingPower, " à : ", vi.GetConfig().Candidates[candidateChoice].CandidateID, "il était", user.TypeOfUser, "vote to a random")
}

func (vi *VotingInstance) ConstructTextForGraphCandidates(out io.Writer, results map[string]float64) {

	counterYesVoter := 0
	counterNoVoter := 0
	counterIndecisiveVoter := 0
	counterThresholdVoter := 0
	counterNormalVoter := 0
	counterNonResponsibleVoter := 0
	counterResponsibleVoter := 0
	for _, user := range vi.GetConfig().Voters {
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

	s := "%"
	fmt.Fprintf(out, "digraph network_activity {\n")
	fmt.Fprintf(out, "labelloc=\"t\";")
	fmt.Fprintf(out, "label = <Votation Diagram of %d nodes.   Results are : ", len(vi.GetConfig().Voters)+len(vi.GetConfig().Candidates))
	for _, cand := range vi.GetConfig().Candidates {
		fmt.Fprintf(out, "%s = %.4v %s,", cand.CandidateID, results[cand.CandidateID], s)
	}
	fmt.Fprintf(out, "<font point-size='10'><br/>(generated: %s)<br/> Il y a %v YesVoter, %v Threshold Voters, %v Non responsibleVoter, %v ResponsibleVoter, %v IndecisiveVoter and %v NormalVoter</font>>; ", time.Now(), counterYesVoter, counterThresholdVoter, counterNonResponsibleVoter, counterResponsibleVoter, counterIndecisiveVoter, counterNormalVoter)
	fmt.Fprintf(out, "graph [fontname = \"helvetica\"];\n")
	fmt.Fprintf(out, "{\n")
	fmt.Fprintf(out, "node [fontname = \"helvetica\" area = 10 style= filled]\n")

	for j, user := range vi.GetConfig().Voters {
		colorOfUser := "black"
		if user.TypeOfUser == "YesVoter" { //YesVoter
			colorOfUser = "darkolivegreen"
		} else if user.TypeOfUser == "NoVoter" { //NoVoter
			colorOfUser = "darkorange1"
		} else if user.TypeOfUser == "IndecisiveVoter" { //IndecisiveVoter
			colorOfUser = "seashell4"
		} else if user.TypeOfUser == "ThresholdVoter" { //ThresholdVoter
			colorOfUser = "gold2"
		} else if user.TypeOfUser == "NonResponsibleVoter" { //NonResponsibleVoter
			colorOfUser = "hotpink1"
		} else if user.TypeOfUser == "ResponsibleVoter" { //ResponsibleVoter
			colorOfUser = "deepskyblue3"
		} else { //NormalVoter
			colorOfUser = "white"
		}
		s := strconv.FormatInt(int64(j), 10)
		fmt.Fprintf(out, "user%s [fillcolor=\"%s\" label=\"user%s\"]\n", s, colorOfUser, s)
	}
	listOfCOlorCand := []string{"", "magenta", "firebrick", "darkseagreen", "darkolivegreen"}
	for k, cand := range vi.GetConfig().Candidates {
		colorOfCand := listOfCOlorCand[int(math.Ceil(float64(k+1)/2))] + strconv.FormatInt(int64(2*k%4)+1, 10)
		fmt.Fprintf(out, "%s [fillcolor=\"%s\" label=\"%s\"]\n", cand.CandidateID, colorOfCand, cand.CandidateID)
	}
	fmt.Fprintf(out, "}\n")
	fmt.Fprintf(out, "edge [fontname = \"helvetica\"];\n")

	for _, user := range vi.GetConfig().Voters {

		colorDeleg := "#8A2BE2"
		colorOfUser := "black"
		if user.TypeOfUser == "YesVoter" { //YesVoter
			colorOfUser = "darkolivegreen"
		} else if user.TypeOfUser == "NoVoter" { //NoVoter
			colorOfUser = "darkorange1"
		} else if user.TypeOfUser == "IndecisiveVoter" { //IndecisiveVoter
			colorOfUser = "seashell4"
		} else if user.TypeOfUser == "ThresholdVoter" { //ThresholdVoter
			colorOfUser = "gold2"
		} else if user.TypeOfUser == "NonResponsibleVoter" { //NonResponsibleVoter
			colorOfUser = "hotpink1"
		} else if user.TypeOfUser == "ResponsibleVoter" { //ResponsibleVoter
			colorOfUser = "deepskyblue3"
		} else { //NormalVoter
			colorOfUser = "black"
		}

		//bruteforce tab with all different possible colors
		var color = []string{"#12B2F5", "#2AEF56", "#FF78EC", "#EAC224", "#F53024", "#A107DE", "#112AE8", "#FF8F00"}

		//creation d'un tableau qui a les cumulative values (plus simple pour le graph)
		cumulativeHistoryOfChoice := make([]voting.Choice, 0)
		new_vote_value := make(map[string]voting.Liquid)
		for _, choice := range user.HistoryOfChoice {
			for name, value := range choice.VoteValue {
				var err error
				new_vote_value[name], err = AddLiquid(new_vote_value[name], value)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		new_choice, err := NewChoice(new_vote_value)
		if err != nil {
			fmt.Println(err.Error())
		}
		cumulativeHistoryOfChoice = append(cumulativeHistoryOfChoice, new_choice)

		for index, cand := range vi.GetConfig().Candidates {
			for _, choice := range cumulativeHistoryOfChoice {
				if choice.VoteValue[cand.CandidateID].Percentage != 0. {
					sizeArrow := 1
					if choice.VoteValue[cand.CandidateID].Percentage > 400 {
						sizeArrow = 400
					} else if choice.VoteValue[cand.CandidateID].Percentage < 20 {
						sizeArrow = 20
					} else {
						sizeArrow = choice.VoteValue[cand.CandidateID].Percentage
					}
					fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
						"[ label = < <font color='%v'><b>%v</b></font><br/>> color=\"%s\" penwidth=%v];\n",
						user.UserID, cand.CandidateID, colorOfUser, choice.VoteValue[cand.CandidateID].Percentage, color[index%len(color)], sizeArrow/20)
				}
			}
		}

		for other, quantity := range user.DelegatedTo {
			sizeArrow := 1
			if quantity.Percentage > 400 {
				sizeArrow = 400
			} else if quantity.Percentage < 20 {
				sizeArrow = 20
			} else {
				sizeArrow = quantity.Percentage
			}
			fmt.Fprintf(out, "\"%v\" -> \"%v\" "+
				"[ label = < <font color='#8A2BE2'><b>%v</b></font><br/>> color=\"%s\" penwidth=%v];\n",
				user.UserID, other, quantity.Percentage, colorDeleg, sizeArrow/20)
		}
	}

	fmt.Fprintf(out, "}\n")

}

//Write the necessary information for GraphViz in the buffer given in parameter
func (vi *VotingInstance) ConstructTextForGraph(out io.Writer) {

	counterYesVoter := 0
	counterNoVoter := 0
	counterIndecisiveVoter := 0
	counterThresholdVoter := 0
	counterNormalVoter := 0
	counterNonResponsibleVoter := 0
	counterResponsibleVoter := 0
	for _, user := range vi.GetConfig().Voters {
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

	results := vi.GetResults()
	s := "%"
	fmt.Fprintf(out, "digraph network_activity {\n")
	fmt.Fprintf(out, "labelloc=\"t\";")
	fmt.Fprintf(out, "label = <Votation Diagram of %d nodes.    Results are Yes = %.4v %s, No = %.4v %s<font point-size='10'><br/>(generated: %s)<br/> Il y a %v YesVoter, %v NoVoter, %v Threshold Voters, %v Non responsibleVoter, %v ResponsibleVoter, %v IndecisiveVoter and %v NormalVoter</font>>; ", len(vi.GetConfig().Voters)+2, results["yes"], s, results["no"], s, time.Now(), counterYesVoter, counterNoVoter, counterThresholdVoter, counterNonResponsibleVoter, counterResponsibleVoter, counterIndecisiveVoter, counterNormalVoter)
	fmt.Fprintf(out, "graph [fontname = \"helvetica\"];\n")
	fmt.Fprintf(out, "{\n")
	fmt.Fprintf(out, "node [fontname = \"helvetica\" area = 10 style= filled]\n")

	for j, user := range vi.GetConfig().Voters {
		colorOfUser := "black"
		if user.TypeOfUser == "YesVoter" { //YesVoter
			colorOfUser = "darkolivegreen"
		} else if user.TypeOfUser == "NoVoter" { //NoVoter
			colorOfUser = "darkorange1"
		} else if user.TypeOfUser == "IndecisiveVoter" { //IndecisiveVoter
			colorOfUser = "seashell4"
		} else if user.TypeOfUser == "ThresholdVoter" { //ThresholdVoter
			colorOfUser = "gold2"
		} else if user.TypeOfUser == "NonResponsibleVoter" { //NonResponsibleVoter
			colorOfUser = "hotpink1"
		} else if user.TypeOfUser == "ResponsibleVoter" { //ResponsibleVoter
			colorOfUser = "deepskyblue3"
		} else { //NormalVoter
			colorOfUser = "white"
		}
		s := strconv.FormatInt(int64(j), 10)
		fmt.Fprintf(out, "user%s [fillcolor=\"%s\" label=\"user%s\"]\n", s, colorOfUser, s)
	}
	fmt.Fprintf(out, "}\n")
	fmt.Fprintf(out, "edge [fontname = \"helvetica\"];\n")

	for _, user := range vi.GetConfig().Voters {

		colorVoteYes := "#22bd27"
		colorVoteNo := "#cf1111"
		colorDeleg := "#8A2BE2"

		//creation d'un tableau qui a les cumulative values (plus simple pour le graph)
		cumulativeHistoryOfChoice := make([]voting.Choice, 0)
		new_vote_value := make(map[string]voting.Liquid)
		for _, choice := range user.HistoryOfChoice {
			for name, value := range choice.VoteValue {
				var err error
				new_vote_value[name], err = AddLiquid(new_vote_value[name], value)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		new_choice, err := NewChoice(new_vote_value)
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

//Transfer of voting power between 2 users
func (vi *VotingInstance) DelegTo(userSend *voting.User, userReceive *voting.User, quantity voting.Liquid) error {

	//CANNOT DELEGATE 0 voting power, do nothing if it is the case
	if quantity.Percentage > 0 {
		new_quantity, err := AddLiquid(quantity, userReceive.DelegatedFrom[userSend.UserID])
		if err != nil {
			return xerrors.Errorf(err.Error())
		}
		userReceive.DelegatedFrom[userSend.UserID] = new_quantity

		userSend.DelegatedTo[userReceive.UserID] = new_quantity

		userReceive.VotingPower += quantity.Percentage
		err = vi.CheckVotingPower(userReceive)
		if err != nil {
			return xerrors.Errorf(err.Error())
		}

		userSend.VotingPower -= quantity.Percentage
		err = vi.CheckVotingPower(userSend)
		if err != nil {
			return xerrors.Errorf(err.Error())
		}
	}

	return nil
}

//VOTING SYSTEM FUNCTIONS :::::

type VotingSystem struct {
	//contain all the votingInstances mapped to their stringID

	//66666666666666666666666666666 ESSAYER DE RETIRER LE POINTEUR
	VotingInstancesList map[string]voting.VotingInstance

	//database
	Database storage.DB
}

//Returns the voting instance list
func (vs VotingSystem) GetVotingInstanceList() map[string]voting.VotingInstance {
	return vs.VotingInstancesList
}

//Creation of a voting system, passing db and map as arguments
func NewVotingSystem(db storage.DB, votingInstanceList map[string]voting.VotingInstance) VotingSystem {
	return VotingSystem{
		Database:            db,
		VotingInstancesList: votingInstanceList,
	}
}

//Creates and add new a voting instance
func (vs VotingSystem) CreateAndAdd(id string, config voting.VotingConfig, status string) (voting.VotingInstance, error) {

	//check if id is null
	if id == "" {
		return nil, xerrors.Errorf("The id is empty.")
	}

	//check if status is open or close only
	if !(status == "open") && !(status == "close") {
		return nil, xerrors.Errorf("The status is incorrect, should be either 'open' or 'close'.")
	}

	//fmt.Println("Votes: ", votes)

	//create the object votingInstance
	var vi = VotingInstance{
		Id:     id,
		Config: config,
		Status: status,
	}

	p := &vi
	*p = vi

	//adding vi to the list of vi's of the voting system
	vs.VotingInstancesList[id] = p

	return p, nil
}

func (vi VotingInstance) GetVotingID() string {
	return vi.Id
}

//Delete the voting instance linked to the id
func (vs VotingSystem) Delete(id string) error {

	vi := vs.VotingInstancesList[id]
	if vi.GetStatus() == "open" {
		//vi.Status = "close
		return xerrors.Errorf("Can't delete the votingInsance because it is still open")
	} else {
		delete(vs.VotingInstancesList, id)
	}
	return nil
}

//Return a list of all the voting instance that are open
func (vs VotingSystem) ListVotings() []string {
	listeDeVotes := make([]string, 0)
	for key := range vs.VotingInstancesList {
		if vs.VotingInstancesList[key].GetStatus() == "open" {
			listeDeVotes = append(listeDeVotes, key)
		}
	}
	return listeDeVotes
}

//Return the voting instance from the id
func (vs VotingSystem) GetVotingInstance(id string) voting.VotingInstance {
	return vs.VotingInstancesList[id]
}

//Create and return a new voting configuration
func NewVotingConfig(voters []*voting.User, title string, desc string, cand []*voting.Candidate, typeOfVotingConfig voting.TypeOfVotingConfig) (voting.VotingConfig, error) {
	if title == "" {
		return voting.VotingConfig{}, xerrors.Errorf("Title is empty")
	}
	if typeOfVotingConfig != "CandidateQuestion" && typeOfVotingConfig != "YesOrNoQuestion" {
		return voting.VotingConfig{}, xerrors.Errorf("TypeOfVotingCOnfig Incorrect")
	}

	return voting.VotingConfig{
		Voters:             voters,
		Title:              title,
		Description:        desc,
		Candidates:         cand,
		TypeOfVotingConfig: typeOfVotingConfig,
	}, nil
}

//Create and return a new User
func (vs VotingSystem) NewUser(userID string, delegTo map[string]voting.Liquid, delegFrom map[string]voting.Liquid, historyOfChoice []voting.Choice, typeOfUser voting.TypeOfUser, preferenceDelegationList []*voting.User) (voting.User, error) {

	// if votingPower > (float64(delegFrom)+1)*PERCENTAGE {
	// 	return voting.Choice{}, xerrors.Errorf("Voting power is too much : %f", votingPower)
	// }

	// //check that the sum overall votes is less or equal to the voting power
	// var sum float64 = 0
	// for _, value := range deleg {
	// 	sum += value.Percentage
	// }
	// for _, value := range choice {
	// 	sum += value.Percentage
	// }
	// if sum > (votingPower + float64(delegFrom)*PERCENTAGE) {
	// 	return voting.Choice{}, xerrors.Errorf("Cumulate voting power distributed is greater than the voting power. Was: %f, must not be greater thant %f.", sum, votingPower)
	// }

	return voting.User{
		UserID:                   userID,
		DelegatedTo:              delegTo,
		DelegatedFrom:            delegFrom,
		VotingPower:              PERCENTAGE,
		HistoryOfChoice:          historyOfChoice,
		TypeOfUser:               typeOfUser,
		PreferenceDelegationList: preferenceDelegationList,
	}, nil
}

//Create and return a new Candidate
func (vs VotingSystem) NewCandidate(candidateID string) (voting.Candidate, error) {

	return voting.Candidate{
		CandidateID: candidateID,
	}, nil
}

//Create and return a new Choice
func NewChoice(voteValue map[string]voting.Liquid) (voting.Choice, error) {
	return voting.Choice{
		VoteValue: voteValue,
	}, nil
}

//Create and return a new Liquid
func NewLiquid(p int) (voting.Liquid, error) {
	if p < 0 {
		return voting.Liquid{}, xerrors.Errorf("Init value is incorrect: was %d, must be positive.", int(p))
	}

	return voting.Liquid{
		Percentage: p,
	}, nil
}

//Return the addition of 2 liquids
func AddLiquid(l1 voting.Liquid, l2 voting.Liquid) (voting.Liquid, error) {
	result, err := NewLiquid(int(l1.Percentage + l2.Percentage))
	if err != nil {
		return voting.Liquid{}, xerrors.Errorf(err.Error())
	}
	return result, err
}

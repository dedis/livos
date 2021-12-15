package impl

import (
	"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting"
	"golang.org/x/xerrors"
)

//VOTING IMPLEMENTATION

const PERCENTAGE = 100

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

//Give the result of the choices of the voting instance in the form: map[no:50 yes:50]
func (vi *VotingInstance) GetResults() map[string]float64 {
	if vi.Config.TypeOfVotingConfig == "YesOrNoQuestion" {
		results := make(map[string]float64, len(vi.Config.Voters))
		var yesPower float64 = 0
		var noPower float64 = 0
		for _, user := range vi.Config.Voters {
			for _, choice := range user.HistoryOfChoice {
				yesPower += choice.VoteValue["yes"].Percentage
				noPower += choice.VoteValue["no"].Percentage
			}
		}
		//in order to get 4 and not 4.6666666... for example
		// var temp1 = float64(int(yesPower/float64(counter))*100) / 100
		// var temp2 = float64(int(noPower/float64(counter))*100) / 100
		results["yes"] = yesPower / float64(len(vi.Config.Voters))
		results["no"] = noPower / float64(len(vi.Config.Voters))

		return results
	} else {
		resultsCandidate := make(map[string]float64, len(vi.Config.Candidates))

		for _, candidate := range vi.Config.Candidates {
			var candidateResult float64 = 0
			for _, user := range vi.Config.Voters {
				for _, choice := range user.HistoryOfChoice {
					candidateResult += choice.VoteValue[candidate.CandidateID].Percentage
				}
			}
			resultsCandidate[candidate.CandidateID] = candidateResult / float64(len(vi.Config.Voters))
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

	var sumOfVotingPower float64 = 0
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
func NewLiquid(p float64) (voting.Liquid, error) {
	if p < 0 {
		return voting.Liquid{}, xerrors.Errorf("Init value is incorrect: was %d, must be positive.", int(p))
	}

	return voting.Liquid{
		Percentage: p,
	}, nil
}

//Return the addition of 2 liquids
func AddLiquid(l1 voting.Liquid, l2 voting.Liquid) (voting.Liquid, error) {
	result, err := NewLiquid(l1.Percentage + l2.Percentage)
	if err != nil {
		return voting.Liquid{}, xerrors.Errorf(err.Error())
	}
	return result, err
}

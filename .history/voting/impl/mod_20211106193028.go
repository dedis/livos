package impl

import (
	"fmt"

	"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting"
	"golang.org/x/xerrors"
)

//the voting implementation

const PERCENTAGE = 100

type VotingInstance struct {
	//voting instance's id
	Id string

	//parameters of the Voting
	Config voting.VotingConfig

	// open / closed
	Status string

	// Votes contains the choice of each voter, references by a userID
	Votes map[string]*voting.Choice

	//db.socket personnalisé pour chacun ?
}

func (vi *VotingInstance) CastVote(userID string, choice voting.Choice) error {
	if vi.Status == "close" {
		return xerrors.Errorf("Impossible the cast the vote, the voting instance is closed.")
	}
	fmt.Println("CHOICE", choice)
	fmt.Println("userID", userID)
	fmt.Println("vi:", *vi)
	vi.Votes[userID] = &voting.Choice{}
	vi.Votes[userID] = &choice
	return nil
}

func (vi VotingInstance) CloseVoting() {
	vi.SetStatus("close")
}

func (vi *VotingInstance) SetStatus(status string) error {
	if (status != "open") || (status != "close") {
		return xerrors.Errorf("The status is incorrect. Should be either 'open' or 'close'. Was: %s", status)
	}
	vi.Status = status
	return nil
}

func (vi *VotingInstance) GetStatus() string {
	return vi.Status
}

//mieux de garder pointeur ?
func (vi VotingInstance) GetConfig() voting.VotingConfig {
	return vi.Config
}

//Give the result of the choices of the voting instance in the form: map[no:50 yes:50]
func (vi VotingInstance) GetResults() map[string]float32 {
	results := make(map[string]float32, len(vi.Votes))
	counter := 0
	var yesPower float32 = 0
	var noPower float32 = 0
	for _, v := range vi.Votes {
		yesPower += v.MyChoice["yes"].Percentage
		noPower += v.MyChoice["no"].Percentage
		counter++
	}
	results["yes"] = yesPower / float32(counter)
	results["no"] = noPower / float32(counter)

	return results
}

type VotingSystem struct {
	//contain all the votingInstances mapped to their stringID
	VotingInstancesList map[string]*VotingInstance

	//database
	Database storage.DB
}

//creation of a voting system, passing db and map as arguments
func NewVotingSystem(db storage.DB, vil map[string]*VotingInstance) VotingSystem {
	return VotingSystem{
		Database:            db,
		VotingInstancesList: vil,
	}
}

//creation of a voting instance
func (vs VotingSystem) CreateAndAdd(id string, config voting.VotingConfig, status string, votes map[string]*voting.Choice) (VotingInstance, error) {

	//check if id is null
	if id == "" {
		return VotingInstance{}, xerrors.Errorf("The id is empty.")
	}

	//check if status is open or close only
	// if status != "open" || status != "close" {
	// 	return VotingInstance{}, xerrors.Errorf("The status is incorrect, should be either 'open' or 'close'.")
	// }

	//fmt.Println("Votes: ", votes)

	//create the object votingInstance
	var vi = VotingInstance{
		Id:     id,
		Config: config,
		Status: status,
		Votes:  votes,
	}

	p := &vi
	*p = vi

	//adding vi to the list of vi's of the voting system
	vs.VotingInstancesList[id] = p

	return *p, nil
}

func (vs VotingSystem) Delete(id string) error {

	vi := vs.VotingInstancesList[id]
	if vi.Status == "open" {
		//vi.Status = "close"
		return xerrors.Errorf("Can't delete the votingInsance because it is still open")
	} else {
		delete(vs.VotingInstancesList, id)
	}
	return nil
}

//Return a list of all the voting instance
func (vs VotingSystem) ListVotings() []string {
	listeDeVotes := make([]string, len(vs.VotingInstancesList))
	for key := range vs.VotingInstancesList {
		if vs.VotingInstancesList[key].Status == "open" {
			listeDeVotes = append(listeDeVotes, key)
		}
	}
	return listeDeVotes
}

//Do we need to make a check to see if the id is null or letters or in fact
//doesn't belong to the list of ids
func (vs VotingSystem) GetVotingInstance(id string) VotingInstance {
	return *vs.VotingInstancesList[id]
}

func NewVotingConfig(voters []string, title string, desc string, cand []string) (voting.VotingConfig, error) {
	if title == "" {
		return voting.VotingConfig{}, xerrors.Errorf("title is empty")
	}

	return voting.VotingConfig{
		Voters:      voters,
		Title:       title,
		Description: desc,
		Candidates:  cand,
	}, nil
}

func NewChoice(deleg map[string]voting.Liquid, choice map[string]voting.Liquid, delegFrom int, votingPower float32) (voting.Choice, error) {
	if delegFrom < 0 {
		return voting.Choice{}, xerrors.Errorf("Delegation number received is negative : %d", delegFrom)
	}

	if votingPower > (float32(delegFrom)+1)*PERCENTAGE {
		return voting.Choice{}, xerrors.Errorf("Voting power is too much : %f", votingPower)
	}

	//check that the sum overall votes is less or equal to the voting power
	var sum float32 = 0
	for _, value := range deleg {
		sum += value.Percentage
	}
	for _, value := range choice {
		sum += value.Percentage
	}
	if sum > (votingPower + float32(delegFrom)*PERCENTAGE) {
		return voting.Choice{}, xerrors.Errorf("Cumulate voting power distributed is greater than the voting power. Was: %f, must not be greater thant %f.", sum, votingPower)
	}

	return voting.Choice{
		DelegatedTo:   deleg,
		MyChoice:      choice,
		DelegatedFrom: delegFrom,
		VotingPower:   votingPower,
	}, nil
}

func NewLiquid(p float32) (voting.Liquid, error) {
	if p > 100 {
		return voting.Liquid{}, xerrors.Errorf("Init value is incorrect: Was %f, must be less than %d", p, PERCENTAGE)
	}

	return voting.Liquid{
		Percentage: p,
	}, nil
}

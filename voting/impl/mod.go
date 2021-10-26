package impl

import (
	"fmt"

	"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting"
)

//the voting implementation

type VotingInstance struct {
	//voting instance's id
	Id string

	//parameters of the Voting
	Config voting.VotingConfig

	// open / closed
	Status string

	// Votes contains the choice of each voter, references by a userID
	Votes map[string]voting.Choice

	//db.socket personnalisé pour chacun ?
}

func (vi VotingInstance) CastVote(userID string, choice voting.Choice) {
	vi.Votes[userID] = choice
}

func (vi VotingInstance) CloseVoting(id string) {
	vi.SetStatus("close")
}

func (vi *VotingInstance) SetStatus(status string) {
	vi.Status = status
}

func (vi *VotingInstance) GetStatus(status string) string {
	return vi.Status
}

//mieux de garder pointeur ?
func (vi VotingInstance) GetConfig() voting.VotingConfig {
	return vi.Config
}

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
	VotingInstancesList map[string]VotingInstance

	//database
	Database storage.DB
}

//creation of a voting system, passing db and map as arguments
func NewVotingSystem(db storage.DB, vil map[string]VotingInstance) VotingSystem {
	return VotingSystem{
		Database:            db,
		VotingInstancesList: vil,
	}
}

//creation of a voting instance
func (vs VotingSystem) Create(id string, config voting.VotingConfig, status string, votes map[string]voting.Choice) VotingInstance {
	//create the object votingInstance
	var vi = VotingInstance{
		Id:     id,
		Config: config,
		Status: status,
		Votes:  votes,
	}

	//adding vi to the list of vi's of the voting system
	vs.VotingInstancesList[id] = vi

	return vi
}

func (vs VotingSystem) Delete(id string) {
	vi := vs.VotingInstancesList[id]
	if vi.Status == "open" {
		//vi.Status = "close"
		fmt.Println("Can't delete the votingInsance because it is still open")
	} else {
		delete(vs.VotingInstancesList, id)
	}
}

func (vs VotingSystem) ListVotings() []string {
	listeDeVotes := make([]string, len(vs.VotingInstancesList))
	for key := range vs.VotingInstancesList {
		listeDeVotes = append(listeDeVotes, key)
	}
	return listeDeVotes

}

func (vs VotingSystem) GetVotingInstance(id string) VotingInstance {
	return vs.VotingInstancesList[id]
}

func NewVotingConfig(voters []string, title string, desc string, cand []string) voting.VotingConfig {
	return voting.VotingConfig{
		Voters:      voters,
		Title:       title,
		Description: desc,
		Candidates:  cand,
	}
}

func NewChoice(deleg map[string]voting.Liquid, choice map[string]voting.Liquid, vp float32) voting.Choice {
	return voting.Choice{
		DelegatedTo: deleg,
		MyChoice:    choice,
		VotingPower: vp,
	}
}

func NewLiquid(p float32) voting.Liquid {
	return voting.Liquid{
		Percentage: p,
	}
}

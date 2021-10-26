package impl

import (
	"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting"
)

//the voting implementation

// how to use the stuct defined in the upper mod.go file ?

//listVote is map[string]Voting, containing all the voting sessions

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

func CastVote(votingID, userID string, choice voting.Choice) {
	//listVote.GetVoting(votingID).db.update(userID, Choice)
}

func CloseVoting(votingID string) {
	//listVote.GetVoting(votingID).close()
}

func GetVoting(votingID string) VotingInstance {
	// return listVote.get(votingID)
	return VotingInstance{}
}

//creation of a voting system, passing db and map as arguments
func NewVotingSystem(db storage.DB, vil map[string]VotingInstance) VotingSystem {
	return VotingSystem{
		Database:            db,
		VotingInstancesList: vil,
	}
}

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

type VotingSystem struct {
	//contain all the votingInstances mapped to their stringID
	VotingInstancesList map[string]VotingInstance

	//database
	Database storage.DB
}

func NewVotingConfig(voters []string, title string, desc string, cand []string) voting.VotingConfig {
	return voting.VotingConfig{
		Voters:      voters,
		Title:       title,
		Description: desc,
		Candidates:  cand,
	}
}

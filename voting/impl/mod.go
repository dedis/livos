package impl

import (
	"github.com/dedis/livos/storage"
	"github.com/dedis/livos/voting"
)

//the voting implementation

// how to use the stuct defined in the upper mod.go file ?

//listVote is map[string]Voting, containing all the voting sessions

func Create(config voting.VotingConfig) string {
	return config.Title //to encrypt to get a unique votingID?
}

func CastVote(votingID, userID string, choice voting.Choice) {
	//listVote.GetVoting(votingID).db.update(userID, Choice)
}

func CloseVoting(votingID string) {
	//listVote.GetVoting(votingID).close()
}

func GetVoting(votingID string) Voting {
	// return listVote.get(votingID)
	return Voting{}
}

func (c Voting) newVoting(vc voting.VotingConfig, status string, votes map[string]voting.Choice, db storage.DB) Voting {
	return Voting{
		Config:   vc,
		Status:   status,
		Votes:    votes,
		Database: db,
	}
}

func newVotingConfig(voters []string, title string, desc string, cand []string) voting.VotingConfig {
	return voting.VotingConfig{
		Voters:      voters,
		Title:       title,
		Description: desc,
		Candidates:  cand,
	}
}

type Voting struct {
	//parameters of the Voting
	Config voting.VotingConfig

	// open / closed
	Status string

	// Votes contains the choice of each voter, references by a userID
	Votes map[string]voting.Choice

	//database
	Database storage.DB
}

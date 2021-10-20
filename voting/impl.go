package voting

import "github.com/dedis/livos/storage"

//the voting implementation

// how to use the stuct defined in the upper mod.go file ?

//listVote is map[string]Voting, containing all the voting sessions

func Create(config VotingConfig) string {
	return config.Title //to encrypt to get a unique votingID?
}

func CastVote(votingID, userID string, choice Choice) {
	//listVote.GetVoting(votingID).db.update(userID, Choice)
}

func CloseVoting(votingID string) {
	//listVote.GetVoting(votingID).close()
}

func GetVoting(votingID string) Voting {
	// return listVote.get(votingID)
	return Voting{}
}

func newVoting(vc VotingConfig, status string, votes map[string]Choice, db storage.DB) Voting {
	return Voting{
		Config:   vc,
		Status:   status,
		Votes:    votes,
		Database: db,
	}
}

func newVotingConfig(voters []string, title string, desc string, cand []string) VotingConfig {
	return VotingConfig{
		Voters:      voters,
		Title:       title,
		Description: desc,
		Candidates:  cand,
	}
}

package voting

import "github.com/dedis/livos/storage"

//the voting interface and struct definition

type VotingSystem interface {
	// Create returns the unique votingID
	Create(config VotingConfig) string

	CastVote(votingID, userID string, choice Choice)

	CloseVoting(votingID string)

	GetVoting(votingID string) Voting
}

type VotingConfig struct {
	// Voters is a list of userID
	Voters      []string
	Title       string
	Description string

	//isYesNo bool

	// Candidates is a list of userID (can be empty if yes/no question)
	Candidates []string
}

type Choice struct {
	// DelegatedTo contains the percentage of the voting power given to the different delegates (represented by a userID)
	// or is empty if there is no delegation
	DelegatedTo map[string]Liquid

	// MyChoice contains map of userID to the percentage of voting power, or is empty if there is a delegation
	MyChoice map[string]Liquid

	// VotingPower contains how many voting percentage is left ?
	VotingPower float32
}

type Liquid struct {
	Percentage float32
}

type Voting struct {
	//parameters of the Voting
	Config VotingConfig

	// open / closed
	Status string

	// Votes contains the choice of each voter, references by a userID
	Votes map[string]Choice

	//database
	Database storage.DB
}

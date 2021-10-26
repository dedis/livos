package voting

//the voting interface and struct definition

type VotingSystem interface {
	ListVotings() []string

	GetVotingInstance(votingID string) VotingInstance

	Create(config VotingConfig) VotingInstance

	Delete(votingID string)
}

type VotingInstance interface {
	CastVote(userID string, choice Choice)

	GetConfig() VotingConfig

	CloseVoting()

	GetResults() Results
}

type VotingConfig struct {
	// Voters is a list of userID
	Voters      []string
	Title       string
	Description string

	// Candidates is a list of userID (can be empty if yes/no question)
	Candidates []string
}

type Choice struct {
	// DelegatedTo contains the percentage of the voting power given to the
	// different delegates (represented by a userID) or is empty if there is no
	// delegation
	DelegatedTo map[string]Liquid

	// MyChoice contains map of userID to the percentage of voting power, or is empty if there is a delegation
	MyChoice map[string]Liquid

	// VotingPower contains how many voting percentage is left ?
	VotingPower float32
}

type Results struct {
}

type Liquid struct {
	Percentage float32
}

package voting

//the voting interface and struct definition

type VotingSystem interface {
	ListVotings() []string

	GetVotingInstance(votingID string) VotingInstance

	Create(config VotingConfig) VotingInstance

	Delete(votingID string)

	//override the method print?
}

type VotingInstance interface {
	CastVote(userID string, choice Choice)

	GetConfig() VotingConfig

	CloseVoting()

	GetResults() map[string]float32

	SetStatus(status string)

	GetStatus() string

	//activate (open)?

	//override the method print?
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

	// MyChoice contains map the YES/NO answer to the percentage of voting
	// power, or is empty if there is a delegation
	MyChoice map[string]Liquid

	// Number of delegation power received
	DelegatedFrom int

	// VotingPower contains how many voting percentage is left ?
	VotingPower float64
}

//for liquidity and delegation
type Liquid struct {
	Percentage float64
}

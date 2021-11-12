package voting

//the voting interface and struct definition

type VotingSystem interface {
	ListVotings() []string

	GetVotingInstance(votingID string) VotingInstance

	CreateAndAdd(config VotingConfig) VotingInstance

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

	CheckVotingPower(user *User) bool

	//activate (open)?

	//override the method print?
}

type VotingConfig struct {
	// Voters is a list of userID
	Voters      []*User
	Title       string
	Description string
	// Candidates is a list of userID (can be empty if yes/no question)
	Candidates []string
}

type Choice struct {
	// VoteValue contains map the YES/NO answer to the percentage of voting
	// power, or is empty if there is a delegation
	VoteValue map[string]Liquid
}

//for liquidity and delegation
type Liquid struct {
	Percentage float64
}

type User interface {
	CheckVotingPower() bool

	SetChoice(choice Choice) error

	DelegTo(user *User, quantity Liquid) error
	DelegFrom(user *User, quantity Liquid) error
}

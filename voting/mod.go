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
	CastVote(user *User)

	GetConfig() VotingConfig

	CloseVoting()

	GetResults() map[string]float32

	SetStatus(status string)

	GetStatus() string

	GetUSer(string) (*User, error)

	CheckVotingPower(user *User) bool

	SetChoice(user *User, choice Choice) error

	DelegTo(user *User, other *User, quantity Liquid) error

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

type User struct {
	//name of the user
	UserID string

	//keep the record of how much was delegated to whom
	DelegatedTo map[string]Liquid

	//keep the record of how much was given to self and from who
	DelegatedFrom map[string]Liquid

	//choice of the user concerning the voting instance
	MyChoice Choice

	//the amount of voting still left to split btw votes or delegations
	VotingPower float64
}

// type User interface {
// 	CheckVotingPower() bool

// 	SetChoice(choice Choice) error

// 	DelegTo(user *User, quantity Liquid) error

// 	DelegFrom(user *User, quantity Liquid) error
// }

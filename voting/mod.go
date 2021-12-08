package voting

//the voting interface and struct definition

type VotingSystem interface {
	ListVotings() []string

	GetVotingInstance(id string) VotingInstance

	CreateAndAdd(id string, config VotingConfig, status string) (VotingInstance, error)

	Delete(votingID string) error

	GetVotingInstanceList() map[string]VotingInstance

	NewUser(userID string, delegTo map[string]Liquid, delegFrom map[string]Liquid, histoChoice []Choice, typeOfUser TypeOfUser, preferenceDelegationList []*User) (User, error)

	//override the method print?
}

type VotingInstance interface {
	GetConfig() VotingConfig

	CloseVoting()

	GetResults() map[string]float64

	SetStatus(status string) error

	GetStatus() string

	GetVotingID() string

	SetTitle(title string) error

	SetDescription(description string) error

	SetVoters(users []*User) error

	SetCandidates(candidates []string) error

	GetUser(string) (*User, error)

	CheckVotingPower(user *User) error

	CheckVotingPowerOfUser(user *User) bool

	SetVote(user *User, choice Choice) error

	DelegTo(user *User, other *User, quantity Liquid) error

	CheckVotingPowerOfVoters() bool

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

	//the amount of voting still left to split btw votes or delegations
	VotingPower float64

	//history of choices that were cast
	HistoryOfChoice []Choice

	//type define behavior of the user
	TypeOfUser TypeOfUser

	//delegation preference list
	PreferenceDelegationList []*User
}

type TypeOfUser string

const (
	YesVoter       TypeOfUser = "YesVoter"
	NoVoter        TypeOfUser = "NoVoter"
	IndeciseVoter  TypeOfUser = "IndeciseVoter"
	ThresholdVoter TypeOfUser = "ThresholdVoter"
	None           TypeOfUser = "None"
)

// func (t TypeOfUser) String() string {
// 	switch t {
// 	case YesVoter:
// 		return "YesVoter"
// 	case NoVoter:
// 		return "NoVoter"
// 	case IndeciseVoter:
// 		return "IndeciseVoter"
// 	case ThresholdVoter:
// 		return "ThresholdVoter"
// 	default:
// 		return "None"
// 	}
// }

// type User interface {
// 	CheckVotingPower() bool

// 	SetChoice(choice Choice) error

// 	DelegTo(user *User, quantity Liquid) error

// 	DelegFrom(user *User, quantity Liquid) error
// }

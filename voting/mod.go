package voting

import (
	"io"
)

//the voting interface and struct definition

type VotingSystem interface {
	ListVotings() []string

	GetVotingInstance(id string) VotingInstance

	CreateAndAdd(id string, config VotingConfig, status string) (VotingInstance, error)

	Delete(votingID string) error

	GetVotingInstanceList() map[string]VotingInstance

	NewUser(userID string, delegTo map[string]Liquid, delegFrom map[string]Liquid, histoChoice []Choice, typeOfUser TypeOfUser, preferenceDelegationList []*User) (User, error)

	NewCandidate(candidateID string) (Candidate, error)

	//override the method print?
}

type VotingInstance interface {
	GetConfig() VotingConfig

	SetConfig(config VotingConfig)

	CloseVoting()

	GetResults() map[string]float64

	GetResultsQuadraticVoting() map[string]float64

	ConstructTextForGraph(out io.Writer)

	ConstructTextForGraphCandidates(out io.Writer, results map[string]float64)

	SplitVPintoActions(user *User, i int, votingPowerToSplit int)

	SplitVotingPower(user *User, i int, votingPowerToSplit int) (int, int, int)

	RandomDelegate(user *User, i int, votingPower int)

	YesVote(user *User, votingPower int)

	NoVote(user *User, votingPower int)

	IndecisiveVote(user *User, i int, quantityToDeleg int)

	RandomVote(user *User, i int)

	ThresholdVote(user *User, i int, threshold int)

	NonResponsibleVote(user *User, i int)

	ResponsibleVote(user *User, i int)

	CandidateVote(user *User, i int, votingPower int)

	BreakTheCycleCandidate(user *User, i int, votingPower int)

	IndecisiveVoteCandidate(user *User, i int, quantityToDeleg int)

	DefaultVoteCandidate(user *User, i int)

	ThresholdVoteCandidate(user *User, i int, threshold int, votingPower int)

	NonResponsibleVoteCandidate(user *User, i int, votingPower int)

	ResponsibleVoteCandidate(user *User, i int, votingPower int)

	RandomWithProbabilities(user *User) int

	SetStatus(status string) error

	GetStatus() string

	GetVotingID() string

	SetTitle(title string) error

	SetDescription(description string) error

	SetVoters(users []*User) error

	SetCandidates(candidates []*Candidate) error

	SetTypeOfVotingConfig(typeOfVotingConfig string) error

	GetUser(string) (*User, error)

	GetCandidate(string) (*Candidate, error)

	CheckVotingPower(user *User) error

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
	Candidates []*Candidate

	//work like a boolean to see if the votingConfig is a yes/no question or a candidate one
	TypeOfVotingConfig TypeOfVotingConfig
}
type TypeOfVotingConfig string

const (
	CandidateQuestion TypeOfVotingConfig = "CandidateQuestion"
	YesOrNoQuestion   TypeOfVotingConfig = "YesOrNoQuestion"
)

type Choice struct {
	// VoteValue contains map the YES/NO answer to the percentage of voting
	// power, or is empty if there is a delegation
	VoteValue map[string]Liquid
}

//for liquidity and delegation
type Liquid struct {
	Percentage int
}

type User struct {
	//name of the user
	UserID string

	//keep the record of how much was delegated to whom
	DelegatedTo map[string]Liquid

	//keep the record of how much was given to self and from who
	DelegatedFrom map[string]Liquid

	//the amount of voting still left to split btw votes or delegations
	VotingPower int

	//history of choices that were cast
	HistoryOfChoice []Choice

	//type define behavior of the user
	TypeOfUser TypeOfUser

	//delegation preference list
	PreferenceDelegationList []*User
}

type Candidate struct {
	//name of the candidat
	CandidateID string

	//It could aso have a party, a programm, a type of candidate
}

type TypeOfUser string

const (
	YesVoter            TypeOfUser = "YesVoter"
	NoVoter             TypeOfUser = "NoVoter"
	CandVoter           TypeOfUser = "CandVoter"
	IndecisiveVoter     TypeOfUser = "IndecisiveVoter"
	ThresholdVoter      TypeOfUser = "ThresholdVoter"
	NonResponsibleVoter TypeOfUser = "NonResponsibleVoter"
	ResponsibleVoter    TypeOfUser = "ResponsibleVoter"
	None                TypeOfUser = "None"
)

//VOTING CONFIG FUNCTIONS :::::

func (vc VotingConfig) SetCandidates(new_candidates []*Candidate) VotingConfig {
	vc.Candidates = new_candidates
	return vc
}

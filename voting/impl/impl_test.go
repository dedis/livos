package impl

import (
	"testing"

	"github.com/dedis/livos/voting"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

//VARIABLES NECESSAIRES
var VoteList = make(map[string]*VotingInstance)
var VoteSystem = NewVotingSystem(nil, VoteList)

//Creation of a empty list of choces (for history)
var histoChoice = make([]voting.Choice, 0)

var userNoemien, _ = VoteSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice)
var userGuillaume, _ = VoteSystem.NewUser("Guillaume", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice)
var userEtienne, _ = VoteSystem.NewUser("Etienne", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice)

var voters = []*voting.User{&userNoemien, &userGuillaume, &userEtienne}
var candidats = make([]string, 0)

var voteConfig, _ = NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)

var vi, _ = VoteSystem.CreateAndAdd("Session01", voteConfig, "open")

var liquid_150, _ = NewLiquid(150)
var liquid_100, _ = NewLiquid(100)
var liquid_50, _ = NewLiquid(50)
var liquid_0, _ = NewLiquid(0)

var yesChoice = make(map[string]voting.Liquid)
var noChoice = make(map[string]voting.Liquid)
var midChoice = make(map[string]voting.Liquid)

func TestNewUser(t *testing.T) {
	userNoemien, err := VoteSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice)
	if err != nil {
		xerrors.Errorf(err.Error())
	}
	require.Equal(t, userNoemien.UserID, "Noemien", "UserID initialization incorrect : was %s, should be %s", userNoemien.UserID, "Noemien")
	require.Equal(t, userNoemien.VotingPower, 100., "VotingPower initialization incorrect : was %f, should be %f", userNoemien.VotingPower, 100.)
	require.Equal(t, userNoemien.HistoryOfChoice, make([]voting.Choice, 0), "HistoryOfChoice initialization incorrect")
	require.Equal(t, userNoemien.DelegatedFrom, make(map[string]voting.Liquid), "DelegatedFrom initialization incorrect")
	require.Equal(t, userNoemien.DelegatedTo, make(map[string]voting.Liquid), "DelegatedTo initialization incorrect")
}

//TODO : add testing error for newUser()

func TestVotingInstanceCreate(t *testing.T) {
	voteConfig, err := NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	require.Equal(t, err, nil, "Cannot create VotingConfig")

	VoteSystem.CreateAndAdd("Session01", voteConfig, "open")
	id := VoteSystem.VotingInstancesList["Session01"].Id
	require.Equal(t, id, "Session01", "The id of the votingInstance just created is incorrect, got: %s, want %s.", id, "Session01")

	status := VoteSystem.VotingInstancesList["Session01"].Status
	require.Equal(t, status, "open", "The status of the votingInstance just created is incorrect, got: %s, want %s.", status, "open")

	config := VoteSystem.VotingInstancesList["Session01"].Config
	require.Equal(t, config.Title, "TestVotingTitle", "The config title of the votingInstance just created is incorrect, got: %s, want %s.", config.Title, "TestVotingTitle")

	require.Equal(t, config.Description, "Quick description", "The config description of the votingInstance just created is incorrect, got: %s, want %s.", config.Description, "Quick description")

	_, err = NewVotingConfig(voters, "", "Quick description", candidats)
	require.Equal(t, err.Error(), "Title is empty")
}

func TestCreateAndAdd(t *testing.T) {
	voteConfig, err := NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	require.Equal(t, err, nil, "Creation of votingConfig is incorrect.")

	_, err = VoteSystem.CreateAndAdd("", voteConfig, "open")
	require.Equal(t, err.Error(), "The id is empty.")

	_, err = VoteSystem.CreateAndAdd("Session01", voteConfig, "")
	require.Equal(t, err.Error(), "The status is incorrect, should be either 'open' or 'close'.")

	vi, _ := VoteSystem.CreateAndAdd("Session01", voteConfig, "open")
	require.Equal(t, VoteSystem.VotingInstancesList["Session01"], vi, "Creation of the voting instance is incorrect")
}

func TestCloseVoting(t *testing.T) {
	vi.CloseVoting()
	require.Equal(t, vi.Status, "close", "Status incorrect. Was: %s, should be: %s", vi.Status, "close")
}

func TestSetStatus(t *testing.T) {
	err := vi.SetStatus("")
	require.Equal(t, err.Error(), "The status is incorrect, should be either 'open' or 'close'.")

	vi.SetStatus("close")
	require.Equal(t, vi.Status, "close", "Status incorrect. Was: %s, should be: %s", vi.Status, "close")

	vi.SetStatus("open")
	require.Equal(t, vi.Status, "open", "Status incorrect. Was: %s, should be: %s", vi.Status, "open")
}

func TestCreationOfLiquid(t *testing.T) {
	liquid_100, _ := NewLiquid(100)
	require.Equal(t, liquid_100.Percentage, 100.)

	_, err := NewLiquid(-10)
	require.Equal(t, err.Error(), "Init value is incorrect: was -10, must be positive.")
}

func TestNewChoice(t *testing.T) {
	tabChoice := map[string]voting.Liquid{}
	tabChoice["yes"] = liquid_100
	tabChoice["no"] = liquid_0

	//pas d'erreur throw, must have some
	choice, _ := NewChoice(tabChoice)
	require.Equal(t, choice.VoteValue["yes"].Percentage, 100., "Choice percentage of yes is incorrect : was %d, should be %d.", choice.VoteValue["yes"].Percentage, 100.)
	require.Equal(t, choice.VoteValue["no"].Percentage, 0., "Choice percentage of no is incorrect : was %d, should be %d.", choice.VoteValue["no"].Percentage, 0.)
}

func TestSetVote(t *testing.T) {
	yesChoice["yes"] = liquid_150
	yesChoice["no"] = liquid_0

	choiceGuillaume, _ := NewChoice(yesChoice)
	err := vi.SetVote(&userGuillaume, choiceGuillaume)
	require.Equal(t, err.Error(), "Voting power can't be negative.")

	yesChoice["yes"] = liquid_100
	choiceGuillaume2, _ := NewChoice(yesChoice)
	vi.SetVote(&userGuillaume, choiceGuillaume2)
	require.Equal(t, userGuillaume.HistoryOfChoice[len(userGuillaume.HistoryOfChoice)-1], choiceGuillaume2)
}

func TestGetResults(t *testing.T) {
	var histoChoice = make([]voting.Choice, 0)

	var userN, _ = VoteSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice)
	var userG, _ = VoteSystem.NewUser("Guillaume", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice)
	var userE, _ = VoteSystem.NewUser("Etienne", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice)

	var voters = []*voting.User{&userN, &userG, &userE}
	var candidats = make([]string, 0)

	var voteConfig, _ = NewVotingConfig(voters, "TestGetResults", "Quick description", candidats)

	var vi, _ = VoteSystem.CreateAndAdd("Session02", voteConfig, "open")

	var liquid_100, _ = NewLiquid(100)
	var liquid_50, _ = NewLiquid(50)
	var liquid_0, _ = NewLiquid(0)

	var yesChoice = make(map[string]voting.Liquid)
	var noChoice = make(map[string]voting.Liquid)
	var midChoice = make(map[string]voting.Liquid)

	yesChoice["yes"] = liquid_100
	yesChoice["no"] = liquid_0
	noChoice["no"] = liquid_100
	noChoice["yes"] = liquid_0
	midChoice["no"] = liquid_50
	midChoice["yes"] = liquid_50

	choiceG, _ := NewChoice(noChoice)
	choiceE, _ := NewChoice(midChoice)
	choiceN, _ := NewChoice(yesChoice)

	vi.SetVote(&userG, choiceG)
	vi.SetVote(&userE, choiceE)
	vi.SetVote(&userN, choiceN)

	propYes := vi.GetResults()["yes"]
	require.Equal(t, propYes, 50., "Yes proportion is incorrect, got: %f, want: %f.", propYes, 50.)

	propNo := vi.GetResults()["no"]
	require.Equal(t, propYes, 50., "No proportion is incorrect, got: %f, want: %f.", propNo, 50.)
}

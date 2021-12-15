package impl

import (
	"testing"

	"github.com/dedis/livos/voting"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

//VARIABLES NECESSAIRES
var VoteList = make(map[string]voting.VotingInstance)
var VoteSystem = NewVotingSystem(nil, VoteList)

//Creation of a empty list of choces (for history)
var histoChoice = make([]voting.Choice, 0)

var userNoemien, _ = VoteSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
var userGuillaume, _ = VoteSystem.NewUser("Guillaume", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
var userEtienne, _ = VoteSystem.NewUser("Etienne", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)

var voters = []*voting.User{&userNoemien, &userGuillaume, &userEtienne}
var candidats = make([]*voting.Candidate, 0)

var voteConfig, _ = NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats, "YesOrNoQuestion")

var vi, _ = VoteSystem.CreateAndAdd("Session01", voteConfig, "open")

var liquid_150, _ = NewLiquid(150)
var liquid_100, _ = NewLiquid(100)
var liquid_50, _ = NewLiquid(50)
var liquid_0, _ = NewLiquid(0)

var yesChoice = make(map[string]voting.Liquid)
var noChoice = make(map[string]voting.Liquid)
var midChoice = make(map[string]voting.Liquid)

func TestNewUser(t *testing.T) {
	userNoemien, err := VoteSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
	if err != nil {
		xerrors.Errorf(err.Error())
	}
	require.Equal(t, userNoemien.UserID, "Noemien", "UserID initialization incorrect : was %s, should be %s", userNoemien.UserID, "Noemien")
	require.Equal(t, userNoemien.VotingPower, 100., "VotingPower initialization incorrect : was %f, should be %f", userNoemien.VotingPower, 100.)
	require.Equal(t, userNoemien.HistoryOfChoice, make([]voting.Choice, 0), "HistoryOfChoice initialization incorrect")
	require.Equal(t, userNoemien.DelegatedFrom, make(map[string]voting.Liquid), "DelegatedFrom initialization incorrect")
	require.Equal(t, userNoemien.DelegatedTo, make(map[string]voting.Liquid), "DelegatedTo initialization incorrect")
}

func TestNewCandidate(t *testing.T) {
	candidateTrump, err := VoteSystem.NewCandidate("Trump")
	if err != nil {
		xerrors.Errorf(err.Error())
	}
	require.Equal(t, candidateTrump.CandidateID, "Trump")
}

func TestVotingInstanceCreate(t *testing.T) {
	voteConfig, err := NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats, "YesOrNoQuestion")
	require.Equal(t, err, nil, "Cannot create VotingConfig")

	VoteSystem.CreateAndAdd("Session01", voteConfig, "open")
	id := VoteSystem.VotingInstancesList["Session01"].GetVotingID()
	require.Equal(t, id, "Session01", "The id of the votingInstance just created is incorrect, got: %s, want %s.", id, "Session01")

	status := VoteSystem.VotingInstancesList["Session01"].GetStatus()
	require.Equal(t, status, "open", "The status of the votingInstance just created is incorrect, got: %s, want %s.", status, "open")

	config := VoteSystem.VotingInstancesList["Session01"].GetConfig()
	require.Equal(t, config.Title, "TestVotingTitle", "The config title of the votingInstance just created is incorrect, got: %s, want %s.", config.Title, "TestVotingTitle")

	require.Equal(t, config.Description, "Quick description", "The config description of the votingInstance just created is incorrect, got: %s, want %s.", config.Description, "Quick description")

}

func TestNewVotingConfig(t *testing.T) {
	voteConfig, _ := NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats, "YesOrNoQuestion")
	require.Equal(t, voteConfig.Voters, voters)
	require.Equal(t, voteConfig.Title, "TestVotingTitle")
	require.Equal(t, voteConfig.Description, "Quick description")
	require.Equal(t, voteConfig.Candidates, candidats)
	require.Equal(t, voteConfig.TypeOfVotingConfig, voting.TypeOfVotingConfig("YesOrNoQuestion"))

	_, err := NewVotingConfig(voters, "", "Quick description", candidats, "YesOrNoQuestion")
	require.Equal(t, err.Error(), "Title is empty")

	_, err = NewVotingConfig(voters, "Hello", "Quick description", candidats, "Blabla")
	require.Equal(t, err.Error(), "TypeOfVotingCOnfig Incorrect")

}

func TestCreateAndAdd(t *testing.T) {
	voteConfig, err := NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats, "YesOrNoQuestion")
	require.Equal(t, err, nil, "Creation of votingConfig is incorrect.")

	_, err = VoteSystem.CreateAndAdd("", voteConfig, "open")
	require.Equal(t, err.Error(), "The id is empty.")

	_, err = VoteSystem.CreateAndAdd("Session01", voteConfig, "")
	require.Equal(t, err.Error(), "The status is incorrect, should be either 'open' or 'close'.")

	vi, _ := VoteSystem.CreateAndAdd("Session01", voteConfig, "open")
	require.Equal(t, VoteSystem.VotingInstancesList["Session01"], vi, "Creation of the voting instance is incorrect")
}

/* func TestListVoting(t *testing.T) {
	vi1, _ := VoteSystem.CreateAndAdd("Session01", voteConfig, "open")
	vi2, _ := VoteSystem.CreateAndAdd("Session02", voteConfig, "open")
	vi3, _ := VoteSystem.CreateAndAdd("Session03", voteConfig, "open")
	verifListVotingSystem := []string{vi1.GetVotingID(), vi2.GetVotingID(), vi3.GetVotingID()}
	listeVotingSystems := VoteSystem.ListVotings()
	require.Equal(t, verifListVotingSystem, listeVotingSystems)
} */

func TestCloseVoting(t *testing.T) {
	vi.CloseVoting()
	require.Equal(t, vi.GetStatus(), "close", "Status incorrect. Was: %s, should be: %s", vi.GetStatus(), "close")
}

func TestSetStatus(t *testing.T) {
	err := vi.SetStatus("")
	require.Equal(t, err.Error(), "The status is incorrect, should be either 'open' or 'close'.")

	vi.SetStatus("close")
	require.Equal(t, vi.GetStatus(), "close", "Status incorrect. Was: %s, should be: %s", vi.GetStatus(), "close")

	vi.SetStatus("open")
	require.Equal(t, vi.GetStatus(), "open", "Status incorrect. Was: %s, should be: %s", vi.GetStatus(), "open")
}

func TestCreationOfLiquid(t *testing.T) {
	liquid_100, _ := NewLiquid(100)
	require.Equal(t, liquid_100.Percentage, 100.)

	_, err := NewLiquid(-10)
	require.Equal(t, err.Error(), "Init value is incorrect: was -10, must be positive.")
}

func TestAdditionOfLiquid(t *testing.T) {
	additionOfLiquids, _ := AddLiquid(liquid_100, liquid_150)
	require.Equal(t, additionOfLiquids.Percentage, 250.)
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

func TestDelegTo(t *testing.T) {
	vi.DelegTo(&userNoemien, &userEtienne, liquid_50)
	require.Equal(t, 50., userEtienne.DelegatedFrom["Noemien"].Percentage)
	require.Equal(t, 50., userNoemien.DelegatedTo["Etienne"].Percentage)
	require.Equal(t, 50., userNoemien.VotingPower, "userNoemien false power")
	require.Equal(t, 150., userEtienne.VotingPower, "userEtienne false power")
}

func TestGetResults(t *testing.T) {
	var histoChoice = make([]voting.Choice, 0)

	var userN, _ = VoteSystem.NewUser("Noemien", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
	var userG, _ = VoteSystem.NewUser("Guillaume", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)
	var userE, _ = VoteSystem.NewUser("Etienne", make(map[string]voting.Liquid), make(map[string]voting.Liquid), histoChoice, voting.None, nil)

	var candidatTrump, _ = VoteSystem.NewCandidate("Trump")
	var candidatObama, _ = VoteSystem.NewCandidate("Obama")
	var candidatJeanMi, _ = VoteSystem.NewCandidate("JeanMi")

	var voters = []*voting.User{&userN, &userG, &userE}
	var candidats = []*voting.Candidate{&candidatObama, &candidatTrump, &candidatJeanMi}

	var voteConfig, _ = NewVotingConfig(voters, "TestGetResults", "Quick description", candidats, "CandidateQuestion")

	var vi, _ = VoteSystem.CreateAndAdd("Session02", voteConfig, "open")

	if voteConfig.TypeOfVotingConfig == "YesOrNoQuestion" {
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
	} else if voteConfig.TypeOfVotingConfig == "CandidateQuestion" {
		var liquid_100, _ = NewLiquid(100)
		var liquid_70, _ = NewLiquid(70)
		var liquid_30, _ = NewLiquid(30)
		var liquid_0, _ = NewLiquid(0)

		var yesChoice = make(map[string]voting.Liquid)
		var noChoice = make(map[string]voting.Liquid)
		var midChoice = make(map[string]voting.Liquid)

		yesChoice["Trump"] = liquid_100
		yesChoice["Obama"] = liquid_0
		noChoice["Obama"] = liquid_100
		noChoice["Trump"] = liquid_0
		midChoice["Obama"] = liquid_70
		midChoice["Trump"] = liquid_30

		choiceG, _ := NewChoice(noChoice)
		choiceE, _ := NewChoice(midChoice)
		choiceN, _ := NewChoice(yesChoice)

		vi.SetVote(&userG, choiceG)
		vi.SetVote(&userE, choiceE)
		vi.SetVote(&userN, choiceN)

		propTrump := vi.GetResults()["Trump"]
		require.Equal(t, propTrump, 130., "Trump proportion is incorrect, got: %f, want: %f.", propTrump, 130.)

		propObama := vi.GetResults()["Obama"]
		require.Equal(t, propObama, 170., "Obama proportion is incorrect, got: %f, want: %f.", propObama, 170.)

		propJeanMi := vi.GetResults()["JeanMi"]
		require.Equal(t, propJeanMi, 0., "JeanMi proportion is incorrect, got: %f, want: %f.", propJeanMi, 0.)
	}
}

func TestGetUser(t *testing.T) {
	user, _ := vi.GetUser("Noemien")
	require.Equal(t, *user, userNoemien, "Get user returned incorrect user")
	require.Equal(t, user.UserID, "Noemien", "Get user returned incorrect userID")

	_, err := vi.GetUser("else")
	require.Equal(t, err.Error(), "Cannot find the user. UserId is incorrect.")
}

func TestGetCandidate(t *testing.T) {

	var candidatTrump, _ = VoteSystem.NewCandidate("Trump")
	var voters = make([]*voting.User, 0)
	var candidats = []*voting.Candidate{&candidatTrump}

	var voteConfig, _ = NewVotingConfig(voters, "TestGetResults", "Quick description", candidats, "CandidateQuestion")
	var vi, _ = VoteSystem.CreateAndAdd("Session02", voteConfig, "open")

	candidate, _ := vi.GetCandidate("Trump")
	require.Equal(t, candidate, &candidatTrump, "Get candidate returned incorrect user")
	require.Equal(t, candidate.CandidateID, "Trump", "Get candidate returned incorrect userID")

	_, err := vi.GetCandidate("else")
	require.Equal(t, err.Error(), "Cannot find the Candidate. CandidateID is incorrect.")
}

func TestDelete(t *testing.T) {
	var VoteList = make(map[string]voting.VotingInstance)
	var VoteSystem2 = NewVotingSystem(nil, VoteList)
	var voters = []*voting.User{}
	var candidats = make([]*voting.Candidate, 0)
	var voteConfig, _ = NewVotingConfig(voters, "TestDelete", "Quick description", candidats, "YesOrNoQuestion")

	VoteSystem2.CreateAndAdd("Session02", voteConfig, "open")
	require.Equal(t, 1, len(VoteSystem2.VotingInstancesList), "Creation of the voting instance not complete. It didn't appear on the voting instance list.")

	err := VoteSystem2.Delete("Session02")
	require.Equal(t, "Can't delete the votingInsance because it is still open", err.Error(), "Deletion of the voting instance incorrect. It's still there.")

	VoteSystem2.VotingInstancesList["Session02"].CloseVoting()
	VoteSystem2.Delete("Session02")
	require.Equal(t, 0, len(VoteSystem2.VotingInstancesList), "Deletion of the voting instance incorrect. It's still there.")
}

//delegFrom

package tests

import (
	"testing"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"github.com/stretchr/testify/require"
)

var VoteList = make(map[string]*impl.VotingInstance)
var VoteSystem = impl.NewVotingSystem(nil, VoteList)
var voters = make([]string, 3)
var candidats = make([]string, 3)
var votes = make(map[string]*voting.Choice)

func TestVotingSystemCreate(t *testing.T) {
	voters = append(voters, "Noemien", "Guillaume", "Etienne")
	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	if err != nil {
		t.Errorf("Cannot create VotingConfig")
	}

	VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)
	id := VoteSystem.VotingInstancesList["Session01"].Id
	if id != "Session01" {
		t.Errorf("The id of the votingInstance just created is incorrect, got: %s, want %s.", id, "Session01")
	}

	require.Equal(t, 3, 3)
	//require.Error(t, )

	status := VoteSystem.VotingInstancesList["Session01"].Status
	if status != "open" {
		t.Errorf("The status of the votingInstance just created is incorrect, got: %s, want %s.", status, "open")
	}

	config := VoteSystem.VotingInstancesList["Session01"].Config
	if config.Title != "TestVotingTitle" {
		t.Errorf("The config title of the votingInstance just created is incorrect, got: %s, want %s.", config.Title, "TestVotingTitle")
	}

	if config.Description != "Quick description" {
		t.Errorf("The config description of the votingInstance just created is incorrect, got: %s, want %s.", config.Description, "Quick description")
	}
}

func TestSetStatus(t *testing.T) {
	voters = append(voters, "Noemien", "Guillaume", "Etienne")
	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	require.Equal(t, err, nil, "Creation of votingConfig is incorrect.")

	VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)
	addVoteInst := VoteSystem.VotingInstancesList["Session01"]

	s := "close"
	addVoteInst.SetStatus(s)
	require.Equal(t, addVoteInst.Status, s, "Status incorrect. Was: %s, should be: %s", addVoteInst.Status, s)
}

func TestCloseVoting(t *testing.T) {
	voters = append(voters, "Noemien", "Guillaume", "Etienne")
	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	require.Equal(t, err, nil, "Creation of votingConfig is incorrect.")

	s := "close"
	VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)
	addVoteInst := VoteSystem.VotingInstancesList["Session01"]
	addVoteInst.CloseVoting()
	require.Equal(t, addVoteInst.Status, s, "Status incorrect. Was: %s, should be: %s", addVoteInst.Status, s)

}

func TestGetResults(t *testing.T) {
	voters = append(voters, "Noemien", "Guillaume", "Etienne")
	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	require.Equal(t, err, nil, "Creation of votingConfig is incorrect.")

	vi, err := VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)
	require.Equal(t, err, nil, "Creation of votingInstance is incorrect.")

	deleg := make(map[string]voting.Liquid)
	yesChoice := make(map[string]voting.Liquid)
	noChoice := make(map[string]voting.Liquid)
	midChoice := make(map[string]voting.Liquid)

	liq100, err := impl.NewLiquid(100)
	require.Equal(t, err, nil, "Creation of liquid is incorrect.")

	liq50, err := impl.NewLiquid(50)
	require.Equal(t, err, nil, "Creation of liquid is incorrect.")

	liqid0, err := impl.NewLiquid(0)
	require.Equal(t, err, nil, "Creation of liquid is incorrect.")

	yesChoice["yes"] = liq100
	yesChoice["no"] = liqid0
	noChoice["no"] = liq100
	noChoice["yes"] = liqid0
	midChoice["no"] = liq50
	midChoice["yes"] = liq50
	choiceNoemien, err := impl.NewChoice(deleg, yesChoice, 0, 100)
	require.Equal(t, err, nil, "Creation of the choice is incorrect.")

	choiceGuillaume, err := impl.NewChoice(deleg, noChoice, 0, 100)
	require.Equal(t, err, nil, "Creation of the choice is incorrect.")

	choiceEtienne, err := impl.NewChoice(deleg, midChoice, 0, 100)
	require.Equal(t, err, nil, "Creation of the choice is incorrect.")

	err = vi.CastVote("Noemien", choiceNoemien)
	require.Equal(t, err, nil, "Impossible to cast a vote on a closed session.")

	err = vi.CastVote("Guillaume", choiceGuillaume)
	require.Equal(t, err, nil, "Impossible to cast a vote on a closed session.")

	err = vi.CastVote("Etienne", choiceEtienne)
	require.Equal(t, err, nil, "Impossible to cast a vote on a closed session.")

	propYes := vi.GetResults()["yes"]
	require.Equal(t, propYes, 50., "Yes proportion is incorrect, got: %f, want: %f.", propYes, 50.)

	propNo := vi.GetResults()["no"]
	require.Equal(t, propYes, 50., "No proportion is incorrect, got: %f, want: %f.", propNo, 50.)

}

func TestCastVotes(t *testing.T) {
	voters = append(voters, "Noemien", "Guillaume", "Etienne")
	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	require.Equal(t, err, nil, "Creation of votingConfig is incorrect.")

	vi, err := VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)
	require.Equal(t, err, nil, "Creation of votingInstance is incorrect.")

	deleg := make(map[string]voting.Liquid)
	yesChoice := make(map[string]voting.Liquid)

	liq100, err := impl.NewLiquid(100)
	require.Equal(t, err, nil, "Creation of liquid is incorrect.")

	liqid0, err := impl.NewLiquid(0)
	require.Equal(t, err, nil, "Creation of liquid is incorrect.")

	yesChoice["yes"] = liq100
	yesChoice["no"] = liqid0

	choiceNoemien, err := impl.NewChoice(deleg, yesChoice, 0, 100)
	require.Equal(t, err, nil, "Creation of the choice is incorrect.")

	err = vi.CastVote("Noemien", choiceNoemien)
	require.Equal(t, err, nil, "Impossible to cast a vote on a closed session.")

	require.Equal(t, vi.Votes["Noemien"].MyChoice["yes"].Percentage, 100., "Proportion in yes is incorrect. Was: %f, should be %f", vi.Votes["Noemien"].MyChoice["yes"].Percentage, 100.)
}

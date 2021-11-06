package tests

import (
	"testing"

	"github.com/dedis/livos/voting"
	"github.com/dedis/livos/voting/impl"
	"github.com/stretchr/testify/require"
	//"github.com/dedis/livos/storage/bbolt"
	//"github.com/dedis/livos/storage.DB"
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

// func TestSetStatus(t *testing.T) {
// 	voters = append(voters, "Noemien", "Guillaume", "Etienne")
// 	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
// 	if err != nil {
// 		t.Errorf("Cannot create VotingConfig")
// 	}
// 	VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)

// 	addVoteInst := VoteSystem.VotingInstancesList["Session01"]

// 	if addVoteInst.Status != "close" {
// 		t.Errorf("Set status was incorrect, got: %s, want %s", addVoteInst.Status, "close")
// 	}

// }

// func TestCloseVoting(t *testing.T) {
// 	voters = append(voters, "Noemien", "Guillaume", "Etienne")
// 	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
// 	if err != nil {
// 		t.Errorf("Cannot create VotingConfig")
// 	}
// 	VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)
// 	VoteSystem.VotingInstancesList["Session01"].SetStatus("close")
// 	status := VoteSystem.VotingInstancesList["Session01"].Status
// 	if status != "close" {
// 		t.Errorf("CloseVoting was incorrect, got: %s, want %s", status, "close")
// 	}
// }

func TestGetResults(t *testing.T) {
	voters = append(voters, "Noemien", "Guillaume", "Etienne")
	voteConfig, err := impl.NewVotingConfig(voters, "TestVotingTitle", "Quick description", candidats)
	if err != nil {
		t.Errorf("Cannot create VotingConfig")
	}
	VoteSystem.CreateAndAdd("Session01", voteConfig, "open", votes)
	vi := VoteSystem.VotingInstancesList["Session01"]

	deleg := make(map[string]voting.Liquid)
	yesChoice := make(map[string]voting.Liquid)
	noChoice := make(map[string]voting.Liquid)
	midChoice := make(map[string]voting.Liquid)

	liq100, err100 := impl.NewLiquid(100)
	liq50, err50 := impl.NewLiquid(50)
	liqid0, err0 := impl.NewLiquid(0)
	if (err100 != nil) || (err50 != nil) || (err0 != nil) {
		t.Error("Creation of liquid is incorrect.")
	}

	yesChoice["yes"] = liq100
	yesChoice["no"] = liqid0
	noChoice["no"] = liq100
	noChoice["yes"] = liqid0
	midChoice["no"] = liq50
	midChoice["yes"] = liq50
	choiceNoemien, errN := impl.NewChoice(deleg, yesChoice, 0, 0)
	choiceGuillaume, errG := impl.NewChoice(deleg, noChoice, 0, 0)
	choiceEtienne, errE := impl.NewChoice(deleg, midChoice, 0, 0)
	if (errN != nil) || (errG != nil) || (errE != nil) {
		t.Error("Choices creation not correct.")
	}
	vi.CastVote("Noemien", choiceNoemien)
	vi.CastVote("Guillaume", choiceGuillaume)
	vi.CastVote("Etienne", choiceEtienne)
	propYes := vi.GetResults()["yes"]
	if propYes != 50. {
		t.Errorf("Yes proportion is incorrect, got: %f, want: %f.", propYes, 50.)
	}
	propNo := vi.GetResults()["no"]
	if propNo != 50. {
		t.Errorf("No proportion is incorrect, got: %f, want: %f.", propNo, 50.)
	}
}

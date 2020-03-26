package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {

	rand := NewRNG(MaxNum)
	game := NewGame(rand)

	if game == nil {
		t.Errorf("Game New() got nil, want not nil")
	}

}

func TestValidateNumberChoice(t *testing.T) {
	assert := assert.New(t)

	rand := NewRNG(MaxNum)
	game := NewGame(rand)

	// Check min
	err := game.validateChoice(MinNum - 1)
	assert.Equal(err, ErrInvalidNumber)

	// Check max
	err = game.validateChoice(MaxNum + 1)
	assert.Equal(err, ErrInvalidNumber)

	// Check valid
	err = game.validateChoice(7)
	assert.Nil(err)
}

func TestAddingPlayer(t *testing.T) {
	assert := assert.New(t)

	rand := NewRNG(MaxNum)
	game := NewGame(rand)

	// New player
	err := game.RegisterPlayer(&Player{"Steve", 5, 8})
	assert.Nil(err)

	// Same player
	err = game.RegisterPlayer(&Player{"Steve", 5, 8})
	assert.Equal(ErrInvalidPlayerName, err)

	// New player but bad first number
	err = game.RegisterPlayer(&Player{"Sarah", 15, 8})
	assert.Equal(err, ErrInvalidNumber)

	// New player but bad second number
	err = game.RegisterPlayer(&Player{"Sarah", 5, 84})
	assert.Equal(err, ErrInvalidNumber)

}

func TestRounds(t *testing.T) {
	assert := assert.New(t)

	rand := NewRNG(MaxNum)
	game := NewGame(rand)
	game.RegisterPlayer(&Player{"Steve", 3, 9})
	game.RegisterPlayer(&Player{"Sarah", 4, 2})

	for i := 0; i < MaxRounds; i++ {
		err := game.PlayRound()
		if err != nil {
			t.Errorf("Error playing round. %s" + err.Error())
		}
		assert.Equal(i+1, game.Round)
	}
	assert.Equal(game.GetState(), GameStateCompleted)
}

func TestVectorExampleGame(t *testing.T) {
	assert := assert.New(t)
	seq := []int{9, 1, 4, 10, 7, 5, 3}
	gen := NewSSNG(seq)

	game := NewGame(gen)
	// Example with 3 players
	game.RegisterPlayer(&Player{"PlayerA", 3, 8})
	game.RegisterPlayer(&Player{"PlayerB", 5, 7})
	game.RegisterPlayer(&Player{"PlayerC", 3, 7})
	game.AddWaitingPlayersToGame()

	// Round 1: -1 -1 -1
	// Round 2: -2 -2 -2
	// Round 3: -2 -3 -1
	// Round 4: -3 -4 -2
	// Round 5: -3  1  3
	// Round 6: -3  6  4
	// Round 7:  2  5  9

	game.PlayRound()
	assert.Equal(game.Players["PlayerA"].Score, -1)
	assert.Equal(game.Players["PlayerB"].Score, -1)
	assert.Equal(game.Players["PlayerC"].Score, -1)
	game.PlayRound()
	assert.Equal(game.Players["PlayerA"].Score, -2)
	assert.Equal(game.Players["PlayerB"].Score, -2)
	assert.Equal(game.Players["PlayerC"].Score, -2)
	game.PlayRound()
	assert.Equal(game.Players["PlayerA"].Score, -2)
	assert.Equal(game.Players["PlayerB"].Score, -3)
	assert.Equal(game.Players["PlayerC"].Score, -1)
	game.PlayRound()
	assert.Equal(game.Players["PlayerA"].Score, -3)
	assert.Equal(game.Players["PlayerB"].Score, -4)
	assert.Equal(game.Players["PlayerC"].Score, -2)
	game.PlayRound()
	assert.Equal(game.Players["PlayerA"].Score, -3)
	assert.Equal(game.Players["PlayerB"].Score, 1)
	assert.Equal(game.Players["PlayerC"].Score, 3)
	game.PlayRound()
	assert.Equal(game.Players["PlayerA"].Score, -3)
	assert.Equal(game.Players["PlayerB"].Score, 6)
	assert.Equal(game.Players["PlayerC"].Score, 4)
	game.PlayRound()
	assert.Equal(game.Players["PlayerA"].Score, 2)
	assert.Equal(game.Players["PlayerB"].Score, 5)
	assert.Equal(game.Players["PlayerC"].Score, 9)

	// Check PlayerC wins here
	winner, _ := game.NominateWinner()
	assert.Equal(winner.Name, "PlayerC")
}

func TestNominatingAlphabeticalWinner(t *testing.T) {
	assert := assert.New(t)
	seq := []int{9, 1, 4, 10, 7, 5, 3}
	gen := NewSSNG(seq)

	game := NewGame(gen)
	// Example with 3 players
	game.RegisterPlayer(&Player{"ZZZ", 3, 8})
	game.RegisterPlayer(&Player{"BBB", 3, 8})
	game.RegisterPlayer(&Player{"NNN", 3, 8})
	game.AddWaitingPlayersToGame()

	game.PlayRound()
	// Should have the same scores and bounds
	game.NominateWinner()
	winner, _ := game.NominateWinner()
	assert.Equal(winner.Name, "BBB")
}

func TestNominatingUpperBoundWinner(t *testing.T) {
	assert := assert.New(t)
	seq := []int{9, 1, 4, 10, 7, 5, 3}
	gen := NewSSNG(seq)

	game := NewGame(gen)
	// Example with 3 players, highest upper bound playerA should win here with 8
	game.RegisterPlayer(&Player{"PlayerA", 3, 8})
	game.RegisterPlayer(&Player{"PlayerB", 5, 7})
	game.RegisterPlayer(&Player{"PlayerC", 3, 7})
	game.AddWaitingPlayersToGame()

	game.PlayRound()
	// PlayerB should win on highest upper bound
	winner, _ := game.NominateWinner()
	assert.Equal(winner.Name, "PlayerA")
}

func TestNominatingLowerBoundWinner(t *testing.T) {
	assert := assert.New(t)
	seq := []int{9, 1, 4, 10, 7, 5, 3}
	gen := NewSSNG(seq)

	game := NewGame(gen)
	// Example with 3 players, highest lower bound playerB should win here with 5
	game.RegisterPlayer(&Player{"PlayerA", 3, 7})
	game.RegisterPlayer(&Player{"PlayerB", 5, 7})
	game.RegisterPlayer(&Player{"PlayerC", 3, 7})
	game.AddWaitingPlayersToGame()

	game.PlayRound()
	// PlayerB should win on highest upper bound
	winner, _ := game.NominateWinner()
	assert.Equal(winner.Name, "PlayerB")
}

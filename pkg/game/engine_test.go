package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ()

func setup(t *testing.T) {

}

func TestEngineStartStop(t *testing.T) {
	assert := assert.New(t)
	// TODO: Move to setup test func
	mockGame := NewMockGame(10)
	engineConfig := &EngineConfig{
		GameSpeed:    1 * time.Minute,
		WaitingCount: 10,
		ManualRun:    true,
	}
	engine := NewEngine(mockGame, engineConfig)
	engine.Start()

	// Tests
	assert.Equal(engine.IsRunning(), true)
	wait := make(chan bool)
	engine.Cancel <- wait
	<-wait
	assert.Equal(engine.IsRunning(), false)

}

func TestEngineStartGamePlayRoundsThenComplete(t *testing.T) {
	assert := assert.New(t)
	// TODO: Move to setup test func
	mockGameMaxRounds := 10
	mockGame := NewMockGame(mockGameMaxRounds)
	engineConfig := &EngineConfig{
		GameSpeed:    10 * time.Nanosecond,
		WaitingCount: 10,
		ManualRun:    true,
	}
	engine := NewEngine(mockGame, engineConfig)
	// Set the manual ticker
	manualTicker := NewManualTicker()
	engine.Ticker = manualTicker.GetTicker()
	engine.Start()

	// Tests
	assert.Equal(engine.IsRunning(), true)
	// Put the game into ready mode, engine should start playing rounds
	mockGame.State = GameStateInProgress
	for i := 0; i < mockGameMaxRounds; i++ {
		manualTicker.Tick()
		<-engine.Event
	}
	assert.Equal(mockGameMaxRounds, mockGame.PlayRoundCalled)
	assert.Equal(GameStateCompleted, mockGame.GetState())

}

func TestEngineAddingPlayerToGame(t *testing.T) {
	assert := assert.New(t)
	// TODO: Move to setup test func
	mockGameMaxRounds := 10
	mockGame := NewMockGame(mockGameMaxRounds)
	engineConfig := &EngineConfig{
		GameSpeed:    10 * time.Nanosecond,
		WaitingCount: 10,
		ManualRun:    true,
	}
	engine := NewEngine(mockGame, engineConfig)
	engine.Start()

	// Tests
	player := &Player{"Steve", 5, 3}
	rc := make(chan *ActionResponse)
	join := &Action{ActionTypeJoinGame, player, rc}
	engine.Action <- join
	resp := <-rc
	assert.True(resp.Success)
}

func TestEngineCountdown(t *testing.T) {
	assert := assert.New(t)
	// TODO: Move to setup test func
	rand := NewRNG(MaxNum)
	game := NewGame(rand)
	engineConfig := &EngineConfig{
		GameSpeed:    10 * time.Minute,
		WaitingCount: 10,
		ManualRun:    true,
	}
	engine := NewEngine(game, engineConfig)
	// Set the manual ticker
	manualTicker := NewManualTicker()
	engine.Ticker = manualTicker.GetTicker()
	engine.Start()

	// Game should be waiting for minimum players to join
	assert.Equal(GameStateWaiting, game.GetState())
	// Tests
	playerOne := &Player{"Steve", 5, 3}
	playerTwo := &Player{"Sarah", 4, 1}
	rc := make(chan *ActionResponse)
	join := &Action{ActionTypeJoinGame, playerOne, rc}
	engine.Action <- join
	<-rc
	<-engine.Event
	// Engine should be waiting for one more player
	assert.Equal(GameStateWaiting, game.GetState())
	join = &Action{ActionTypeJoinGame, playerTwo, rc}
	engine.Action <- join
	<-rc
	<-engine.Event
	// Turn the engine manually once
	manualTicker.Tick()
	// Add waiting players to game
	<-engine.Event
	<-engine.Event

	// Game should be in ready state as it has enough players
	assert.Equal(GameStateReady, game.GetState())
	// Start countdown
	manualTicker.Tick()
	event := <-engine.Event
	assert.Equal(CountdownStarted.String(), event.Type)
	// continue the countdown 10, 9, 8,...,1
	for i := 1; i < engine.Config.WaitingCount; i++ {
		manualTicker.Tick()
		event = <-engine.Event
		assert.Equal(CountingDown.String(), event.Type)
	}
	// Next turn should start the game
	manualTicker.Tick()
	event = <-engine.Event
	assert.Equal(GameStarted.String(), event.Type)
	assert.Equal(GameStateInProgress, game.GetState())
}

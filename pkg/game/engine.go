package game

import (
	"log"
	"time"
)

// ActionType - wrapper around int
type ActionType int

// Player - Not in the game
type Player struct {
	Name   string
	First  int
	Second int
}

// Action - external actions that may affect the state of the game
type Action struct {
	Type   ActionType
	Player *Player
	Reply  chan *ActionResponse
}

// ActionTypes
const (
	ActionTypeJoinGame    ActionType = 0
	ActionTypeObserveGame ActionType = 1
)

// ActionResponse - Result of action returned to original caller
type ActionResponse struct {
	Success bool
	Message string
}

// Engine - Runs the game and mutates game state
type Engine struct {
	Event        chan *Event
	Action       chan *Action
	Cancel       chan chan bool
	Game         GameI
	Config       *EngineConfig
	running      bool
	count        int
	countingDown bool
	Ticker       <-chan time.Time
}

// EngineConfig - Parameters that alter the behavior of the game
type EngineConfig struct {
	WaitingCount int
	GameSpeed    time.Duration
	ManualRun    bool
}

// NewEngine - Initiates a new Engine with given Game and Config
func NewEngine(game GameI, config *EngineConfig) *Engine {

	event := make(chan *Event)
	action := make(chan *Action)
	cancel := make(chan chan bool)

	return &Engine{
		Event:  event,
		Cancel: cancel,
		Action: action,
		Game:   game,
		Config: config,
	}
}

// Start - The engine, enters the game loop
// Handles time ticks, external actions, game mutations, and broadcasting events
func (eng *Engine) Start() {

	// log.Println("GameEngine - Starting...")
	eng.running = true
	go func() {

		// If we are not manual set a timed ticker
		if !eng.Config.ManualRun {
			// TODO: Swap this for NewTicker as this is a memory leak if not cancelled
			eng.Ticker = time.Tick(eng.Config.GameSpeed)
		}

		// log.Println("GameEngine - Entering game loop...")
		for {
			// log.Println("GameEngine - In")
			select {
			case <-eng.Ticker:
				gameState := eng.Game.GetState()
				switch gameState {

				case GameStateWaiting:
					// log.Println("GameEngine - Waiting")
					// Add new players to the game
					joined, _ := eng.Game.AddWaitingPlayersToGame()
					if len(joined) > 0 {
						eng.Event <- NewEvent(PlayerJoined, joined)
						// log.Printf("Waiting - Added players to the game [+%v]\n", joined)
					}
					eng.Game.GetReady()
					eng.Event <- NewEvent(GameWaiting, nil)
				case GameStateReady:
					// log.Println("GameEngine - Ready Countdown")
					// Add new players on countdown:
					joined, _ := eng.Game.AddWaitingPlayersToGame()
					if len(joined) > 0 {
						// log.Printf("Ready - Added players to the game [+%v]\n", joined)
						eng.Event <- NewEvent(PlayerJoined, joined)
					}
					if !eng.countingDown {
						// Start counting down if we haven't already
						eng.startCountdown()
						eng.Event <- NewEvent(CountdownStarted, eng.count)
					} else if eng.isCountdownComplete() {
						// Check if countdown is complete
						eng.Game.Start()
						eng.Event <- NewEvent(GameStarted, eng.count)
					} else {
						// Otherwise, keep counting down
						eng.countdown()
						eng.Event <- NewEvent(CountingDown, eng.count)
					}

				case GameStateInProgress:
					// log.Println("GameEngine - Playing Round")
					err := eng.Game.PlayRound()
					if err != nil {
						// If something bad's happened here, we want to know about it!
						log.Fatal(err.Error())
						return
					}
					eng.Event <- NewEvent(PlayedRound, eng.Game.GetRoundResult())
					// log.Printf("Played Round: %+v\n", eng.Game)

				case GameStateCompleted:
					// log.Println("GameEngine - Completed")
					winner, err := eng.Game.NominateWinner()
					if err != nil {
						// If something bad's happened here, we want to know about it!
						log.Fatal(err.Error())
						return
					}
					eng.Event <- NewEvent(GameCompleted, winner)
					eng.Game.Reset()
					eng.resetCountdown()
					eng.Event <- NewEvent(GameReset, eng.Game)

				case GameStateCancelled:
					// log.Println("GameEngine - Cancelled")
					// TODO: Implement
					return
				}

			case action := <-eng.Action:

				switch action.Type {

				case ActionTypeJoinGame:
					// log.Println("Received Join Game Action")
					err := eng.Game.RegisterPlayer(action.Player)
					if err != nil {
						action.Reply <- &ActionResponse{false, err.Error()}
						log.Printf("Unable to add player: %s\n", err.Error())
					} else {
						action.Reply <- &ActionResponse{true, ""}
						eng.Event <- NewEvent(PlayerRegistered, eng.Game)
					}
				}

			case wait := <-eng.Cancel:
				eng.Game.Cancel()
				eng.running = false
				wait <- true
				return
			}
			// log.Println("GameEngine - Out")
		}
	}()
}

func (eng *Engine) IsRunning() bool {
	return eng.running
}

func (eng *Engine) startCountdown() {
	eng.resetCountdown()
	eng.countingDown = true
}

func (eng *Engine) resetCountdown() {
	eng.count = eng.Config.WaitingCount
	eng.countingDown = false
}

func (eng *Engine) countdown() {
	eng.count--
}

func (eng *Engine) isCountdownComplete() bool {
	if !eng.countingDown {
		return false
	} else {
		if eng.count <= 1 {
			eng.cancelCountdown()
			return true
		} else {
			return false
		}
	}
}

func (eng *Engine) cancelCountdown() {
	eng.countingDown = false
}

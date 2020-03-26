package game

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type State int
type EventType int

const (

	// Levers
	Name               = "Small Browser Based Game"
	MaxRounds          = 30
	MinPlayersRequired = 2
	MinNum             = 1
	MaxNum             = 10
	BlackJack          = 21

	// Scores
	ExactMatchScore   = 5
	InsideBoundsScore = 5
	OutOfBoundsScore  = -1

	// Game States
	GameStateReady      State = 1
	GameStateInProgress State = 2
	GameStateCompleted  State = 3
	GameStateWaiting    State = 4
	GameStateCancelled  State = 5
)

var (
	ErrInvalidNumber     = fmt.Errorf("Invalid number: Choose a number between %d - %d", MinNum, MaxNum)
	ErrInvalidPlayerName = fmt.Errorf("Invalid name: There is already a player here with that name")
	ErrNotEnoughPlayers  = errors.New("Invalid action: Not enough players in the game")
	ErrGameInProgress    = errors.New("Invalid action: Game is in progress")
	ErrGameComplete      = errors.New("Invalid action: Game in complete")
	ErrNoSingleWinner    = errors.New("Invalid state: Not single winner nominated")
)

type GameI interface {
	GetReady() error
	Start() error
	PlayRound() error
	UpdatePlayerScores(number int)
	NominateWinner() (GamePlayer, error)
	Reset() error
	RegisterPlayer(player *Player) error
	CheckPlayerExists(name string) error
	GetState() State
	Cancel() error
	AddWaitingPlayersToGame() ([]*GamePlayer, error)
	GetRoundResult() RoundResult
}

type GamePlayer struct {
	Name   string `json:"name"`
	Upper  int    `json:"upper"`
	Lower  int    `json:"lower"`
	Score  int    `json:"score"`
	Winner bool   `json:"winner"`
}

type Game struct {
	Players     map[string]GamePlayer `json:"players"`
	Round       int                   `json:"round"`
	Numbers     [MaxRounds]int        `json:"numbers"`
	Rand        NumberGenerator       `json:"-"`
	TopScore    int                   `json:"top_score"`
	Winner      GamePlayer            `json:"winner"`
	state       State
	registered  map[string]GamePlayer
	waitingRoom []*GamePlayer
}

// RoundResult - Sorted leader board for API
type RoundResult struct {
	LeaderBoard []GamePlayer `json:"leader_board"`
	Round       int          `json:"round"`
}

func NewGame(rand NumberGenerator) *Game {
	players := make(map[string]GamePlayer, 0)
	registered := make(map[string]GamePlayer, 0)
	waitingRoom := make([]*GamePlayer, 0)

	return &Game{
		Players:     players,
		Round:       0,
		Numbers:     [MaxRounds]int{},
		Rand:        rand,
		TopScore:    math.MinInt8,
		Winner:      GamePlayer{},
		state:       GameStateWaiting,
		registered:  registered,
		waitingRoom: waitingRoom,
	}
}

func (g *Game) GetReady() error {

	if g.state == GameStateInProgress {
		return ErrGameInProgress
	}

	if len(g.Players) >= MinPlayersRequired {
		g.state = GameStateReady
	}

	return nil
}

func (g *Game) Start() error {

	if g.state == GameStateInProgress {
		return ErrGameInProgress
	}

	if len(g.Players) < MinPlayersRequired {
		return ErrNotEnoughPlayers
	}

	g.state = GameStateInProgress
	g.Round = 0

	return nil
}

func (g *Game) PlayRound() error {

	if g.Round >= MaxRounds {
		return ErrGameComplete
	}

	roundNumber := g.Rand.GetInt()
	// Update Scores and Leader Board
	g.UpdatePlayerScores(roundNumber)

	g.Numbers[g.Round] = roundNumber
	g.Round++

	if g.Round >= MaxRounds {
		g.state = GameStateCompleted
	}

	return nil
}

func (g *Game) UpdatePlayerScores(number int) {
	// set top to a low value
	topScore := math.MinInt8
	// Loop through players in game
	for name, player := range g.Players {
		if player.Upper == number || player.Lower == number {
			player.Score += ExactMatchScore
			// fmt.Printf("Player %s ExactMatch. New Score[%d]\n", name, player.Score)
		} else if number <= player.Upper && number >= player.Lower {
			player.Score += InsideBoundsScore - (player.Upper - player.Lower)
			// fmt.Printf("Player %s InsideBoundsScore. New Score[%d]\n", name, player.Score)
		} else {
			player.Score += OutOfBoundsScore
			// fmt.Printf("Player %s OutOfBoundsScore. New Score[%d]\n", name, player.Score)
		}
		// Update topScore
		if player.Score > topScore {
			topScore = player.Score
		}
		// Check for rogue win case
		if player.Score == BlackJack {
			player.Winner = true
			g.Winner = player
			g.state = GameStateCompleted
		}
		// Write updates back to the map!
		g.Players[name] = player
		g.TopScore = topScore
	}
}

// TODO: This fails if players draw and their scores are negative! Fix it!
func (g *Game) NominateWinner() (GamePlayer, error) {
	// Check for BlackJack winner
	if g.Winner.Name != "" {
		return g.Winner, nil
	}

	winner := GamePlayer{}
	winners := []GamePlayer{}
	bigUp := 0
	bigLow := 0
	// Find players with top score and capture highest upper/lower bounds
	for _, player := range g.Players {
		if player.Score == g.TopScore {
			winners = append(winners, player)
			if player.Upper > bigUp {
				bigUp = player.Upper
			}
			if player.Lower > bigLow {
				bigLow = player.Lower
			}
		}
	}
	// Check for one winner
	if len(winners) == 1 {
		winners[0].Winner = true
		g.Players[winners[0].Name] = winners[0]
		return winners[0], nil
	}
	// Keep winners with highest upper bound
	bigUpWinners := []GamePlayer{}
	for _, player := range winners {
		if player.Upper == bigUp {
			bigUpWinners = append(bigUpWinners, player)
		}
	}
	// Check is there one winner
	if len(bigUpWinners) == 1 {
		bigUpWinners[0].Winner = true
		g.Players[bigUpWinners[0].Name] = bigUpWinners[0]
		return bigUpWinners[0], nil
	}

	// Keep winners with highest lower bound
	// Prep name array in case more than one
	names := []string{}
	bigLowWinners := []GamePlayer{}
	for _, player := range bigUpWinners {
		if player.Lower == bigLow {
			bigLowWinners = append(bigLowWinners, player)
			names = append(names, player.Name)
		}
	}
	// Check is there one winner
	if len(bigLowWinners) == 1 {
		bigLowWinners[0].Winner = true
		g.Players[bigLowWinners[0].Name] = bigLowWinners[0]
		return bigLowWinners[0], nil
	}

	// Return the player who's first in alphabetical order
	sort.Strings(names)
	for _, player := range bigLowWinners {
		if player.Name == names[0] {
			player.Winner = true
			g.Players[names[0]] = player
			return player, nil
		}
	}

	return winner, ErrNoSingleWinner
}

func (g *Game) Reset() error {
	g.Round = 0
	g.state = GameStateWaiting
	g.Winner = GamePlayer{}
	for k, player := range g.Players {
		player.Score = 0
		g.Players[k] = player
	}
	// Get players from the waiting room
	return nil
}

func (g *Game) RegisterPlayer(player *Player) error {

	// Check name not in use
	if err := g.CheckPlayerExists(player.Name); err != nil {
		return err
	}
	// Check choices
	if err := g.validateChoice(player.First); err != nil {
		return err
	}
	if err := g.validateChoice(player.Second); err != nil {
		return err
	}

	gp := GamePlayer{Name: player.Name}
	if player.First >= player.Second {
		gp.Upper = player.First
		gp.Lower = player.Second
	} else {
		gp.Upper = player.Second
		gp.Lower = player.First
	}
	g.registered[player.Name] = gp
	g.waitingRoom = append(g.waitingRoom, &gp)

	return nil
}

func (g *Game) AddWaitingPlayersToGame() ([]*GamePlayer, error) {

	emptyWaitingRoom := make([]*GamePlayer, 0)

	if g.state == GameStateInProgress {
		return emptyWaitingRoom, ErrGameInProgress
	}

	waiting := g.waitingRoom
	g.waitingRoom = emptyWaitingRoom

	// TODO: Check if the name already exists here and return an error
	for _, waitingPlayer := range waiting {
		g.Players[waitingPlayer.Name] = *waitingPlayer
	}

	return waiting, nil
}

func (g *Game) GetRoundResult() RoundResult {
	// Make a slice
	leaderBoard := make([]GamePlayer, 0)
	for _, player := range g.Players {
		leaderBoard = append(leaderBoard, player)
	}
	// Sort the slice
	sort.Slice(leaderBoard, func(i, j int) bool {
		return leaderBoard[i].Score > leaderBoard[j].Score
	})

	return RoundResult{
		Round:       g.Round,
		LeaderBoard: leaderBoard,
	}
}

func (g *Game) GetState() State {
	return g.state
}

func (g *Game) CheckPlayerExists(name string) error {
	_, exists := g.registered[name]
	if exists {
		return ErrInvalidPlayerName
	}

	return nil
}

func (g *Game) Cancel() error {
	g.state = GameStateCancelled

	return nil
}

func (g *Game) validateChoice(i int) error {
	if i < MinNum || i > MaxNum {
		return ErrInvalidNumber
	}

	return nil
}

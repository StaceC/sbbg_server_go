package game

type MockGame struct {
	Players                  map[string]GamePlayer `json:"players"`
	Round                    int                   `json:"round"`
	Numbers                  [MaxRounds]int        `json:"numbers"`
	Rand                     NumberGenerator       `json:"-"`
	TopScore                 int                   `json:"top_score"`
	State                    State                 `json:"state"`
	StartCalled              int
	PlayRoundCalled          int
	UpdatePlayerScoredCalled int
	maxRounds                int
}

func NewMockGame(maxRounds int) *MockGame {
	return &MockGame{
		State:     GameStateWaiting,
		maxRounds: maxRounds,
	}
}

func (gm *MockGame) GetReady() error {
	return nil
}

func (gm *MockGame) Start() error {
	gm.StartCalled++
	return nil
}

func (gm *MockGame) PlayRound() error {
	gm.Round++
	gm.PlayRoundCalled++

	if gm.Round >= gm.maxRounds {
		gm.State = GameStateCompleted
	}
	return nil
}

func (gm *MockGame) UpdatePlayerScores(number int) {
	gm.UpdatePlayerScoredCalled++
}

func (gm *MockGame) NominateWinner() (GamePlayer, error) {
	return GamePlayer{}, nil
}

func (gm *MockGame) Reset() error {
	return nil
}

func (gm *MockGame) RegisterPlayer(player *Player) error {
	return nil
}

func (gm *MockGame) CheckPlayerExists(name string) error {
	return nil
}

func (gm *MockGame) GetState() State {
	return gm.State
}

func (gm *MockGame) Cancel() error {
	return nil
}

func (gm *MockGame) AddWaitingPlayersToGame() ([]*GamePlayer, error) {
	return make([]*GamePlayer, 0), nil
}

func (gm MockGame) GetRoundResult() RoundResult {
	return RoundResult{}
}

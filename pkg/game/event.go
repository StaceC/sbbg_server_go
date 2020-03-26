package game

const (
	// Events
	PlayerJoined     EventType = 0
	PlayerLeft       EventType = 1
	PlayedRound      EventType = 2
	GameCreated      EventType = 3
	GameStarted      EventType = 4
	GameCompleted    EventType = 5
	GameReady        EventType = 6
	GameWaiting      EventType = 7
	CountdownStarted EventType = 8
	CountingDown     EventType = 9
	GameReset        EventType = 10
	PlayerRegistered EventType = 11
)

func (et EventType) String() string {
	names := [...]string{
		"Player Joined",
		"Player Left",
		"Played Round",
		"Game Created",
		"Game Started",
		"Game Completed",
		"Game Ready",
		"Game Waiting",
		"Countdown Started",
		"Counting Down",
		"Game Reset",
		"Player Registered",
	}

	return names[et]
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewEvent(eventType EventType, data interface{}) *Event {
	return &Event{eventType.String(), data}
}

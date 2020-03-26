package game

import "time"

// ManualTicker - For testing the engine.
// Use this when you don't want the timer to run
// or you want to turn the engine manually.
type ManualTicker struct {
	Ticker chan time.Time
}

// NewManualTicker - Make time channel
func NewManualTicker() *ManualTicker {
	ticker := make(chan time.Time)
	return &ManualTicker{ticker}
}

// GetTicker - Returns the out side of the channel
func (mt *ManualTicker) GetTicker() <-chan time.Time {
	return mt.Ticker
}

// Tick - Fires the in side of the channel.
// Call this to turn the engine one time.
func (mt *ManualTicker) Tick() {
	mt.Ticker <- time.Now()
	return
}

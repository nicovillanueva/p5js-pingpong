package server

import (
	"math/rand"
)

// PingPongMatch is a match between two players
type PingPongMatch struct {
	matchID    int
	kind       string
	redPlayer  int
	bluePlayer int
	plays      []string
}

// Arbiter controls the flow of the matches.
// Do not instantiate directly; see `NewArbiter()`
type Arbiter struct {
	ActiveGames []PingPongMatch
}

// NewArbiter spawns a new game controller/arbiter
func NewArbiter() *Arbiter {
	return &Arbiter{
		ActiveGames: make([]PingPongMatch, 0),
	}
}

// NewPublicMatch is an open match that does not need logging in
func (a *Arbiter) NewPublicMatch() int {
	i := rand.Int() // TODO: Determine id from db
	a.ActiveGames = append(a.ActiveGames, PingPongMatch{
		matchID: i,
		kind:    "public",
		plays:   make([]string, 0),
	})
	return i
}

// NewPrivateMatch starts a match between two user ids
func (a *Arbiter) NewPrivateMatch(red, blue int) int {
	i := rand.Int() // TODO: Determine id from db
	a.ActiveGames = append(a.ActiveGames, PingPongMatch{
		matchID:    i,
		kind:       "private",
		redPlayer:  red,
		bluePlayer: blue,
		plays:      make([]string, 0),
	})
	return i
}

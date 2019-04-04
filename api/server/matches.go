package server

import "errors"

const (
	matchKindOpen    = "open"
	matchKindPrivate = "private"
	matchKindSolo    = "solo"
)

// PingPongMatch is a match between two players
type PingPongMatch struct {
	kind    string
	plays   []Play
	players []int64
	lastBy  int64
}

// Play contains a sketch, it's author and a thumbnail
type Play struct {
	Sketch   string `json:"sketch"`
	Author   int64  `json:"author"`
	Snapshot string `json:",omitempty"`
}

// AddPlay adds a sketch to the game
func (ppm *PingPongMatch) AddPlay(play string, userID int64) {
	ppm.plays = append(ppm.plays, Play{
		Sketch: play,
		Author: userID,
	})
	ppm.lastBy = userID
}

// LastPlay returns the last play and it's author
func (ppm *PingPongMatch) LastPlay() Play {
	if len(ppm.plays) == 0 {
		logger.Errorf("tried to obtain a play from an empty match!")
		return Play{}
	}
	return ppm.plays[len(ppm.plays)-1]
}

// AllPlays pulls all plays
func (ppm *PingPongMatch) AllPlays() []Play {
	return ppm.plays
}

// Play returns a play and it's author by it's index
func (ppm *PingPongMatch) Play(idx int) (Play, error) {
	if idx > len(ppm.plays) {
		return Play{}, errors.New("play not found in match")
	}
	return ppm.plays[idx], nil
}

// IsAllowed finds out if a user id is allowed to post in this game
func (ppm *PingPongMatch) IsAllowed(user int64) bool {
	if ppm.kind == matchKindSolo && user != ppm.players[0] {
		return false
	}
	if ppm.kind == matchKindOpen {
		return true
	}
	for i := 0; i < len(ppm.players); i++ {
		if ppm.players[i] == user {
			return true
		}
	}
	return false
}

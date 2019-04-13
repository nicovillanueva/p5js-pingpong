package server

import (
	"errors"
	"time"
)

// PingPongMatch is a match between N players
type PingPongMatch struct {
	matchID          int64
	ownerID          int64
	startDate        time.Time
	players          PlayersRepo
	requiresApproval bool
	sketches         []Sketch
	lastBy           int64
	theme            string
}

// PlayersRepo is a bunch o'people divided by status
type PlayersRepo struct {
	Current      []int64
	JoinRequests []int64
}

// NewPlayer adds a player to the currently playing players
func (pr *PlayersRepo) NewPlayer(p int64) {
	pr.Current = append(pr.Current, p)
}

// IsQueued says if an id is already in the requests list
func (pr *PlayersRepo) IsQueued(p int64) bool {
	for _, q := range pr.JoinRequests {
		if p == q {
			return true
		}
	}
	return false
}

// IsPlaying says if an id is already playing
func (pr *PlayersRepo) IsPlaying(p int64) bool {
	for _, q := range pr.Current {
		if p == q {
			return true
		}
	}
	return false
}

// NewJoinRequest adds a join request to the match
func (pr *PlayersRepo) NewJoinRequest(p int64) bool {
	if pr.IsQueued(p) || pr.IsPlaying(p) {
		return false
	}
	pr.JoinRequests = append(pr.JoinRequests, p)
	return true
}

// Approve approves a request
func (pr *PlayersRepo) Approve(u int64) {
	if !pr.IsQueued(u) || pr.IsPlaying(u) {
		return
	}
	// remove from requests
	// for i, r := range pr.JoinRequests {
	// 	if r == u {
	// 		// pr.JoinRequests = pr.JoinRequests[]
	// 	}
	// }
	// add to current
}

// // PlayersGroup is a bunch o'people
// type PlayersGroup []int64

// func (pg *PlayersGroup) contains(other int64) bool {
// 	for _, i := range *pg {
// 		if other == i {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (pg *PlayersGroup) remove(p int64) {
// 	// TODO: remove from the queue
// 	logger.Warnf("current requests %+v:", pg)
// }

// Sketch contains a sketch, it's author and a thumbnail
type Sketch struct {
	Sketch   string `json:"sketch"`
	Author   int64  `json:"author"`
	Snapshot string `json:",omitempty"`
}

// AddPlay adds a sketch to the game
func (ppm *PingPongMatch) AddPlay(play string, userID int64) {
	ppm.sketches = append(ppm.sketches, Sketch{
		Sketch: play,
		Author: userID,
	})
	ppm.lastBy = userID
}

// LastPlay returns the last play and it's author
func (ppm *PingPongMatch) LastPlay() Sketch {
	if len(ppm.sketches) == 0 {
		logger.Errorf("tried to obtain a play from an empty match!")
		return Sketch{}
	}
	return ppm.sketches[len(ppm.sketches)-1]
}

// AllPlays pulls all sketches
func (ppm *PingPongMatch) AllPlays() []Sketch {
	return ppm.sketches
}

// Sketch returns a play and it's author by it's index
func (ppm *PingPongMatch) Sketch(idx int) (Sketch, error) {
	if idx > len(ppm.sketches) {
		return Sketch{}, errors.New("play not found in match")
	}
	return ppm.sketches[idx], nil
}

// IsAllowed finds out if a user id is allowed to post in this game
func (ppm *PingPongMatch) IsAllowed(user int64) bool {
	if ppm.requiresApproval && !ppm.players.IsPlaying(user) {
		return false
	}
	return true
}

// AddJoinRequest adds a request to the request pool
func (ppm *PingPongMatch) AddJoinRequest(userID int64) (bool, string) {
	if ppm.players.IsQueued(userID) {
		return false, "request already in queue"
	}
	if ppm.players.IsPlaying(userID) {
		return false, "user already playing"
	}
	ppm.players.NewJoinRequest(userID)
	return true, "request added"
}

// Summary is the public summary (even for anonymous users) of a match
func (ppm *PingPongMatch) Summary() (MatchStatus, error) {
	if ppm.sketches == nil {
		return MatchStatus{}, errors.New("still starting match")
	}
	return MatchStatus{
		MatchID:        ppm.matchID,
		SketchesLength: len(ppm.sketches),
		PlayingUsers:   ppm.players.Current,
		StartedDate:    ppm.startDate,
	}, nil
}

// ApproveRequest moves a request from pending, to playing
func (ppm *PingPongMatch) ApproveRequest(userID int64) {
	if !ppm.players.IsQueued(userID) || ppm.players.IsPlaying(userID) {
		return
	}

	ppm.players.Approve(userID)
	// ppm.players = append(ppm.players, userID)
	// ppm.joinRequests.remove(userID)
}

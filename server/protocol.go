package server

import "time"

// ReturnRequest sends a return
type ReturnRequest struct {
	Sketch string `json:"sketch"`
	UserID int64  `json:"user_id"`
}

// NewMatchRequest starts a match
type NewMatchRequest struct {
	MaxPlayers      int    `json:"max_players"`
	RequireApproval bool   `json:"requires_approval"`
	Sketch          string `json:"sketch"`
	FirstPlayer     int64  `json:"user_id"`
	Theme           string `json:"theme"`
}

// MatchStatus is the public summary of a match
type MatchStatus struct {
	MatchID        int64     `json:"match_id"`
	SketchesLength int       `json:"sketches_length"`
	PlayingUsers   []int64   `json:"players"`
	StartedDate    time.Time `json:"started_on"`
}

// MatchResponse is the response with information about the searched match
type MatchResponse struct {
	MatchStatus
	Status string `json:"status"`
}

// MatchStatusOwner is the summary presented to the match owner
type MatchStatusOwner struct {
	MatchStatus
	PendingRequests []int64 `json:"pending_requests,omitempty"`
}

// MatchWithSketches is a summary of a match, and it's sketches
type MatchWithSketches struct {
	MatchStatus
	Sketches []Sketch `json:"sketches"`
}

// AuthenticatedRequest is a request tied to a given user
type AuthenticatedRequest struct {
	UserID int64 `query:"user" json:"user_id"`
	Token  int64 `query:"token"`
}

// ApproveJoinRequest approves a given request
// There's no denying; at most they go stale
type ApproveJoinRequest struct {
	RequestID int64 `json:"join_request_id"`
}

// JoinStatus says if the joining request was added
type JoinStatus struct {
	Added  bool   `json:"added"`
	Reason string `json:"reason"`
}

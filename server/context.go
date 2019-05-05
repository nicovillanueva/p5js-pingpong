package server

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// MatchContext context that pulls matches
type MatchContext struct {
	echo.Context
}

// CreateMatch creates a match out of it's context
func (mc *MatchContext) CreateMatch() (*PingPongMatch, error) {
	var newMatchReq NewMatchRequest
	if err := mc.Context.Bind(&newMatchReq); err != nil {
		return nil, err
	}
	return &PingPongMatch{
		ownerID:   newMatchReq.FirstPlayer,
		startDate: time.Now(),
		sketches: []Sketch{{
			Sketch: newMatchReq.Sketch,
			Author: newMatchReq.FirstPlayer,
		}},
		players: PlayersRepo{
			Current: []int64{newMatchReq.FirstPlayer},
		},
		requiresApproval: newMatchReq.RequireApproval,
		lastBy:           newMatchReq.FirstPlayer,
		theme:            newMatchReq.Theme,
	}, nil
}

// RetrieveMatch pulls a match from it's context
func (mc *MatchContext) RetrieveMatch() (*PingPongMatch, error) {
	mid, err := mc.matchID()
	if err != nil {
		return nil, err
	}
	m, err := store.FindMatch(mid)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (mc *MatchContext) matchID() (int64, error) {
	mid := mc.Context.Param("matchId")
	matchID, err := strconv.ParseInt(mid, 10, 64)
	if err != nil {
		return int64(0), err
	}
	return matchID, nil
}

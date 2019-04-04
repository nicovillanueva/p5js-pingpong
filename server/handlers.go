package server

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

// ReturnRequest sends a return
type ReturnRequest struct {
	Sketch string `json:"sketch"`
	UserID int64  `json:"user_id"`
}

// MatchRequest starts a match
type MatchRequest struct {
	Sketch      string `json:"sketch"`
	FirstPlayer int64  `json:"user_id"`
	Mode        string `json:"mode"`
}

// UserResponse responds with the user referenced and the status of the request
type UserResponse struct {
	User   string `json:"user"`
	Status string `json:"status"`
}

// MatchResponse is the response with information about the searched match
type MatchResponse struct {
	Status   string `json:"status"`
	MatchID  int64  `json:"match_id"`
	Sketches []Play `json:"sketch,omitempty"`
}

func handleRenderMatch(c echo.Context) error {
	mid := c.Param("matchId")
	i, err := strconv.ParseInt(mid, 10, 64)
	if err != nil {
		return c.JSON(403, MatchResponse{
			Status: err.Error(),
		})
	}
	m, err := store.FindMatch(i)
	if err != nil {
		return c.JSON(403, MatchResponse{
			MatchID: i,
			Status:  err.Error(),
		})
	}
	_, getAll := c.QueryParams()["all"]
	var p []Play
	if getAll {
		p = m.AllPlays()
	} else {
		p = []Play{m.LastPlay()}
	}
	return c.JSON(200, MatchResponse{
		MatchID:  i,
		Status:   "ok",
		Sketches: p,
	})
}

func handleNewMatch(c echo.Context) error {
	var matchReq MatchRequest
	if err := c.Bind(&matchReq); err != nil {
		return c.JSON(501, MatchResponse{
			Status: err.Error(),
		})
	}
	match := PingPongMatch{
		kind: matchReq.Mode,
		plays: []Play{{
			Sketch: matchReq.Sketch,
			Author: matchReq.FirstPlayer,
		}},
		players: []int64{matchReq.FirstPlayer},
	}
	matchID, err := store.SaveMatch(match)
	if err != nil {
		return c.JSON(403, MatchResponse{
			Status: err.Error(),
		})
	}
	return c.JSON(200, MatchResponse{
		MatchID: matchID,
		Status:  "created",
	})
}

func handleReturn(c echo.Context) error {
	var returnRequest ReturnRequest
	mid := c.Param("matchId")
	matchID, err := strconv.ParseInt(mid, 10, 64)
	if err != nil {
		return c.JSON(400, MatchResponse{
			Status: err.Error(),
		})
	}
	if err := c.Bind(&returnRequest); err != nil {
		return c.JSON(400, MatchResponse{
			MatchID: matchID,
			Status:  err.Error(),
		})
	}
	m, err := store.FindMatch(matchID)
	if err != nil {
		return c.JSON(404, MatchResponse{
			MatchID: matchID,
			Status:  err.Error(),
		})
	}
	if !m.IsAllowed(returnRequest.UserID) {
		return c.JSON(406, MatchResponse{
			MatchID: matchID,
			Status:  "user is not a player here (private match)",
		})
	}
	if m.lastBy == returnRequest.UserID {
		return c.JSON(409, MatchResponse{
			MatchID: matchID,
			Status:  "cannot send two sketches in a row",
		})
	}
	m.AddPlay(returnRequest.Sketch, returnRequest.UserID)
	m.lastBy = returnRequest.UserID
	logger.Infof("added play to match %d by user %d, currently %d sketches long", matchID, returnRequest.UserID, len(m.plays))
	return c.JSON(200, MatchResponse{
		MatchID: matchID,
		Status:  "sent",
	})
}

func handleNewUser(c echo.Context) error {
	var payload struct {
		User string `json:"u"`
		Pass string `json:"p"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(400, UserResponse{
			User:   payload.User,
			Status: err.Error(),
		})
	}
	if err := store.StoreUser(payload.User, payload.Pass); err != nil {
		return c.JSON(403, UserResponse{
			User:   payload.User,
			Status: err.Error(),
		})
	}
	return c.JSON(200, UserResponse{
		User:   payload.User,
		Status: "created",
	})
}

func handleJoinMatch(c echo.Context) error {
	logger.Warn("pending method")
	return nil
}

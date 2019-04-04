package server

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

var returnRequest struct {
	Sketch string `json:"sketch"`
	UserID int64  `json:"user_id"`
}

var MatchRequest struct {
	Sketch      string `json:"sketch"`
	FirstPlayer int64  `json:"user_id"`
	Mode        string `json:"mode"`
}

type userResponse struct {
	User   string `json:"user"`
	Status string `json:"status"`
}
type matchResponse struct {
	Status   string `json:"status"`
	MatchID  int64  `json:"match_id"`
	Sketches []Play `json:"sketch,omitempty"`
}

func handleRenderMatch(c echo.Context) error {
	mid := c.Param("matchId")
	i, err := strconv.ParseInt(mid, 10, 64)
	if err != nil {
		return c.JSON(403, matchResponse{
			Status: err.Error(),
		})
	}
	m, err := store.FindMatch(i)
	if err != nil {
		return c.JSON(403, matchResponse{
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
	return c.JSON(200, matchResponse{
		MatchID:  i,
		Status:   "ok",
		Sketches: p,
	})
}

func handleNewMatch(c echo.Context) error {
	if err := c.Bind(&MatchRequest); err != nil {
		return c.JSON(501, matchResponse{
			Status: err.Error(),
		})
	}
	match := PingPongMatch{
		kind: MatchRequest.Mode,
		plays: []Play{{
			Sketch: MatchRequest.Sketch,
			Author: MatchRequest.FirstPlayer,
		}},
		players: []int64{MatchRequest.FirstPlayer},
	}
	matchID, err := store.SaveMatch(match)
	if err != nil {
		return c.JSON(403, matchResponse{
			Status: err.Error(),
		})
	}
	return c.JSON(200, matchResponse{
		MatchID: matchID,
		Status:  "created",
	})
}

func handleReturn(c echo.Context) error {

	mid := c.Param("matchId")
	matchID, err := strconv.ParseInt(mid, 10, 64)
	if err != nil {
		return c.JSON(400, matchResponse{
			Status: err.Error(),
		})
	}
	if err := c.Bind(&returnRequest); err != nil {
		return c.JSON(400, matchResponse{
			MatchID: matchID,
			Status:  err.Error(),
		})
	}
	m, err := store.FindMatch(matchID)
	if err != nil {
		return c.JSON(404, matchResponse{
			MatchID: matchID,
			Status:  err.Error(),
		})
	}
	if !m.IsAllowed(returnRequest.UserID) {
		return c.JSON(406, matchResponse{
			MatchID: matchID,
			Status:  "user is not a player here (private match)",
		})
	}
	if m.lastBy == returnRequest.UserID {
		return c.JSON(409, matchResponse{
			MatchID: matchID,
			Status:  "cannot send two sketches in a row",
		})
	}
	m.AddPlay(returnRequest.Sketch, returnRequest.UserID)
	m.lastBy = returnRequest.UserID
	logger.Infof("added play to match %d by user %d, currently %d sketches long", matchID, returnRequest.UserID, len(m.plays))
	return c.JSON(200, matchResponse{
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
		return c.JSON(400, userResponse{
			User:   payload.User,
			Status: err.Error(),
		})
	}
	if err := store.StoreUser(payload.User, payload.Pass); err != nil {
		return c.JSON(403, userResponse{
			User:   payload.User,
			Status: err.Error(),
		})
	}
	return c.JSON(200, userResponse{
		User:   payload.User,
		Status: "created",
	})
}

func handleJoinMatch(c echo.Context) error {
	logger.Warn("pending method")
	return nil
}

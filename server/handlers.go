package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary Start a new pingpong match
// @Description Receives the match settings and starts it. Includes the serve (first sketch)
// @Accept json
// @Produce json
// @Param match_settings body server.NewMatchRequest true "Starting settings for the match"
// @Success 200 {object} server.MatchResponse
// @Failure 502 {object} server.MatchResponse
// @Failure 403 {object} server.MatchResponse
// @Failure 425 {object} server.MatchResponse
// @Router /match [POST]
func handleNewMatch(ctx echo.Context) error {
	c := ctx.(*MatchContext)
	match, err := c.CreateMatch()
	if err != nil {
		return c.Context.JSON(http.StatusBadGateway, MatchResponse{
			Status: err.Error(),
		})
	}
	if err = store.SaveMatch(match); err != nil {
		// TODO: test errored save
		return c.JSON(http.StatusForbidden, MatchResponse{
			// conflicting IDs
			Status: err.Error(),
		})
	}
	s, err := match.Summary()
	if err != nil {
		return c.JSON(http.StatusTooEarly, MatchResponse{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, MatchResponse{
		Status:      "created",
		MatchStatus: s,
	})

}

func handleMatchStatus(ctx echo.Context) error {
	var err error
	var m *PingPongMatch
	var status MatchStatus
	var authReq AuthenticatedRequest

	c := ctx.(*MatchContext)
	if err = c.Bind(&authReq); err != nil {
		return c.JSON(403, MatchResponse{
			Status: err.Error(),
		})
	}
	m, err = c.RetrieveMatch()
	if err != nil {
		return c.JSON(403, MatchResponse{
			Status: err.Error(),
		})
	}
	if status, err = m.Summary(); err != nil {
		return c.JSON(425, MatchResponse{
			Status: err.Error(),
		})
	}
	if authReq.UserID == m.ownerID && authReq.Token == 1 {
		logger.Debugf("match owner %d requesting match status", authReq.UserID)
		return c.JSON(http.StatusOK, MatchStatusOwner{
			MatchStatus:     status,
			PendingRequests: m.players.JoinRequests,
		})
	}
	logger.Debugf("non-owner %d requesting match status", authReq.UserID)
	return c.JSON(http.StatusOK, status)
}

func handleNewSketch(ctx echo.Context) error {
	c := ctx.(*MatchContext)
	var returnRequest ReturnRequest
	if err := c.Bind(&returnRequest); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	var m *PingPongMatch
	m, err := c.RetrieveMatch()
	if err != nil {
		return err
	}
	logger.Warnf("return req: %+v", returnRequest)
	s, _ := m.Summary()
	if !m.IsAllowed(returnRequest.UserID) {
		return c.JSON(406, MatchResponse{
			MatchStatus: s,
			Status:      "user is not yet a player here",
		})
	}
	if m.lastBy == returnRequest.UserID {
		return c.JSON(409, MatchResponse{
			MatchStatus: s,
			Status:      "cannot send two sketches in a row",
		})
	}
	m.AddPlay(returnRequest.Sketch, returnRequest.UserID)
	return c.JSON(200, MatchResponse{
		MatchStatus: s,
		Status:      "received",
	})
}

func handleJoinMatch(ctx echo.Context) error {
	c := ctx.(*MatchContext)
	m, err := c.RetrieveMatch()
	if err != nil {
		return c.JSON(403, MatchResponse{
			Status: err.Error(),
		})
	}
	var authReq AuthenticatedRequest
	if err := c.Bind(&authReq); err != nil {
		return c.JSON(403, MatchResponse{
			Status: err.Error(),
		})
	}
	added, reason := m.AddJoinRequest(authReq.UserID)
	if !added {
		return c.JSON(409, JoinStatus{added, reason})
	}
	return c.JSON(200, JoinStatus{added, reason})
}

func handleApproveRequests(ctx echo.Context) error {
	c := ctx.(*MatchContext)
	var joinRequest *ApproveJoinRequest
	var match *PingPongMatch
	err := c.Bind(&joinRequest)
	if err != nil {
		return err
	}
	match, err = c.RetrieveMatch()
	if err != nil {
		return c.Context.JSON(404, MatchResponse{
			Status: err.Error(),
		})
	}
	s, _ := match.Summary()
	match.ApproveRequest(joinRequest.RequestID)
	return c.JSON(200, s)
}

// get N sketches
func handleGetSketches(ctx echo.Context) error {
	c := ctx.(*MatchContext)
	match, err := c.RetrieveMatch()
	if err != nil {
		return err
	}
	// s, _ := match.Summary()
	// return c.JSON(200, MatchWithSketches{
	return c.JSON(200, SketchList{
		// MatchStatus: s,
		Sketches: []Sketch{match.LastPlay()},
	})
}

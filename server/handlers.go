package server

import (
	"github.com/labstack/echo/v4"
)

func handleNewMatch(ctx echo.Context) error {
	c := ctx.(*MatchContext)
	match, err := c.CreateMatch()
	if err != nil {
		return err
	}
	if err = store.SaveMatch(match); err != nil {
		// TODO: test errored save
		return c.JSON(403, MatchResponse{
			Status: err.Error(),
		})
	}
	s, err := match.Summary()
	if err != nil {
		return c.JSON(425, MatchResponse{
			Status: err.Error(),
		})
	}
	return c.JSON(200, MatchResponse{
		Status:      "created",
		MatchStatus: s,
	})

}

func handleMatchStatus(ctx echo.Context) error {
	var err error
	var m *PingPongMatch
	var status MatchStatus
	c := ctx.(*MatchContext)
	var authReq AuthenticatedRequest
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
		return c.JSON(200, MatchStatusOwner{
			MatchStatus:     status,
			PendingRequests: m.players.JoinRequests,
		})
	}
	logger.Debugf("non-owner %d requesting match status", authReq.UserID)
	return c.JSON(200, status)
}

func handleNewSketch(ctx echo.Context) error {
	c := ctx.(*MatchContext)
	var returnRequest ReturnRequest
	if err := c.Bind(&returnRequest); err != nil {
		return c.JSON(500, err)
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
		return err
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
	s, _ := match.Summary()
	return c.JSON(200, MatchWithSketches{
		MatchStatus: s,
		Sketches:    []Sketch{match.LastPlay()},
	})
}

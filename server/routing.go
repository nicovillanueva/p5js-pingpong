package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

// APIVersion is a specific API version with all it's handlers
// Call `.apply()` to enable the routes in the server
type APIVersion struct {
	Number    int
	Enabled   bool
	Handlings []routeDefinition
}

type routeDefinition struct {
	Method  string
	Path    string
	Handler func(echo.Context) error
}

func (a *APIVersion) apply(e *echo.Echo) {
	if !a.Enabled {
		return
	}
	g := e.Group(fmt.Sprintf("/v%d", a.Number))
	for _, r := range a.Handlings {
		g.Add(r.Method, r.Path, r.Handler)
	}
}

func setupRoutes(e *echo.Echo) {
	// Root
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		})
	})

	a := APIVersion{
		Number:  1,
		Enabled: true,
		Handlings: []routeDefinition{
			{"POST", "/match", handleNewMatch},
			{"GET", "/match/:matchId", handleMatchStatus},
			{"PUT", "/match/:matchId/players", handleJoinMatch},
			{"PATCH", "/match/:matchId/players", handleApproveRequests},
			{"PUT", "/match/:matchId/sketches", handleNewSketch},
			{"GET", "/match/:matchId/sketches", handleGetSketches},
		},
	}
	a.apply(e)
}

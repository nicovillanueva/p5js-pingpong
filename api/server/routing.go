package server

import (
	"github.com/labstack/echo/v4"
)

var apiGroups map[string]*echo.Group

func setupRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		})
	})

	apiGroups = make(map[string]*echo.Group)
	apiGroups["users"] = e.Group("/users")
	apiGroups["matches"] = e.Group("/matches")

	setupUsers()
	setupMatches()
}

func setupUsers() {
	apiGroups["users"].Add("PUT", "/", func(c echo.Context) error {
		var u user
		if err := c.Bind(&u); err != nil {
			c.Echo().Logger.Warnf("could not bind data: %v", err)
		}
		uid, err := u.save()
		if err != nil {
			return c.JSON(500, errCannotSave{})
		}
		return c.JSON(200, struct {
			UserID int `json:"user_id"`
		}{uid})
	})
}

func setupMatches() {
	apiGroups["matches"].Add("GET", "/", func(c echo.Context) error {
		// TODO: template out match
		return nil
	})

	apiGroups["matches"].Add("PUT", "/serve/public", func(c echo.Context) error {
		// TODO: new public match
		mid := activeArbiter.NewPublicMatch()
		return c.JSON(200, struct {
			MatchID int `json:"match_id"`
		}{mid})
	})

	// apiGroups["matches"].Add("PUT", "/serve/private", func(c echo.Context) error {
	// 	// TODO: new private match
	// 	activeArbiter.NewPrivateMatch(c.Get("userid"), c.Get(""))
	// 	return nil
	// })

	apiGroups["matches"].Add("PATCH", "/return", func(c echo.Context) error {
		// TODO: new entry for a match
		return nil
	})
}

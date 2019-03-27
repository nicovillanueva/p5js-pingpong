package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/gommon/log"
)

// Verbose defines the logging level. True goes into Debug, Info otherwise.
// Set previous to starting the server
const Verbose = false

var activeArbiter *Arbiter

func init() {
	activeArbiter = NewArbiter()
}

// Start serves the server
func Start() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Logger.SetLevel(func() log.Lvl {
		if Verbose {
			return log.DEBUG
		}
		return log.INFO
	}())

	setupRoutes(e)

	e.Logger.Debug("active routes:")
	for _, r := range e.Routes() {
		e.Logger.Debugf("%+v", *r)
	}

	e.Logger.Fatal(e.Start(":8000"))
}

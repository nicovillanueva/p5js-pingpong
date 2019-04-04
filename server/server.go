package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/gommon/log"
)

// MatchRequest starts a match
// type MatchRequest struct {
// 	Sketch      string `json:"sketch"`
// 	FirstPlayer int64  `json:"user_id"`
// 	Mode        string `json:"mode"`
// }

// Verbose defines the logging level. True goes into Debug, Info otherwise.
// Set previous to starting the server
const Verbose = true

var store *DataStore
var logger echo.Logger

func init() {
	store = NewDataStore()
}

// Start serves the server
func Start() {
	e := echo.New()
	logger = e.Logger

	e.Use(middleware.Recover())
	e.Logger.SetLevel(func() log.Lvl {
		if Verbose {
			return log.DEBUG
		}
		return log.INFO
	}())
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	setupRoutes(e)

	e.Logger.Debug("active routes:")
	for _, r := range e.Routes() {
		e.Logger.Debugf("%+v", *r)
	}

	e.Logger.Fatal(e.Start(":8000"))
}

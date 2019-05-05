package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/gommon/log"

	_ "github.com/nicovillanueva/p5js-pingpong/server/docs" // generated swagger docs
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Verbose defines the logging level. True goes into Debug, Info otherwise.
// Set previous to starting the server
const Verbose = true

var store *DataStore
var logger echo.Logger

func init() {
	store = NewDataStore()
}

// Start serves the server
// @title PingPong API
// @version 0.1
// @description The p5jspingpong API
// @host localhost:8000
// @BasePath /
func Start() {
	e := echo.New()
	logger = e.Logger

	// custom context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &MatchContext{c}
			return h(cc)
		}
	})

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
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8000"))
}

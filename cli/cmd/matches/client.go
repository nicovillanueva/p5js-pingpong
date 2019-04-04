package cmd

import (
	"github.com/nicovillanueva/p5js-pingpong/api/server"
	"net/http"
)

const baseURL = "http://localhost:8000"

func requestNewMatch() {
	// var req struct {
	// 	sketch
	// 	player
	// 	mode
	// }
	// server.MatchRequest
	http.PUT(baseURL + "/matches")
}

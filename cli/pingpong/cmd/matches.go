package cmd

import (
	// "bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Send a sketch for an open match",
	Run: func(cmd *cobra.Command, args []string) {
		newMatch()
	},
}
var returnCmd = &cobra.Command{
	Use:   "return",
	Short: "Return a hit in a given match",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func newMatch() {
	var endpoint = "http://localhost:8000/"
	client := &http.Client{}
	// r, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer([]byte(`testing`)))
	r, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}

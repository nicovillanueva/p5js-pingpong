package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(matchCmd)
}

// rootCmd represents the base command when called without any subcommands
var matchCmd = &cobra.Command{
	Use: "match",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("from match")
	},
}

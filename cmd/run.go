/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	internal "github.com/aliamerj/ravly/internal/discovery"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Ravly server",
	RunE:   runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
  fmt.Println("Ravly is now running..")

  	for {
		if err := internal.BroadcastPresence(); err != nil {
			return fmt.Errorf("Error sending upd: %w", err)
		}
		time.Sleep(3 * time.Second)
	}

}

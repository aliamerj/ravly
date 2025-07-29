/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/aliamerj/ravly/internal/discovery"
	"github.com/aliamerj/ravly/internal/transfer"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Ravly server",
	Run:   runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) {
	fmt.Println("Ravly is now running..")
	go discoveryBroadcast()
	go fileReceiver()

  select{}
}

func discoveryBroadcast() {
	for {
		if err := discovery.BroadcastPresence(); err != nil {
			fmt.Println("❌ Error broadcasting:", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func fileReceiver() {
	if err := transfer.StartServer("0.0.0.0:8989"); err != nil {
		fmt.Println("❌ Error receiving files:", err)
	}
}

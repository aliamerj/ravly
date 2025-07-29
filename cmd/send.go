/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/aliamerj/ravly/internal/transfer"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [ip] [file]",
	Short: "Send a file to another device",
  Args: cobra.ExactArgs(2),
	RunE: runSend,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}

func runSend(cmd *cobra.Command, args []string) error {
  ip := args[0]
  filePath := args[1]
  return  transfer.SendFile(ip, filePath)
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	internal "github.com/aliamerj/ravly/internal/discovery"
	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover other Ravly nodes on your network",
	RunE:  runDiscover,
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}

func runDiscover(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Discovering Ravly peers...")
	if err := internal.ListenForPeers(func(addr string) {
		fmt.Printf("👋 Found peer at %s\n", addr)
	}); err != nil {
		return fmt.Errorf("Error discover: %w", err)
	}
	return nil
}

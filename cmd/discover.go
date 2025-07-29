/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/aliamerj/ravly/internal/discovery"
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
	fmt.Println("🌟 Ravly Peer Discovery")
	fmt.Println("-----------------------")
	fmt.Println("🔍 Listening for peers... (Ctrl+C to exit)")

	// Track unique peers with last seen time
	peers := make(map[string]string)
	lastSeen := make(map[string]time.Time)
	var mu sync.Mutex
	var count int

	return discovery.ListenForPeers(func(addr, host string) {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()

		if _, exists := peers[addr]; !exists {
			peers[addr] = host
			count++
		}
		lastSeen[addr] = now

		clearPreviousOutput(count)

		fmt.Println("Discovered peers:")
		fmt.Println("-----------------")

		ips := make([]string, 0, len(peers))
		for ip := range peers {
			ips = append(ips, ip)
		}
		sort.Strings(ips)

		for _, ip := range ips {
			elapsed := time.Since(lastSeen[ip]).Round(time.Second)
			fmt.Printf("🏠 %-15s %-20s (seen %v ago)\n", ip+":", peers[ip], elapsed)
		}

		fmt.Printf("\n🔍 Listening for peers... (found %d, Ctrl+C to exit)", count)
	})
}

func clearPreviousOutput(lineCount int) {
	// Move cursor up and clear from cursor to end of screen
	// 5 extra lines for headers/footers
	fmt.Printf("\033[%dA", lineCount+5)
	fmt.Print("\033[J")
}

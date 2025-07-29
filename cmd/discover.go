/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aliamerj/ravly/daemon"
	"github.com/aliamerj/ravly/daemon/proto"
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
	conn, err := daemon.Connect(DAEMON_ADDRESS)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := proto.NewDaemonServiceClient(conn)

	fmt.Println("🌟 Live Discovery (press Ctrl+C to stop)")
	fmt.Println("----------------------------------------")

	seen := make(map[string]bool)

	for {
		resp, err := client.ListPeers(context.Background(), &proto.Empty{})
		if err != nil {
			fmt.Println("❌ Error:", err)
			break
		}

		for _, p := range resp.Peers {
			if !seen[p.Ip] {
				fmt.Printf("🔎 New peer: %-15s %-20s\n", p.Ip, p.Hostname)
				seen[p.Ip] = true
			}
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}

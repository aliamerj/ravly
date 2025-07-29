/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/aliamerj/ravly/daemon"
	"github.com/aliamerj/ravly/daemon/proto"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send [ip] [file]",
	Short: "Send a file to another device",
	Args:  cobra.ExactArgs(2),
	RunE:  runSend,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}

func runSend(cmd *cobra.Command, args []string) error {
	ip := args[0]
	filePath := args[1]

	conn, err := daemon.Connect(DAEMON_ADDRESS)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := proto.NewDaemonServiceClient(conn)
	if _, err := client.SendFile(context.Background(), &proto.SendFileRequest{
		FilePath: filePath,
		TargetIp: ip,
	}); err != nil {
		return err
	}
	return nil
}

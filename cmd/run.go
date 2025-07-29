/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/aliamerj/ravly/daemon"
	"github.com/aliamerj/ravly/daemon/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	name          string
	transferPort  int
	discoveryPort int
	recvDir       string
	autoAccept    bool

	// runCmd represents the run command
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start Ravly services with configuration",
		RunE:  runRun,
	}
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		recvDir = "."
	}
	defaultDownloadDir := filepath.Join(homeDir, "Downloads")

	runCmd.PersistentFlags().StringVarP(&name, "name", "n", hostname, "Device display name")
	runCmd.PersistentFlags().IntVarP(&transferPort, "transfer-port", "t", 9898, "File transfer port")
	runCmd.PersistentFlags().IntVarP(&discoveryPort, "discovery-port", "p", 9999, "Discovery port")
	runCmd.PersistentFlags().StringVarP(&recvDir, "recv-dir", "r", defaultDownloadDir, "Receiving directory")
	runCmd.PersistentFlags().BoolVarP(&autoAccept, "auto-accept", "a", false, "Auto-accept transfers")
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	if err := startDaemonServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start daemon: %v\n", err)
		os.Exit(1)
	}

	conn, err := daemon.Connect(DAEMON_ADDRESS)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon after retries: %w", err)
	}
	defer conn.Close()

	cfg := proto.Config{
		Name:          name,
		Port:          int32(transferPort),
		DiscoveryPort: int32(discoveryPort),
		RecvDir:       recvDir,
		AutoAccept:    autoAccept,
	}

	client := proto.NewDaemonServiceClient(conn)

	// Send the config to the daemon
	if _, err := client.SetConfig(context.Background(), &cfg); err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}
	go client.BroadcastPresence(context.Background(), &proto.Empty{})
  
    // TODO: show the Receiving files
	 client.ListReceivedFiles(context.Background(), &proto.Empty{})

	fmt.Println("✅ Config successfully set on the daemon.")
	fmt.Println("   Name:", cfg.Name)
	fmt.Println("   Port:", cfg.Port)
	fmt.Println("   Discovery Port:", cfg.DiscoveryPort)
	fmt.Println("   Receive Dir:", cfg.RecvDir)
	fmt.Println("   Auto-accept:", cfg.AutoAccept)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("🛑 Shutting down CLI gracefully.")

	return nil
}

func startDaemonServer() error {
	lis, err := net.Listen("tcp", DAEMON_ADDRESS)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	srv := daemon.New()
	proto.RegisterDaemonServiceServer(grpcServer, srv)

	go func() {
		fmt.Printf("🚀 Daemon started on %s\n", DAEMON_ADDRESS)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Fprintf(os.Stderr, "gRPC server stopped: %v\n", err)
		}
	}()

	return nil
}

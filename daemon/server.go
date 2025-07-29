package daemon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aliamerj/ravly/daemon/proto"
	"github.com/aliamerj/ravly/internal/discovery"
	"github.com/aliamerj/ravly/internal/transfer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	proto.UnimplementedDaemonServiceServer

	mu      sync.RWMutex
	config  *proto.Config
	tracker *transfer.FileTracker
}

func New() *Server {
	return &Server{
		tracker: &transfer.FileTracker{},
	}
}

func Connect(address string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	for i := range 5 {
		conn, err = grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Failed to Start daemon, try  ", i+1)
	}

	return conn, err
}

func (s *Server) ListPeers(ctx context.Context, _ *proto.Empty) (*proto.PeerList, error) {
	var (
		mu    sync.Mutex
		peers = make(map[string]proto.Peer)
	)

	// Listen briefly and collect discovered peers
	done := make(chan struct{})
	go func() {
		_ = discovery.ListenForPeers(int(s.config.DiscoveryPort), func(p discovery.Peer) {
			mu.Lock()
			peers[p.Ip] = proto.Peer{
				Ip:       p.Ip,
				Hostname: p.Name,
				LastSeen: time.Now().Format(time.RFC3339),
			}
			mu.Unlock()
		})
		close(done)
	}()

	// Wait a short time then return the peers (e.g. 3 seconds)
	select {
	case <-time.After(3 * time.Second):
		// enough time to receive some peers
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mu.Lock()
	defer mu.Unlock()
	var list proto.PeerList
	for _, p := range peers {
		list.Peers = append(list.Peers, &p)
	}
	return &list, nil
}

func (s *Server) SetConfig(ctx context.Context, cfg *proto.Config) (*proto.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = cfg
	return &proto.Empty{}, nil
}

func (s *Server) GetConfig(ctx context.Context, _ *proto.Empty) (*proto.Config, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.config == nil {
		return nil, fmt.Errorf("Config not set")
	}
	return s.config, nil
}

func (s *Server) BroadcastPresence(ctx context.Context, _ *proto.Empty) (*proto.Empty, error) {
	for {
		if err := discovery.BroadcastPresence(s.config.Name, int(s.config.DiscoveryPort)); err != nil {
			return &proto.Empty{}, nil
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *Server) SendFile(ctx context.Context, req *proto.SendFileRequest) (*proto.Empty, error) {
	if err := transfer.SendFile(req.TargetIp, req.FilePath, int(s.config.Port)); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *Server) ListReceivedFiles(ctx context.Context, _ *proto.Empty) (*proto.ReceivedFileList, error) {
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Port)

	go func() {
		if err := transfer.StartServer(ctx, addr, s.config.RecvDir, s.tracker); err != nil {
			fmt.Printf("❌ File receive server error: %v\n", err)
		}
	}()

	files := s.tracker.List()
	var protoFiles []*proto.ReceivedFile
	for _, f := range files {
		protoFiles = append(protoFiles, &proto.ReceivedFile{
			FileName:   f.Name,
			SenderIp:   f.Sender,
			ReceivedAt: time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	return &proto.ReceivedFileList{Files: protoFiles}, nil
}

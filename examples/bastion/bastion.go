package bastion

import (
	"github.com/AnyISalIn/yrpc"
	shared "github.com/AnyISalIn/yrpc/shared"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var serverLogger = log.New(os.Stdout, "[server] ", shared.LogFlags)

func (a *agent) ForwardStream(serviceMethod string, rwc io.ReadWriteCloser) error {
	//serverLogger.Printf("forwarding %s", serviceMethod)
	conn, err := a.peer.Stream(serviceMethod)
	if err != nil {
		serverLogger.Printf("failed to stream agent exec %v", err)
		return err
	}
	defer conn.Close()

	//serverLogger.Printf("bridging %T <-> %T", rwc, conn)

	Bridge(conn, rwc)

	return nil
}

type Server struct {
	agentMap  map[string]*agent
	agentLock sync.Mutex
}

func NewServer() *Server {
	return &Server{
		agentMap:  map[string]*agent{},
		agentLock: sync.Mutex{},
	}
}

func (s *Server) Serve(listener net.Listener) {
	srv := yrpc.NewServer(nil)
	bastion := new(Bastion)
	bastion.agentMap = s.agentMap
	if err := srv.Register(bastion); err != nil {
		serverLogger.Fatalf("failed to register %v", err)
	}

	for {
		peer, err := srv.AcceptPeer(listener)
		if err != nil {
			continue
		}
		if s.isAgent(peer) {
			go s.handleAgent(peer)
		}
	}
}

func (s *Server) isAgent(peer *yrpc.Peer) bool {
	var role string
	if err := peer.Call(AGENT_ROLE, &shared.Empty{}, &role); err != nil {
		serverLogger.Fatalf("failed to get client role %v", err)
	}
	return role == "agent"
}

func (s *Server) handleAgent(peer *yrpc.Peer) {
	var agentName string
	if err := peer.Call(AGENT_WHOAMI, &shared.Empty{}, &agentName); err != nil {
		serverLogger.Fatalf("failed to get agent name %v", err)
	}
	serverLogger.Printf("accept agent %s", agentName)

	s.agentLock.Lock()
	agt, got := s.agentMap[agentName]
	if !got {
		agt = &agent{peer: peer, name: agentName}
	} else {
		agt.peer = peer
	}
	s.agentMap[agentName] = agt
	s.agentLock.Unlock()

	//buf := bufio.NewReadWriter([]byte{})
	//if err := agt.PerformCmd("ls", false, buf); err != nil {
	//	serverLogger.Printf("failed to perform cmd on %s", agentName)
	//} else {
	//	serverLogger.Printf(buf.String())
	//}

}

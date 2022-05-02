package bastion

import (
	"encoding/gob"
	"github.com/AnyISalIn/yrpc"
	shared "github.com/AnyISalIn/yrpc/shared"
	"io"
	"log"
	"net"
	"sync"
)

func (a *agent) PerformCmd(cmd string, tty bool, rwc io.ReadWriteCloser) error {
	conn, err := a.peer.Stream(AGENT_EXEC)
	if err != nil {
		log.Printf("[server] failed to stream agent exec %v", err)
		return err
	}
	defer conn.Close()

	log.Printf("[server] stream exec on agent %s", a.name)

	args := &ExecArgs{Cmd: cmd, TTY: tty}

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(args); err != nil {
		log.Printf("[server] failed to encode agent exec %v", err)
		return err
	}

	log.Printf("[server] encoded args on agent %s, waiting reply", a.name)

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
		log.Fatalf("[server] failed to register %v", err)
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
		log.Fatalf("[server] failed to get client role %v", err)
	}
	return role == "agent"
}

func (s *Server) handleAgent(peer *yrpc.Peer) {
	var agentName string
	if err := peer.Call(AGENT_WHOAMI, &shared.Empty{}, &agentName); err != nil {
		log.Fatalf("[server] failed to get agent name %v", err)
	}
	log.Printf("[server] accept agent %s", agentName)

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
	//	log.Printf("[server] failed to perform cmd on %s", agentName)
	//} else {
	//	log.Printf(buf.String())
	//}

}

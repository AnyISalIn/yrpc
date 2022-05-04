package yrpc

import (
	"github.com/AnyISalIn/yrpc/rpc"
	"github.com/AnyISalIn/yrpc/rpc/yrpc"
	"github.com/hashicorp/yamux"
	"log"
	"net"
	"os"
)

type Server struct {
	config     *ServerConfig
	accepted   bool
	connCn     chan net.Conn
	registered []any
}

func NewServer(config *ServerConfig) *Server {
	if config == nil {
		config = DefaultServerConfig()
	}
	return &Server{config: config, accepted: false, connCn: make(chan net.Conn, 32)}
}

type ServerConfig struct {
	Impl        rpc.Interface
	Logger      *log.Logger
	YamuxConfig *yamux.Config
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{Impl: yrpc.New(), Logger: log.New(os.Stdout, "[server] ", 0)}
}

func (s *Server) Register(any any) error {
	return s.config.Impl.Register(any)
}

func (s *Server) AcceptPeer(lis net.Listener) (*Peer, error) {
	if !s.accepted {
		go func() {
			for {
				conn, err := lis.Accept()
				if err != nil {
					s.config.Logger.Printf("failed to accept conn %v", err)
					continue
				}
				s.connCn <- conn
			}
		}()
		s.accepted = true
	}

WAIT:
	select {
	case conn := <-s.connCn:
		session, err := yamux.Server(conn, s.config.YamuxConfig)
		if err != nil {
			s.config.Logger.Printf("failed to create session %v", err)
			goto WAIT
		}
		peer, err := NewPeer(session, s.config.Impl)
		if err != nil {
			return nil, err
		}
		return peer, err
	}
}

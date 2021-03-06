package yrpc

import (
	internal "github.com/AnyISalIn/yrpc/rpc"
	shared "github.com/AnyISalIn/yrpc/shared"
	"github.com/hashicorp/yamux"
	"github.com/ugorji/go/codec"
	"io"
	"log"
	"net"
	"os"
)

type Peer struct {
	session  *yamux.Session
	rpcImp   internal.Interface
	logger   *log.Logger
	shutdown bool
}

const (
	PeerCall uint8 = iota
	PeerStream
)

type PeerRequest struct {
	Type uint8
}

func NewPeerWithDial(network string, addr string, yamuxConfig *yamux.Config, rpcImp internal.Interface, any ...any) (peer *Peer, err error) {
	peer = &Peer{rpcImp: rpcImp, logger: log.New(os.Stdout, "[peer] ", shared.LogFlags)}
	session, err := dial(network, addr, yamuxConfig)
	if err != nil {
		return nil, err
	}
	peer.session = session

	//for _, a := range any {
	//	if err := rpcImp.Register(a); err != nil {
	//		return nil, err
	//	}
	//}

	go peer.run()
	return peer, nil
}

func NewPeer(session *yamux.Session, rpcImp internal.Interface, any ...any) (peer *Peer, err error) {
	peer = &Peer{rpcImp: rpcImp, logger: log.New(os.Stdout, "[peer] ", shared.LogFlags)}

	//for _, a := range any {
	//	if err := rpcImp.Register(a); err != nil {
	//		return nil, err
	//	}
	//}

	peer.session = session
	go peer.run()
	return peer, nil
}

// accept rpc call
func (p *Peer) run() {
	for !p.shutdown {
		stream, err := p.session.AcceptStream()
		if err != nil {
			if err == io.EOF {
				p.logger.Printf("%v", err)
				return
			}
			p.logger.Printf("failed to accept stream %v", err)
			return
		}

		go func() {
			defer stream.Close()
			req := new(PeerRequest)
			decoder := codec.NewDecoder(stream, &codec.MsgpackHandle{})

			if err := decoder.Decode(req); err != nil {
				p.logger.Printf("failed to decode %v", err)
				return
			}

			p.logger.Printf("received request %s", req.Type)

			//ack := new(shared.ACK)
			//if err := encoder.Encode(ack); err != nil {
			//	p.logger.Printf("failed to encode %v", err)
			//	return
			//}

			if req.Type == PeerCall {
				p.logger.Printf("%s -> %s [%d] calling", stream.LocalAddr(), stream.RemoteAddr(), stream.StreamID())
				if err := p.rpcImp.RemoteCall(stream); err != nil {
					p.logger.Printf("failed to call %v", err)
				}
			} else if req.Type == PeerStream {
				p.logger.Printf("%s -> %s [%d] streaming", stream.LocalAddr(), stream.RemoteAddr(), stream.StreamID())
				if err := p.rpcImp.RemoteStream(stream); err != nil {
					p.logger.Printf("failed to stream %v", err)
				}
			} else {
				panic("")
			}
		}()
	}
}

func (p *Peer) Call(serviceMethod string, args any, reply any) error {
	stream, err := p.session.OpenStream()
	defer stream.Close()

	if err != nil {
		return err
	}

	req := new(PeerRequest)
	req.Type = PeerCall
	encoder := codec.NewEncoder(stream, &codec.MsgpackHandle{})
	if err := encoder.Encode(req); err != nil {
		return err
	}

	if err := p.rpcImp.Call(stream, serviceMethod, args, reply); err != nil {
		return err
	}
	return nil
}

func (p *Peer) Stream(serviceMethod string) (io.ReadWriteCloser, error) {
	stream, err := p.session.OpenStream()
	//defer stream.Close()

	if err != nil {
		return nil, err
	}
	// trigger serve conn
	req := new(PeerRequest)
	req.Type = PeerStream
	encoder := codec.NewEncoder(stream, &codec.MsgpackHandle{})

	//decoder := gob.NewDecoder(stream)
	if err := encoder.Encode(req); err != nil {
		return nil, err
	}

	p.logger.Printf("sending request %s", req.Type)

	//ack := new(shared.ACK)
	//if err := decoder.Decode(ack); err != nil {
	//	return nil, err
	//}

	if err := p.rpcImp.Stream(stream, serviceMethod); err != nil {
		stream.Close()
		p.logger.Printf("failed to stream %v", err)
	}
	return stream, nil
}

func (p *Peer) Close() error {
	p.shutdown = true
	return p.session.Close()
}

func dial(network string, addr string, config *yamux.Config) (*yamux.Session, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	session, err := yamux.Client(conn, config)
	if err != nil {
		return nil, err
	}
	return session, nil
}

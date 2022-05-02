package yrpc

import (
	"encoding/gob"
	shared "github.com/AnyISalIn/yrpc/shared"
	"io"
	"log"
	"net/rpc"
	"reflect"
	"strings"
)

type YRPC struct {
	rpcServer     *rpc.Server
	streamHandler *Handler
}

func New() *YRPC {
	return &YRPC{rpcServer: rpc.NewServer(), streamHandler: NewHandler()}
}

func (y *YRPC) Register(revr any) error {
	if err := y.tryRegisterStream(revr); err != nil {
		return err
	}
	if err := y.tryResterCall(revr); err != nil {
		if strings.Contains(err.Error(), "exported methods of suitable type") {
			return nil
		}
		return err
	}
	return nil
}

func (y *YRPC) tryResterCall(revr any) error {
	return y.rpcServer.Register(revr)
}

func (y *YRPC) tryRegisterStream(revr any) error {
	handlers, err := suitableStreamMethods(reflect.TypeOf(revr), reflect.ValueOf(revr))
	if err != nil {
		return err
	}
	for name, handler := range handlers {
		if err := y.streamHandler.Register(name, handler); err != nil {
			return err
		}
	}
	return nil
}

func (y *YRPC) RemoteCall(conn io.ReadWriteCloser) {
	y.rpcServer.ServeConn(conn)
}

func (y *YRPC) RemoteStream(conn io.ReadWriteCloser) {
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	req := new(Request)
	if err := decoder.Decode(req); err != nil {
		log.Printf("[yrpc] can't decode conn %v", err)
	}
	ack := new(shared.ACK)
	if err := encoder.Encode(ack); err != nil {
		log.Printf("[yrpc] can't encode conn %v", err)
	}

	if err := y.streamHandler.ServeConn(conn, req); err != nil {
		log.Printf("[yrpc] can't serve stream conn %v", err)
	}
}

func (y *YRPC) Call(conn io.ReadWriteCloser, serviceMethod string, args any, reply any) error {
	cli := rpc.NewClient(conn)
	return cli.Call(serviceMethod, args, reply)
}

func (y *YRPC) Stream(conn io.ReadWriteCloser, serviceMethod string) error {
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	req := new(Request)
	req.ServiceMethod = serviceMethod
	if err := encoder.Encode(req); err != nil {
		return err
	}
	ack := new(shared.ACK)
	if err := decoder.Decode(ack); err != nil {
		return err
	}

	return nil
}

func (y *YRPC) Supports() []string {
	panic("implement me")
}

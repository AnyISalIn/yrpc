package yrpc

import (
	"errors"
	"fmt"
	shared "github.com/AnyISalIn/yrpc/shared"
	"github.com/ugorji/go/codec"
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
	log.SetFlags(shared.LogFlags)
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

func (y *YRPC) RemoteCall(conn io.ReadWriteCloser) error {
	y.rpcServer.ServeConn(conn)
	return nil
}

func (y *YRPC) RemoteStream(conn io.ReadWriteCloser) error {
	decoder := codec.NewDecoder(conn, &codec.MsgpackHandle{})

	log.Printf("remote stream received req\n")
	req := new(Request)
	if err := decoder.Decode(req); err != nil {
		return errors.New(fmt.Sprintf("[yrpc] can't decode conn %v", err))
	}

	log.Printf("remote stream serve conn\n")
	return y.streamHandler.ServeConn(conn, req)
}

func (y *YRPC) Call(conn io.ReadWriteCloser, serviceMethod string, args any, reply any) error {
	cli := rpc.NewClient(conn)
	return cli.Call(serviceMethod, args, reply)
}

func (y *YRPC) Stream(conn io.ReadWriteCloser, serviceMethod string) error {
	encoder := codec.NewEncoder(conn, &codec.MsgpackHandle{})

	req := new(Request)
	req.ServiceMethod = serviceMethod
	if err := encoder.Encode(req); err != nil {
		return err
	}
	return nil
}

func (y *YRPC) Supports() []string {
	panic("implement me")
}

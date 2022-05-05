package yrpc

import (
	"errors"
	"fmt"
	"github.com/AnyISalIn/yrpc/rpc"
	"io"
	"sync"
)

type Request struct {
	ServiceMethod string
	//Seq           uint64
}

type Handler struct {
	handlers map[string]rpc.StreamHandler
	mutex    sync.Mutex
}

func NewHandler() *Handler {
	return &Handler{handlers: make(map[string]rpc.StreamHandler), mutex: sync.Mutex{}}
}

func (s *Handler) Register(name string, fn rpc.StreamHandler) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.handlers[name] = fn
	return nil
}

func (s *Handler) ServeConn(conn io.ReadWriteCloser, req *Request) error {
	s.mutex.Lock()
	method, got := s.handlers[req.ServiceMethod]
	s.mutex.Unlock()

	if !got {
		return fmt.Errorf("can't found servicemethod %s", req.ServiceMethod)
	}

	if err := method(conn); err != nil {
		return errors.New(fmt.Sprintf("[streamHandler] method return errors %v", err))
	}
	return nil
}

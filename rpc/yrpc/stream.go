package yrpc

import (
	"fmt"
	"github.com/AnyISalIn/yrpc/rpc"
	"io"
	"log"
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

	go func() {
		if err := method(conn); err != nil {
			log.Printf("[streamHandler] method return errors %v", err)
		}
	}()
	return nil
}

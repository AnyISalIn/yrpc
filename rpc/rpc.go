package rpc

import "io"

type StreamHandler func(rwc io.ReadWriteCloser) error

type Server interface {
	RemoteCall(conn io.ReadWriteCloser)
	RemoteStream(conn io.ReadWriteCloser)

	Register(method any) error
}

type Client interface {
	Call(conn io.ReadWriteCloser, serviceMethod string, args any, reply any) error
	Stream(conn io.ReadWriteCloser, serviceMethod string) error
}

type Interface interface {
	Server
	Client
	Supports() []string // inspect supports methods
}

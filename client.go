package yrpc

import (
	"github.com/AnyISalIn/yrpc/rpc"
	"github.com/AnyISalIn/yrpc/rpc/yrpc"
	"github.com/hashicorp/yamux"
	"log"
	"os"
)

type Client struct {
	config     *ClientConfig
	registered []any

	peer *Peer
}

func NewClient(config *ClientConfig) *Client {
	if config == nil {
		config = DefaultClientConfig()
	}
	return &Client{config: config}
}

type ClientConfig struct {
	Impl        rpc.Interface
	Logger      *log.Logger
	YamuxConfig *yamux.Config
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Impl:   yrpc.New(),
		Logger: log.New(os.Stdout, "[client] ", 0),
	}
}

func (c *Client) Dial(network, addr string) (*Peer, error) {
	peer, err := NewPeerWithDial(network, addr, c.config.YamuxConfig, c.config.Impl)
	if err != nil {
		return nil, err
	}
	c.peer = peer
	return c.peer, nil
}

func (c *Client) Register(method any) error {
	return c.config.Impl.Register(method)
}

func (c *Client) Close() error {
	if c.peer != nil {
		return c.peer.Close()
	}
	return nil
}

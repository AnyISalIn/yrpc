package bastion

import (
	"fmt"
	"github.com/AnyISalIn/yrpc"
	shared "github.com/AnyISalIn/yrpc/shared"
	"github.com/ugorji/go/codec"
	"io"
)

const (
	BASTION_PING           = "Bastion.Ping"
	BASTION_FORWARD_STREAM = "Bastion.ForwardStream"
)

type agent struct {
	name string
	peer *yrpc.Peer
}

type Bastion struct {
	agentMap map[string]*agent
}

func (b *Bastion) Ping(args *shared.Empty, reply *shared.Empty) error {
	return nil
}

type ForwardArgs struct {
	AgentID       string
	ServiceMethod string
}

func (b *Bastion) ForwardStream(rwc io.ReadWriteCloser) error {
	defer rwc.Close()

	args := new(ForwardArgs)
	decoder := codec.NewDecoder(rwc, &codec.MsgpackHandle{})

	if err := decoder.Decode(args); err != nil {
		return err
	}

	serverLogger.Printf("reply ack to client")

	agt, got := b.agentMap[args.AgentID]
	if !got {
		return fmt.Errorf("can't find agent id %s", args.AgentID)
	}

	return agt.ForwardStream(args.ServiceMethod, rwc)
}

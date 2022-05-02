package bastion

import (
	"encoding/gob"
	"fmt"
	"github.com/AnyISalIn/yrpc"
	shared "github.com/AnyISalIn/yrpc/shared"
	"io"
)

const (
	BASTION_PING = "Bastion.Ping"
	BASTION_EXEC = "Bastion.Exec"
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

// Forward Exec Client <-> Server <-> Agent
type BastionExecArgs struct {
	AgentID string
	Cmd     string
	TTY     bool
}

func (b *Bastion) Exec(rwc io.ReadWriteCloser) error {
	defer rwc.Close()
	args := new(BastionExecArgs)
	decoder := gob.NewDecoder(rwc)

	if err := decoder.Decode(args); err != nil {
		return err
	}

	agt, got := b.agentMap[args.AgentID]
	if !got {
		return fmt.Errorf("can't find agent id %s", args.AgentID)
	}

	return agt.PerformCmd(args.Cmd, args.TTY, rwc)
}

package bastion

import (
	shared "github.com/AnyISalIn/yrpc/shared"
	"github.com/creack/pty"
	"github.com/ugorji/go/codec"
	"io"
	"log"
	"os"
	"os/exec"
)

const (
	AGENT_WHOAMI = "Agent.Whoami"
	AGENT_ROLE   = "Agent.Role"
	AGENT_EXEC   = "Agent.Exec"
)

var agentLogger = log.New(os.Stdout, "[agent] ", shared.LogFlags)

type Agent struct{}

func (a *Agent) Whoami(args *Empty, reply *string) error {
	name, _ := os.Hostname()
	*reply = name
	return nil
}

func (a *Agent) Role(args *Empty, reply *string) error {
	*reply = "agent"
	return nil
}

type ExecArgs struct {
	Cmd string
	TTY bool
}

func (a *Agent) Exec(rwc io.ReadWriteCloser) error {
	defer rwc.Close()
	//encoder := gob.NewEncoder(rwc)
	agentLogger.Printf("[agent] handle stream exec, waiting exec args ...")
	decoder := codec.NewDecoder(rwc, &codec.MsgpackHandle{})

	args := new(ExecArgs)
	if err := decoder.Decode(args); err != nil {
		return err
	}

	agentLogger.Printf("[agent] received exec args %v", args)

	if args.TTY {
		cmd := exec.Command(args.Cmd)
		ptmx, err := pty.Start(cmd)
		if err != nil {
			agentLogger.Printf("[agent] failed to start tty %v", err)
			return err
		}

		defer ptmx.Close()

		Bridge(rwc, ptmx)

		return nil
	}
	cmd := exec.Command(args.Cmd)
	cmd.Stdout = rwc
	cmd.Stderr = rwc

	if err := cmd.Run(); err != nil {
		agentLogger.Printf("[agent] failed to run cmd %v", err)
	}
	return nil
}

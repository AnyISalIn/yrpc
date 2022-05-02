package bastion

import (
	"encoding/gob"
	"github.com/creack/pty"
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
	log.Printf("[agent] handle stream exec, waiting exec args ...")
	decoder := gob.NewDecoder(rwc)
	args := new(ExecArgs)
	if err := decoder.Decode(args); err != nil {
		return err
	}

	log.Printf("[agent] received exec args %v", args)

	if args.TTY {
		cmd := exec.Command(args.Cmd)
		ptmx, err := pty.Start(cmd)
		if err != nil {
			log.Printf("[agent] failed to start tty %v", err)
			return err
		}

		defer func() { _ = ptmx.Close() }() // Best effort.
		go func() { _, _ = io.Copy(ptmx, rwc) }()
		_, _ = io.Copy(rwc, ptmx)

		return nil
	}
	cmd := exec.Command(args.Cmd)
	cmd.Stdout = rwc
	cmd.Stderr = rwc

	if err := cmd.Run(); err != nil {
		log.Printf("[agent] failed to run cmd %v", err)
	}
	return nil
}

package main

import (
	"flag"
	"github.com/AnyISalIn/yrpc"
	"github.com/AnyISalIn/yrpc/examples/bastion"
	shared "github.com/AnyISalIn/yrpc/shared"
	"github.com/ugorji/go/codec"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var logger = log.New(os.Stdout, "[client] ", shared.LogFlags)

type Agent struct{}

func (a *Agent) Role(args *bastion.Empty, reply *string) error {
	*reply = "client"
	return nil
}

func main() {
	serverAddrPtr := flag.String("server", "127.0.0.1:8043", "bastion server address")
	agentId := flag.String("d", "MacBook-Pro-4.local", "managed agent id")
	tty := flag.Bool("t", false, "tty mode")

	flag.Parse()

	client := yrpc.NewClient(nil)
	defer client.Close()
	agt := new(Agent)
	if err := client.Register(agt); err != nil {
		logger.Fatal(err)
	}

	peer, err := client.Dial("tcp", *serverAddrPtr)
	if err != nil {
		logger.Fatal(err)
	}

	forwardArgs := bastion.ForwardArgs{
		AgentID:       *agentId,
		ServiceMethod: bastion.AGENT_EXEC,
	}

	execArgs := bastion.ExecArgs{
		Cmd: strings.Join(flag.Args(), " "),
		TTY: *tty,
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
		}
	}()

	conn, err := peer.Stream(bastion.BASTION_FORWARD_STREAM)
	if err != nil {
		logger.Panic(conn)
	}
	defer conn.Close()

	encoder := codec.NewEncoder(conn, &codec.MsgpackHandle{})
	if err := encoder.Encode(forwardArgs); err != nil {
		logger.Panic(err)
	}

	logger.Printf("sending forward request to bastion")

	//encoder = gob.NewEncoder(conn)
	if err := encoder.Encode(execArgs); err != nil {
		logger.Panic(err)
	}

	logger.Printf("sending exec request to agent")

	if *tty {
		bastion.Bridge(conn, os.Stdin)
		return
	}

	bastion.Bridge(conn, os.Stdout)
}

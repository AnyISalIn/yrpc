package main

import (
	"encoding/gob"
	"flag"
	"github.com/AnyISalIn/yrpc"
	"github.com/AnyISalIn/yrpc/examples/bastion"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

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
		log.Fatal(err)
	}

	peer, err := client.Dial("tcp", *serverAddrPtr)
	if err != nil {
		log.Fatal(err)
	}

	args := bastion.BastionExecArgs{
		AgentID: *agentId,
		Cmd:     strings.Join(flag.Args(), " "),
		TTY:     *tty,
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
		}
	}()

	conn, err := peer.Stream(bastion.BASTION_EXEC)
	if err != nil {
		log.Panic(conn)
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(args); err != nil {
		log.Panic(err)
	}

	if *tty {
		bastion.Bridge(conn, os.Stdin)
		return
	}

	bastion.Bridge(conn, os.Stdout)
}

package main

import (
	"flag"
	"github.com/AnyISalIn/yrpc"
	"github.com/AnyISalIn/yrpc/examples/bastion"
	shared "github.com/AnyISalIn/yrpc/shared"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "[agent] ", shared.LogFlags)

func main() {
	serverAddrPtr := flag.String("server", "127.0.0.1:8043", "bastion server address")
	flag.Parse()

	client := yrpc.NewClient(nil)
	defer client.Close()

	agent := new(bastion.Agent)
	if err := client.Register(agent); err != nil {
		logger.Fatalf("[agent] failed to register %v", err)
	}

	for {
		peer, err := client.Dial("tcp", *serverAddrPtr)
		if err != nil {
			logger.Printf("[agent] failed to dial %s, err %v", serverAddrPtr, err)
			goto RETRY
		}

		if err := peer.Call(bastion.BASTION_PING, &bastion.Empty{}, &bastion.Empty{}); err != nil {
			logger.Printf("[agent] failed to exec %s, err %v", bastion.BASTION_PING, err)
		} else {
			logger.Printf("[agent] bastion is ready")
		}

		select {}

	RETRY:
		logger.Printf("[agent] retry in 30s")
		time.Sleep(30 * time.Second)
	}
}

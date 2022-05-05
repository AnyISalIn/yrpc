package main

import (
	"flag"
	"github.com/AnyISalIn/yrpc/examples/bastion"
	shared "github.com/AnyISalIn/yrpc/shared"
	"log"
	"net"
	"os"
)

var logger = log.New(os.Stdout, "[server] ", shared.LogFlags)

func main() {

	serverAddrPtr := flag.String("server", "127.0.0.1:8043", "bastion server address")
	flag.Parse()

	srv := bastion.NewServer()
	listener, err := net.Listen("tcp", *serverAddrPtr)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("listening on %s", *serverAddrPtr)

	srv.Serve(listener)
}

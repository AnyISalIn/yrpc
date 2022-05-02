package main

import (
	"flag"
	"github.com/AnyISalIn/yrpc/examples/bastion"
	"log"
	"net"
)

func main() {

	serverAddrPtr := flag.String("server", "127.0.0.1:8043", "bastion server address")
	flag.Parse()

	srv := bastion.NewServer()
	listener, err := net.Listen("tcp", *serverAddrPtr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening on %s", *serverAddrPtr)

	srv.Serve(listener)
}

package iothub

import (
	"github.com/AnyISalIn/yrpc"
	shared "github.com/AnyISalIn/yrpc/shared"
	"log"
	"net"
)

func ServerLoop() {
	listener, err := net.Listen("tcp", "127.0.0.1:9084")
	if err != nil {
		log.Fatal(err)
	}

	srv := yrpc.NewServer(nil)
	if err := srv.Register(new(Hub)); err != nil {
		log.Fatalf("[server] failed to register %v", err)
	}

	for {
		peer, err := srv.AcceptPeer(listener)
		if err != nil {
			continue
		}
		go HandlePeer(peer)
	}
}

func HandlePeer(peer *yrpc.Peer) {
	var devCategory string
	var devId string

	if err := peer.Call(DEVICE_CATEGORY, &shared.Empty{}, &devCategory); err != nil {
		log.Fatalf("[server] failed to get category %v", err)
	}
	if err := peer.Call(DEVICE_ID, &shared.Empty{}, &devId); err != nil {
		log.Fatalf("[server] failed to get deviceId %v", err)
	}

	log.Printf("[server] Device %s/%s: connected", devCategory, devId)

	log.Printf("[server] Device %s/%s: let it off", devCategory, devId)
	if err := peer.Call(DEVICE_OFF, &shared.Empty{}, &shared.Empty{}); err != nil {
		log.Fatalf("[server] failed to let device off %v", err)
	}

	var state *bool
	if err := peer.Call(DEVICE_STATE, &shared.Empty{}, &state); err != nil {
		log.Fatalf("[server] failed to query device state %v", err)
	}
	log.Printf("[server] Device %s/%s: state %v", devCategory, devId, *state)
}

package iothub

import (
	shared "github.com/AnyISalIn/yrpc/shared"
	"log"
)

type Hub struct {
}

const (
	HUB_PING = "Hub.Ping"
)

func (h *Hub) Ping(deviceId *string, res *shared.Empty) error {
	if deviceId != nil {
		log.Printf("[server] received heartbeat from %s", *deviceId)
	}
	return nil
}

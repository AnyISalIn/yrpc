package iothub

import (
	"github.com/AnyISalIn/yrpc"
	shared "github.com/AnyISalIn/yrpc/shared"
	"log"
	"time"
)

func ClientLoop(devId string) {
	deviceId := devId

	client := yrpc.NewClient(nil)
	defer client.Close()

	if err := client.Register(&Device{id: deviceId}); err != nil {
		log.Fatalf("[client] failed to register %v", err)
	}

	for {
		peer, err := client.Dial("tcp", "127.0.0.1:9084")
		if err != nil {
			log.Printf("[client] failed to dial %s", err)
			goto RETRY
		}

		go func() {
			checkTicker := time.NewTicker(time.Second * 10)
			for {
				select {
				case <-checkTicker.C:
					if err := peer.Call(HUB_PING, &deviceId, &shared.Empty{}); err != nil {
						log.Printf("[client] failed to ping hub %v", err)
					} else {
						log.Printf("[client] ping hub successfuly")
					}
				}
			}
		}()

		// LOOPING
		select {}

	RETRY:
		timer := time.NewTimer(time.Second * 30)
		<-timer.C
		log.Println("[client] retry in 30s")
	}
}

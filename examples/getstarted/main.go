package main

import (
	"bufio"
	"fmt"
	"github.com/AnyISalIn/yrpc"
	"io"
	"log"
	"net"
	"time"
)

type ExampleServer struct{}

type ExampleClient struct{}

func (c *ExampleClient) Logs(rwc io.ReadWriteCloser) error {
	defer rwc.Close()

	var logChan = make(chan []byte, 10)
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				logChan <- []byte(fmt.Sprintf("%s -> new log\n", time.Now()))
			}
		}
	}()

	for {
		select {
		case msg := <-logChan:
			if _, err := rwc.Write(msg); err != nil {
				log.Printf("failed to write %v", err)
				return err
			}
		}
	}
}

var addr = "127.0.0.1:3439"

func serverLoop() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	server := yrpc.NewServer(nil)
	if err := server.Register(new(ExampleServer)); err != nil {
		log.Fatal(err)
	}

	for {
		peer, err := server.AcceptPeer(listener)
		if err != nil {
			log.Printf("failed to accept peer %v", err)
		}

		go func() {
			if rwc, err := peer.Stream("ExampleClient.Logs"); err != nil {
				log.Fatal(err)
			} else {
				defer rwc.Close()
				reader := bufio.NewReader(rwc)
				for {
					line, _, err := reader.ReadLine()
					if err != nil {
						log.Printf("read buffer error %v", err)
						return
					}

					log.Printf("ExampleClient.Logs -> %s", line)
				}
			}
		}()

	}

}
func clientLoop() {
	client := yrpc.NewClient(nil)
	defer client.Close()

	if err := client.Register(new(ExampleClient)); err != nil {
		log.Fatal(err)
	}

	for {
		_, err := client.Dial("tcp", addr)
		if err != nil {
			goto RETRY
		}

		select {}
	RETRY:
		log.Printf("retry in 30s")
		time.Sleep(time.Second * 30)
	}
}

func main() {
	go serverLoop()
	time.Sleep(time.Millisecond * 500)
	clientLoop()

}

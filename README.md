# YRPC - A simplify bidirectional Streaming RPC Framework

> based on [yamux](https://github.com/hashicorp/yamux) and net/rpc

## Examples:

- [bastion](https://github.com/AnyISalIn/yrpc/tree/main/examples/bastion)
- [iothub demo](https://github.com/AnyISalIn/yrpc/tree/main/examples/iothub)

## Quick Started

```go
package main

import (
	"github.com/AnyISalIn/yrpc"
	"log"
	"net"
	"time"
)

type ExampleServer struct{}

func (s *ExampleServer) Hello(args *string, reply *string) error {
	*reply = "Server: hello" + *args
	return nil
}

type ExampleClient struct{}

func (c *ExampleClient) Hello(args *string, reply *string) error {
	*reply = "Client: hello" + *args
	return nil
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
			var message = "world"
			var reply string
			if err := peer.Call("ExampleClient.Hello", &message, &reply); err != nil {
				log.Fatal(err)
			}
			log.Printf("ExampleClient.Hello -> %s", reply)
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
		var message = "world"
		var reply string

		peer, err := client.Dial("tcp", addr)
		if err != nil {
			goto RETRY
		}

		if err := peer.Call("ExampleServer.Hello", &message, &reply); err != nil {
			log.Fatal(err)
		}

		log.Printf("ExampleServer.Hello -> %s", reply)

		select {}
	RETRY:
		log.Printf("retry in 30s")
		time.Sleep(time.Second * 30)
	}
}

func main() {
	go serverLoop()
	clientLoop()

}
```

## Streaming RPC

```go
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
	clientLoop()

}

```

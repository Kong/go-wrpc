package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/kong/go-wrpc/examples/proto"
	"github.com/kong/go-wrpc/wrpc"
)

func main() {
	input := flag.String("input", "", "string to echo")
	flag.Parse()
	if *input == "" {
		flag.Usage()
		os.Exit(1)
		return
	}

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}
	c, _, err := wrpc.DefaultDialer.Dial(context.Background(), u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	peer := &wrpc.Peer{}
	peer.AddConn(c)

	err = peer.Register(&proto.EchoServer{})
	if err != nil {
		log.Fatalln(err)
	}

	client := &proto.EchoClient{Peer: peer}

	resp, err := client.Echo(context.Background(),
		&proto.EchoRPCRequest{S: *input})
	if err != nil {
		log.Fatal("echo:", err)
	}
	log.Println("response:", resp.S)
}

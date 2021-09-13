package main

import (
	"context"
	"log"
	"net/http"

	"github.com/kong/go-wrpc/examples/proto"
	"github.com/kong/go-wrpc/wrpc"
)

type echo struct{}

func (e echo) Echo(_ context.Context, req *proto.EchoRPCRequest) (
	*proto.EchoRPCResponse,
	error) {
	log.Println("request:", req.S)
	return &proto.EchoRPCResponse{
		S: "echo-" + req.S,
	}, nil
}

func main() {
	echoServer := &proto.EchoServer{
		Echo: echo{},
	}

	peer := &wrpc.Peer{}
	err := peer.Register(echoServer)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		err := peer.Upgrade(w, r)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

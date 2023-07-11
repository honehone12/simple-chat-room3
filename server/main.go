package main

import (
	"fmt"
	"log"
	"net"
	"simple-chat-room3/common"
	pb "simple-chat-room3/pb"

	"google.golang.org/grpc"
)

func launchLobyServer(errCh chan<- error) {
	grpcServer := grpc.NewServer()
	lServer := LobyServer{}
	pb.RegisterLobyServiceServer(grpcServer, lServer)

	addr := fmt.Sprintf(common.Localhost, common.Port-1)

	l, err := net.Listen(common.Transport, addr)
	if err != nil {
		errCh <- err
		return
	}

	log.Printf("starting loby server at %s\n", addr)
	errCh <- grpcServer.Serve(l)
}

func launchChatRoomServer(i int, errCh chan<- error) {
	grpcServer := grpc.NewServer()
	crServer := ChatRoomServer{
		msgMemMap:  make(map[string]*MsgMemory),
		sortedKeys: make([]string, 0),
	}
	pb.RegisterChatRoomServiceServer(grpcServer, &crServer)

	addr := fmt.Sprintf(common.Localhost, common.Port+i)

	l, err := net.Listen(common.Transport, addr)
	if err != nil {
		errCh <- err
		return
	}

	log.Printf("starting chat room server[%d] at %s\n", i, addr)
	errCh <- grpcServer.Serve(l)
}

func main() {
	errChan := make(chan error)

	go launchLobyServer(errChan)

	for i := 0; i < common.NumServer; i++ {
		go launchChatRoomServer(i, errChan)
	}

	log.Fatalln(<-errChan)
}

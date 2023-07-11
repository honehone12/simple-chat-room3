package main

import (
	"errors"
	"fmt"
	"simple-chat-room2/common"
	pb "simple-chat-room2/pb"
)

type Display struct {
	errCh chan error
}

func NewDisplay() Display {
	return Display{
		errCh: make(chan error),
	}
}

func (d Display) ErrChan() <-chan error {
	return d.errCh
}

func displayMessage(player string, msg string) {
	fmt.Printf("[%s] ", player)
	fmt.Printf("%s\n", msg)
}

func reverseLines(lns int) {
	for i := 0; i < lns; i++ {
		fmt.Printf("\r\033[1A")
	}
}

func (d Display) Display(stream pb.ChatRoomService_ChatClient) {
	for {
		bundle, err := stream.Recv()
		if err != nil {
			d.errCh <- err
			break
		}
		if !bundle.GetOk() {
			d.errCh <- errors.New(bundle.GetErrMsg().GetMsg())
			break
		}

		msgs := bundle.GetChatMsgs()
		len := len(msgs)
		fmt.Println(common.Space64)
		for i := 0; i < len; i++ {
			m := msgs[i]
			displayMessage(m.GetName(), m.GetMsg())
		}
		fmt.Println(common.Space64)
		reverseLines(len + 2)
	}
}

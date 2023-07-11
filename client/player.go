package main

import (
	"context"
	"errors"
	"fmt"
	"simple-chat-room3/common"
	pb "simple-chat-room3/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func displayMessage(player string, msg string) {
	fmt.Printf("[%s] ", player)
	fmt.Printf("%s\n", msg)
}

func reverseLines(lns int) {
	for i := 0; i < lns; i++ {
		fmt.Printf("\r\033[1A")
	}
}

type PlayerAgent struct {
	name           string
	consumingLines int
	input          KeyInput
	errCh          chan error
}

func NewPlayerAgent(name string) *PlayerAgent {
	return &PlayerAgent{
		name:           name,
		consumingLines: 0,
		input:          NewKeyInput(),
		errCh:          make(chan error),
	}
}

func (p *PlayerAgent) cleanUpDisplay() {
	for i := 0; i < p.consumingLines; i++ {
		fmt.Println(common.Space64)
	}
	reverseLines(p.consumingLines)
}

func (p *PlayerAgent) StartKeyInput() {
	go p.input.Input()
}

func (p *PlayerAgent) CatchError() error {
	inputE := p.input.ErrChan()
	var err error
	select {
	case err = <-inputE:
	case err = <-p.errCh:
	}
	return err
}

// cache conns and clients ??

func (p *PlayerAgent) JoinChatRoom(room pb.Room) {
	addr, err := p.askLoby(room)
	if err != nil {
		p.errCh <- err
		return
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		p.errCh <- err
		return
	}
	defer conn.Close()
	crClient := pb.NewChatRoomServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := crClient.Join(ctx, &pb.JoinRequest{Name: p.name})
	if err != nil {
		p.errCh <- err
		return
	}
	if res.GetErrMsg().GetErr() {
		p.errCh <- errors.New(res.GetErrMsg().GetMsg())
		return
	}

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	stream, err := crClient.Chat(ctx)
	if err != nil {
		p.errCh <- err
		return
	}

	err = p.sync(stream)
	if err != nil {
		p.errCh <- err
	}
}

func (p *PlayerAgent) askLoby(room pb.Room) (string, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(common.Localhost, common.Port-1),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	lobyClient := pb.NewLobyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := lobyClient.Room(ctx, &pb.RoomRequest{Room: room})
	if err != nil {
		return "", err
	}
	if res.GetErrMsg().GetErr() {
		return "", errors.New(res.GetErrMsg().GetMsg())
	}

	return res.GetIpAddr(), nil
}

func (p *PlayerAgent) sync(stream pb.ChatRoomService_ChatClient) error {
	ticker := time.NewTicker(time.Millisecond * common.InputSyncMil)
	roomChangeCh := p.input.RoomChangeChan()
	defer ticker.Stop()

	for {
		select {
		case room := <-roomChangeCh:
			go p.JoinChatRoom(room)
			p.cleanUpDisplay()
			return nil
		case now := <-ticker.C:
			var s string
			s, err := p.input.buffer.String()
			if err != nil {
				return err
			}

			err = stream.Send(&pb.ChatClientMsg{
				UnixMil: now.UnixMilli(),
				ChatMsg: &pb.ChatMsg{
					Name: p.name,
					Msg:  s,
				},
			})
			if err != nil {
				return err
			}

			bundle, err := stream.Recv()
			if err != nil {
				return err
			}
			if bundle.GetErrMsg().GetErr() {
				return errors.New(bundle.GetErrMsg().GetMsg())
			}

			msgs := bundle.GetChatMsgs()
			len := len(msgs)

			// gnome's default size only
			// display will be broken in other terminal sizes
			fmt.Println(common.Space64)
			for i := 0; i < len; i++ {
				m := msgs[i]
				displayMessage(m.GetName(), m.GetMsg())
			}
			fmt.Println(common.Space64)
			p.consumingLines = len + 2
			reverseLines(p.consumingLines)
		}
	}
}

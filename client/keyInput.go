package main

import (
	"simple-chat-room2/common"
	pb "simple-chat-room2/pb"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	keyBufferSize = 10
	nonAlphaNum   = 0x0
	space         = 0x20
)

type KeyInput struct {
	buffer *Buffer

	errCh chan error
}

func NewKeyInput() KeyInput {
	return KeyInput{
		buffer: NewBuffer(),
		errCh:  make(chan error),
	}
}

func (i KeyInput) ErrChan() <-chan error {
	return i.errCh
}

func (i KeyInput) Input() {
	keyEvents, err := keyboard.GetKeys(keyBufferSize)
	if err != nil {
		i.errCh <- err
		return
	}

	for e := range keyEvents {
		if e.Err != nil {
			err = e.Err
			break
		}

		if e.Key == keyboard.KeyEsc {
			err = nil
			break
		} else if e.Key == keyboard.KeyBackspace || e.Key == keyboard.KeyBackspace2 {
			i.buffer.Back()
		} else if e.Key == keyboard.KeySpace {
			i.buffer.Add(space)
		} else if e.Rune != nonAlphaNum {
			// still not sure what will be lost with cast
			i.buffer.Add(byte(e.Rune))
		}
	}

	// can not handle this error (this case needs force quit anyway)
	_ = keyboard.Close()
	// wait for close then send err to main, to prevent quit app before close
	i.errCh <- err
}

func (i KeyInput) Sync(name string, stream pb.ChatRoomService_ChatClient) {
	ticker := time.NewTicker(time.Millisecond * common.InputSyncMil)

	var err error
	for now := range ticker.C {
		var s string
		s, err = i.buffer.String()
		if err != nil {
			break
		}

		err = stream.Send(&pb.ChatClientMsg{
			UnixMil: now.UnixMilli(),
			ChatMsg: &pb.ChatMsg{
				Name: name,
				Msg:  s,
			},
		})
		if err != nil {
			break
		}
	}

	ticker.Stop()
	i.errCh <- err
}

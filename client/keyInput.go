package main

import (
	pb "simple-chat-room3/pb"

	"github.com/eiannone/keyboard"
)

const (
	keyBufferSize = 10
	nonAlphaNum   = 0x0
	space         = 0x20
)

func CloseKeyInput() error {
	return keyboard.Close()
}

type KeyInput struct {
	buffer *Buffer

	roomChangeCh chan pb.Room
	errCh        chan error
}

func NewKeyInput() KeyInput {
	return KeyInput{
		buffer: NewBuffer(),

		roomChangeCh: make(chan pb.Room),
		errCh:        make(chan error),
	}
}

func (i KeyInput) RoomChangeChan() <-chan pb.Room {
	return i.roomChangeCh
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
			i.errCh <- err
			break
		}

		if e.Key == keyboard.KeyEsc {
			i.errCh <- nil
			break
		} else if e.Key == keyboard.KeyBackspace || e.Key == keyboard.KeyBackspace2 {
			i.buffer.Back()
		} else if e.Key == keyboard.KeySpace {
			i.buffer.Add(space)
		} else if e.Key >= keyboard.KeyF12 && e.Key <= keyboard.KeyF1 {
			i.buffer.Clear()
			room := pb.Room(keyboard.KeyF1 - e.Key)
			i.roomChangeCh <- room
		} else if e.Rune != nonAlphaNum {
			// en only anyway
			i.buffer.Add(byte(e.Rune))
		}
	}
}

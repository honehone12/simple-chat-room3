package main

import (
	"errors"
	"simple-chat-room2/common"
	"unicode/utf8"
)

var (
	ErrInvalidAsUtf8 = errors.New("bytes are invalid as utf8")
)

type Buffer struct {
	inner  [common.InputBufferSize]byte
	cursor byte
}

func NewBuffer() *Buffer {
	b := Buffer{
		inner:  [common.InputBufferSize]byte{},
		cursor: 0,
	}
	for i := 0; i < common.InputBufferSize; i++ {
		b.inner[i] = space
	}
	return &b
}

func (b *Buffer) Add(value byte) {
	if b.cursor < common.InputBufferSize {
		b.inner[b.cursor] = value
		b.cursor++
	}
}

func (b *Buffer) Back() {
	if b.cursor > 0 {
		b.cursor--
		b.inner[b.cursor] = space
	}
}

func (b *Buffer) String() (string, error) {
	if utf8.Valid(b.inner[:]) {
		return string(b.inner[:]), nil
	} else {
		return "", ErrInvalidAsUtf8
	}
}

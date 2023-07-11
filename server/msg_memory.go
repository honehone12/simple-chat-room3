package main

import (
	"simple-chat-room3/common"
)

type MsgMemory struct {
	timeStamp int64
	msg       string
}

func NewMsgMemory(timeStamp int64) *MsgMemory {
	return &MsgMemory{
		timeStamp: timeStamp,
		msg:       common.Space64,
	}
}

func (m *MsgMemory) Set(timeStamp int64, msg string) {
	if timeStamp > m.timeStamp {
		m.timeStamp = timeStamp
		m.msg = msg
	}
}

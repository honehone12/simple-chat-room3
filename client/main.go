package main

import (
	"flag"
	"log"
	pb "simple-chat-room3/pb"
)

func main() {
	playerName := flag.String("name", "", "player's name")
	flag.Parse()

	if *playerName == "" {
		log.Fatalln("name is needed")
	}

	player := NewPlayerAgent(*playerName)
	player.StartKeyInput()
	defer CloseKeyInput()

	go player.JoinChatRoom(pb.Room_F1)

	err := player.CatchError()
	player.cleanUpDisplay()
	if err != nil {
		log.Fatalln(err)
	}
}

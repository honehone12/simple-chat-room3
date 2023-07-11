package main

import (
	"context"
	"fmt"
	"simple-chat-room3/common"
	pb "simple-chat-room3/pb"
)

type LobyServer struct {
	pb.UnimplementedLobyServiceServer
}

func Room(ctx context.Context, req *pb.RoomRequest) (*pb.RoomResponse, error) {
	r := req.GetRoom()
	res := &pb.RoomResponse{
		IpAddr: "",
		ErrMsg: nil,
	}

	if r >= pb.Room_F1 && r <= pb.Room_F12 {
		res.IpAddr = fmt.Sprintf(common.Localhost, common.Port+r)
	} else {
		res.ErrMsg = &pb.ErrorMsg{
			Err: true,
			Msg: "unknown room",
		}
	}

	return res, nil
}

.PHONY:
protoc:
	protoc ./pb/chat_room.proto --go_out=./pb --go-grpc_out=./pb

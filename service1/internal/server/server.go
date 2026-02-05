package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "service1/internal/gen/message"
	"service1/internal/service"
)

func Init() {
	lis, err := net.Listen("tcp",":50051")
	if err!=nil{
		log.Fatalf("failed to listen :%v",err)
	}

	grpcServer:=grpc.NewServer()

	reflection.Register(grpcServer)

	msgService:=service.NewMessageService()

	pb.RegisterMessageServiceServer(grpcServer,msgService)

	fmt.Println("running on port 50051")

	if err:=grpcServer.Serve(lis);err!=nil{
		log.Fatalf("error in serving grpc server :%v",err)
	}
}
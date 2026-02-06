package service

import (
	"context"
	"encoding/json"
	"service1/internal/database"
	pb "service1/internal/gen/message"
	"service1/internal/models"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type MessageService struct{
	pb.UnimplementedMessageServiceServer
	respoitory *database.MessageRepository
}

func NewMessageService(repo *database.MessageRepository) *MessageService{
	return &MessageService{respoitory: repo}
}

func(s *MessageService) AddNumber(ctx context.Context, req *pb.Request) (*pb.Response,error){
	if req.Number_1<0 || req.Number_2<0{
		return nil,status.Error(codes.InvalidArgument,"both numbers can't be zero")
	}

	result :=req.Number_1+req.Number_2

	requestID:=uuid.New().String()

	messagePayload:=models.MessagePayLoad{
		RequestID:requestID,
		Num_1: req.Number_1,
		Nums_2: req.Number_2,
		Result: result,

	}
	payload,err:=json.Marshal(messagePayload)

	if err!=nil{
		return nil,status.Errorf(codes.Internal, "failed to serialize kafka payload: %v", err)
	}

	// Write in both the sum table and event table using repo
	if err:=s.respoitory.CreateMessageWithOutbox(ctx,messagePayload,payload);err!=nil{
		return nil,status.Errorf(codes.Internal,"error in database insterting :%v",err)
	}


	// Kafka is no longer part int the flow here in the outbox pattern we will publish then using worker
	// if err:=s.producer.Publish(ctx,payload);err!=nil{
	// 	return &pb.Response{Result: result},status.Errorf(codes.Internal, "failed to publish to Kafka: %v", err)
	// }


	return &pb.Response{Result: req.Number_1+req.Number_2},nil
}
package service

import (
	"context"
	"encoding/json"
	pb "service1/internal/gen/message"
	producer "service1/internal/kafka-producer"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type MessageService struct{
	pb.UnimplementedMessageServiceServer
	producer producer.Producer
}

func NewMessageService(p producer.Producer) *MessageService{
	return &MessageService{producer: p}
}

func(s *MessageService) AddNumber(ctx context.Context, req *pb.Request) (*pb.Response,error){
	if req.Number_1<0 || req.Number_2<0{
		return nil,status.Error(codes.InvalidArgument,"both numbers can't be zero")
	}

	result :=req.Number_1+req.Number_2

	requestID:=uuid.New().String()

	payload,err:=json.Marshal(map[string]any{
		"request_id":requestID,
		"number_1":req.Number_1,
		"number_2":req.Number_2,
		"result":result,
	})

	if err!=nil{
		return nil,status.Errorf(codes.Internal, "failed to serialize kafka payload: %v", err)
	}

	if err:=s.producer.Publish(ctx,payload);err!=nil{
		return &pb.Response{Result: result},status.Errorf(codes.Internal, "failed to publish to Kafka: %v", err)
	}
	return &pb.Response{Result: req.Number_1+req.Number_2},nil
}
package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"service1/internal/database"
	pb "service1/internal/gen/message"
	producer "service1/internal/kafka-producer"
	"service1/internal/service"
	"service1/internal/worker"
)

func Init() {

	ctx,cancel:=context.WithCancel(context.Background())

	defer cancel()

	KafkaProducer,err:=producer.NewKafkaProducer(producer.Config{
		Brokers: []string{"localhost:9092"},
		Topic: "message.created",
		ClientID: "service1",
	})

	if err!=nil{
		log.Fatalf("failed to create Kafka producer: %v", err)
	}


	defer KafkaProducer.Close()


	lis, err := net.Listen("tcp",":50051")
	if err!=nil{
		log.Fatalf("failed to listen :%v",err)
	}

	conn,err:=database.Open()
	if err!=nil{
		log.Fatalf("Error connection in the database: %v",err)
	}

	repo:=database.NewMessageRepository(conn)

	// Start my background worker run whenever the server boots end when server is shut down
	outboxWorker:=worker.NewOutboxWriter(KafkaProducer,repo)
	go outboxWorker.StreamOverKafka(ctx)


	grpcServer:=grpc.NewServer()

	reflection.Register(grpcServer)

	msgService:=service.NewMessageService(repo)

	pb.RegisterMessageServiceServer(grpcServer,msgService)

	log.Println("running on port 50051")

	if err:=grpcServer.Serve(lis);err!=nil{
		log.Fatalf("error in serving grpc server :%v",err)
	}
}
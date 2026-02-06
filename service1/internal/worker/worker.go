package worker

import (
	"context"
	"log"
	"service1/internal/database"
	producer "service1/internal/kafka-producer"
	"time"
)

type OutboxWorker struct {
	Producer producer.Producer
	Repository *database.MessageRepository
}

func NewOutboxWriter(p producer.Producer,repo *database.MessageRepository)*OutboxWorker{
	return &OutboxWorker{Producer: p,Repository: repo}
}


func (w *OutboxWorker)StreamOverKafka(ctx context.Context){

	ticker:=time.NewTicker(2*time.Second)
	defer ticker.Stop()

	for {
		select{
		case <- ctx.Done():
			log.Println("outbox worker stopped")
			return
		case <- ticker.C:
			events,err:=w.Repository.GetAllUnProcessedMessages(ctx,10)
			if err!=nil{
				log.Fatalf("error :%v",err)
				continue
			}

		for _,event :=range events{
			if err:=w.Producer.Publish(ctx,event.Payload);err!=nil{
				log.Fatalf("error in publishing events through :%v",err)
			}
			if err:=w.Repository.UpdateMessageProcessingStatus(ctx,event.ID);err!=nil{
				log.Fatalf("error in updating the message status :%v",err)
			}
		}		

		}
	}
}
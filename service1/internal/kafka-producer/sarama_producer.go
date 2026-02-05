package producer

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type KafkaProducer struct{
	producer sarama.AsyncProducer
	topic string
}

func NewKafkaProducer (cfg Config) (*KafkaProducer,error){
	config:=sarama.NewConfig()

	config.ClientID=cfg.ClientID
	config.Producer.Return.Errors=true
	config.Producer.Return.Successes=true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Metadata.AllowAutoTopicCreation = true

	asyncProducer,err:=sarama.NewAsyncProducer(cfg.Brokers,config)
	if err!=nil{
		return nil,err
	}

	kafkaProducer:=&KafkaProducer{
		producer: asyncProducer,
		topic: cfg.Topic,
	}

	go kafkaProducer.handleErrors()
	go kafkaProducer.handleSuccess()

	return kafkaProducer,nil
}


/*
	Sarama.AsyncProducer works differently than a normal function call:

	It has internal channels for communication:

	Input() → you send messages here to be published

	Successes() → delivered messages come here

	Errors() → failed messages come here

*/
func(kp *KafkaProducer)Publish(ctx context.Context,value []byte)error{
	msg:=&sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(value),
	}

	select{
	case kp.producer.Input() <- msg: // message queued successfully
		fmt.Println("Message is queued successfully")
		return nil
	case <- ctx.Done():    // context(request) canceled 
		return ctx.Err()
	}
}

func (kp *KafkaProducer)Close()error{
	return kp.producer.Close()
}



// Adding metric and logs
func(kp *KafkaProducer)handleSuccess(){
	for msg:=range kp.producer.Successes(){
		_=msg
	}
}

func(kp *KafkaProducer)handleErrors(){
	for err:=range kp.producer.Errors(){
		log.Printf("kafka publish error: %v", err)
	}
}
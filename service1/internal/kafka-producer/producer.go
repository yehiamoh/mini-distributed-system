package producer

import "context"

type Producer interface {
	Publish(ctx context.Context,value []byte) error
	Close()error
}
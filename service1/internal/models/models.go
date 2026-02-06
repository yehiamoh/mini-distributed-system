package models

import "time"

type MessagePayLoad struct {
	RequestID string
	Num_1     int32
	Nums_2    int32
	Result    int32
	CreatedAt time.Time
}
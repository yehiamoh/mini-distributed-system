package models

import "time"

type OutboxEvent struct {
	ID         int
	EventType  string
	Payload    []byte // pgx can scan JSONB into []byte
	Processed  bool
	RetryCount int
	CreatedAt  time.Time
}
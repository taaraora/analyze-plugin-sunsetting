package storage

import "context"

type Interface interface {
	GetAll(ctx context.Context, prefix string) ([]Message, error)
	Get(ctx context.Context, prefix string, key string) (Message, error)
	Put(ctx context.Context, prefix string, key string, value Message) error
	Delete(ctx context.Context, prefix string, key string) error
	WatchPrefix(ctx context.Context, prefix string) <-chan WatchEvent
	Close() error
}

type Message interface {
	Payload() []byte
}

type WatchEvent interface {
	Message
	Type() WatchEventType
	Err() error
}

type WatchEventType string

const (
	Added    WatchEventType = "ADDED"
	Modified WatchEventType = "MODIFIED"
	Deleted  WatchEventType = "DELETED"
	Error    WatchEventType = "ERROR"
	Unknown  WatchEventType = "UNKNOWN"
)

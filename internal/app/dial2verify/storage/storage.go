package storage

import "context"

type Storage interface {
	CheckPhone(ctx context.Context, phone string) (bool, error)
	Close() error
}

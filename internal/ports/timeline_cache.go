package ports

import "context"

type TimelineCache interface {
	Get(ctx context.Context, userID string) ([]byte, error)
	Set(ctx context.Context, userID string, data []byte) error
	Delete(ctx context.Context, userID string) error
}

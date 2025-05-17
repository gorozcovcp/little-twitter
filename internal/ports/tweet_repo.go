package ports

import (
	"context"
	"time"

	"github.com/gorozcovcp/little-twitter/internal/domain/model"
)

type TweetRepository interface {
	Save(ctx context.Context, tweet *model.Tweet) error
	GetByUsersBefore(ctx context.Context, userIDs []string, before time.Time, limit int) ([]model.Tweet, error)
}

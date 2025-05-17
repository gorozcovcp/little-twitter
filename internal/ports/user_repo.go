package ports

import (
	"context"

	"github.com/gorozcovcp/little-twitter/internal/domain/model"
)

type UserRepository interface {
	Follow(ctx context.Context, userID, followID string) error
	GetByID(ctx context.Context, userID string) (*model.User, error)
	GetFollowersOf(ctx context.Context, userID string) ([]model.User, error)
}

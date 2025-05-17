package service

import (
	"context"

	"github.com/gorozcovcp/little-twitter/internal/ports"
)

type UserService struct {
	UserRepo      ports.UserRepository
	TimelineCache ports.TimelineCache
}

func NewUserService(ur ports.UserRepository, tc ports.TimelineCache) *UserService {
	return &UserService{
		UserRepo:      ur,
		TimelineCache: tc,
	}
}

func (s *UserService) Follow(ctx context.Context, userID, followID string) error {
	err := s.UserRepo.Follow(ctx, userID, followID)
	if err != nil {
		return err
	}

	// Invalida el cache del usuario que sigue a otro
	s.TimelineCache.Delete(ctx, userID)
	return nil
}

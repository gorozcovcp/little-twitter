package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorozcovcp/little-twitter/internal/domain/model"
	"github.com/gorozcovcp/little-twitter/internal/ports"
)

type TweetService struct {
	TweetRepo     ports.TweetRepository
	UserRepo      ports.UserRepository
	TimelineCache ports.TimelineCache
}

type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var ErrTweetTooLong = &AppError{"Tweet exceeds 280 characters"}

func NewTweetService(tr ports.TweetRepository, ur ports.UserRepository, tc ports.TimelineCache) *TweetService {
	return &TweetService{
		TweetRepo:     tr,
		UserRepo:      ur,
		TimelineCache: tc,
	}
}

func (s *TweetService) PostTweet(ctx context.Context, userID, content string) error {
	if len(content) > 280 {
		return ErrTweetTooLong
	}
	tweet := &model.Tweet{
		UserID:  userID,
		Content: content,
		Created: time.Now(),
	}

	err := s.TweetRepo.Save(ctx, tweet)
	if err != nil {
		return err
	}

	// Invalidar timeline cache de los seguidores del usuario
	followers, err := s.UserRepo.GetFollowersOf(ctx, userID)
	if err == nil {
		for _, follower := range followers {
			s.TimelineCache.Delete(ctx, follower.ID)
		}
	}

	return nil
}

func (s *TweetService) GetTimeline(ctx context.Context, userID string, before time.Time, limit int) ([]model.Tweet, error) {
	cached, err := s.TimelineCache.Get(ctx, userID)
	if err == nil && cached != nil {
		fmt.Printf("GetTimeline-----------------cached: %s\n", string(cached))
		var tweets []model.Tweet
		if jsonErr := json.Unmarshal(cached, &tweets); jsonErr == nil {
			return tweets, nil
		}
	}

	user, err := s.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tweets, err := s.TweetRepo.GetByUsersBefore(ctx, user.Follows, before, limit)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(tweets)
	s.TimelineCache.Set(ctx, userID, data)

	return tweets, nil
}

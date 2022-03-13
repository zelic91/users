package leaderboard

import (
	"context"
	"errors"
	"zelic91/users/gen/models"
	"zelic91/users/shared"
)

const (
	ErrInvalidScore = "invalid score"
)

type Service struct {
	Repo Repo
}

func NewService(repo Repo) Service {
	return Service{
		Repo: repo,
	}
}

func (s Service) GetLeaderboard(ctx context.Context, user *shared.UserClaims, limit *int32, offset *int32) (*models.Leaderboard, error) {
	items, err := s.Repo.GetLeaderboard(user.App, limit, offset)

	if err != nil {
		return nil, err
	}

	result := models.Leaderboard{}

	for _, item := range items {
		result = append(result, &models.LeaderboardItem{
			Username: item.Username,
			Score:    item.Score,
			Rank:     item.Rank,
		})
	}

	return &result, nil
}

func (s Service) SubmitScore(ctx context.Context, user *shared.UserClaims, params *models.ScoreRequest) error {
	userID := user.ID
	score := params.Score
	time := params.Time
	h := params.H

	if !shared.ValidateScoreRequest(*score, *time, *h) {
		return errors.New(ErrInvalidScore)
	}

	return s.Repo.SubmitScore(userID, *score)
}

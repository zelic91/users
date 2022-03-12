package leaderboard

import (
	"context"
	"zelic91/users/gen/models"
	"zelic91/users/gen/restapi/operations"
	"zelic91/users/gen/restapi/operations/leaderboards"
	"zelic91/users/gen/restapi/operations/scores"
	"zelic91/users/shared"

	"github.com/go-openapi/runtime/middleware"
)

type LeaderboardService interface {
	GetLeaderboard(ctx context.Context, limit *int32, offset *int32) (*models.Leaderboard, error)
	SubmitScore(ctx context.Context, user *shared.UserClaims, params *models.ScoreRequest) error
}

func SetupLeaderboard(api *operations.UsersAPI, service LeaderboardService) {
	api.LeaderboardsGetLeaderboardHandler = leaderboards.GetLeaderboardHandlerFunc(func(params leaderboards.GetLeaderboardParams, uc *shared.UserClaims) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		limit := params.Limit
		offset := params.Offset

		leaderboard, err := service.GetLeaderboard(ctx, limit, offset)

		if err != nil {
			return shared.HandleError(err)
		}

		return leaderboards.NewGetLeaderboardOK().WithPayload(*leaderboard)
	})

	api.ScoresSubmitScoreHandler = scores.SubmitScoreHandlerFunc(func(params scores.SubmitScoreParams, user *shared.UserClaims) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		err := service.SubmitScore(ctx, user, params.Body)
		if err != nil {
			return shared.HandleError(err)
		}

		return scores.NewSubmitScoreOK()
	})
}

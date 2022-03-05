package swagger

import (
	"context"
	"zelic91/users/gen/models"
	"zelic91/users/gen/restapi/operations"
	"zelic91/users/gen/restapi/operations/authentication"

	"github.com/go-openapi/runtime/middleware"
)

type AuthService interface {
	SignUp(ctx context.Context, params *models.SignUpRequest) (*models.AuthResponse, error)
	SignIn(ctx context.Context, params *models.SignInRequest) (*models.AuthResponse, error)
}

func SetupAuth(api *operations.UsersAPI, authService AuthService) {
	api.AuthenticationSignUpHandler = authentication.SignUpHandlerFunc(func(params authentication.SignUpParams) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		resp, err := authService.SignUp(ctx, params.Body)

		if err != nil {
			return handleError(err)
		}

		return authentication.NewSignUpOK().WithPayload(resp)
	})

	api.AuthenticationSignInHandler = authentication.SignInHandlerFunc(func(params authentication.SignInParams) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		resp, err := authService.SignIn(ctx, params.Body)

		if err != nil {
			return handleError(err)
		}

		return authentication.NewSignInOK().WithPayload(resp)
	})
}

package swagger

import (
	"zelic91/users/auth"
	"zelic91/users/gen/models"
	"zelic91/users/gen/restapi/operations"
	"zelic91/users/gen/restapi/operations/profile"

	"github.com/go-openapi/runtime/middleware"
)

func SetupProfile(api *operations.UsersAPI) {
	api.ProfileGetProfileHandler = profile.GetProfileHandlerFunc(func(gpp profile.GetProfileParams, u *auth.UserClaims) middleware.Responder {
		resp := models.Profile{
			Email:    "thuongnh.uit@gmail.com",
			Username: "Thuong Dep Trai",
		}
		return profile.NewGetProfileOK().WithPayload(&resp)
	})
}

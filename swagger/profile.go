package swagger

import (
	"zelic91/users/gen/models"
	"zelic91/users/gen/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func Profile(api *operations.UsersAPI) {
	api.GetProfileHandler = operations.GetProfileHandlerFunc(func(gpp operations.GetProfileParams, u *models.User) middleware.Responder {
		profile := models.Profile{
			Email:    "thuongnh.uit@gmail.com",
			Username: "Thuong Dep Trai",
		}
		return operations.NewGetProfileOK().WithPayload(&profile)
	})
}

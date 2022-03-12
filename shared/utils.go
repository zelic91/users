package shared

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *ErrorResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(http.StatusBadRequest)
	if err := producer.Produce(rw, r); err != nil {
		panic(err)
	}
}

func HandleError(err error) middleware.Responder {
	return &ErrorResponse{
		Status:  500,
		Message: err.Error(),
	}
}

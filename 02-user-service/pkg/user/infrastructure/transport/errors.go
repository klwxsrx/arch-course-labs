package transport

import (
	"errors"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/user/app"
	"net/http"
)

var errorsToHTTPCode = map[error]int{
	app.ErrUserAlreadyExists:                   http.StatusConflict,
	app.ErrUserWithSpecifiedEmailAlreadyExists: http.StatusConflict,
	app.ErrUserNotFound:                        http.StatusNotFound,
	errEmptyJSONBody:                           http.StatusBadRequest,
	errInvalidParameter:                        http.StatusBadRequest,
}

func translateError(err error) (errorCode int) {
	if err == nil {
		return http.StatusOK
	}
	for e, code := range errorsToHTTPCode {
		if errors.Is(err, e) {
			return code
		}
	}
	return http.StatusInternalServerError
}

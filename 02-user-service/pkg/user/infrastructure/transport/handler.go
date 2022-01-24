package transport

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/common/infrastructure/transport"
	"github.com/klwxsrx/arch-course-labs/user-service/pkg/user/app"
	"io"
	"net/http"
)

var (
	errInvalidParameter = errors.New("invalid parameter")
	errEmptyJSONBody    = errors.New("empty json body")
)

type requestHandlerFunc func(request *http.Request, service app.UserService) (result interface{}, err error)

type route struct {
	Name    string
	Method  string
	Pattern string
	Handler requestHandlerFunc
}

func getRoutes() []route {
	return []route{
		{
			"createUser",
			http.MethodPost,
			"/user",
			createUserHandler,
		},
		{
			"getUser",
			http.MethodGet,
			"/user/{userID}",
			getUserHandler,
		},
		{
			"updateUser",
			http.MethodPut,
			"/user/{userID}",
			updateUserHandler,
		},
		{
			"deleteUser",
			http.MethodDelete,
			"/user/{userID}",
			deleteUserHandler,
		},
		{
			"health",
			http.MethodGet,
			"/healthz",
			livenessCheckHandler,
		},
	}
}

type jsonUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

type jsonUserData struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func createUserHandler(request *http.Request, service app.UserService) (result interface{}, err error) {
	var user jsonUser
	err = parseJSONFromBody(request, &user)
	if err != nil {
		return nil, err
	}

	return nil, service.Create(user.ID, user.Email, user.FirstName, user.LastName)
}

func getUserHandler(request *http.Request, service app.UserService) (result interface{}, err error) {
	id, err := parseUUID(mux.Vars(request)["userID"])
	if err != nil {
		return nil, err
	}

	user, err := service.Get(id)
	if err != nil {
		return nil, err
	}

	return &jsonUser{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func updateUserHandler(request *http.Request, service app.UserService) (result interface{}, err error) {
	id, err := parseUUID(mux.Vars(request)["userID"])
	if err != nil {
		return nil, err
	}

	var user jsonUserData
	err = parseJSONFromBody(request, &user)
	if err != nil {
		return nil, err
	}

	return nil, service.Update(id, user.Email, user.FirstName, user.LastName)
}

func deleteUserHandler(request *http.Request, service app.UserService) (result interface{}, err error) {
	id, err := parseUUID(mux.Vars(request)["userID"])
	if err != nil {
		return nil, err
	}

	return nil, service.Delete(id)
}

func livenessCheckHandler(_ *http.Request, _ app.UserService) (result interface{}, err error) {
	return struct {
		Status string `json:"status"`
	}{"OK"}, nil
}

func parseUUID(str string) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, errInvalidParameter
	}
	return id, nil
}

func parseJSONFromBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if errors.Is(err, io.EOF) {
		return errEmptyJSONBody
	}
	return err
}

func getHandlerFunc(userService app.UserService, handler requestHandlerFunc, logger log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := handler(r, userService)
		if err != nil {
			logger.WithError(err).Error("failed to handle request")
			httpCode := translateError(err)
			w.WriteHeader(httpCode)
			return
		}

		if result == nil {
			return
		}
		err2 := json.NewEncoder(w).Encode(result)
		if err2 != nil {
			logger.WithError(err2).Error("failed to encode json-result")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func NewHTTPHandler(userService app.UserService, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	for _, route := range getRoutes() {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(getHandlerFunc(userService, route.Handler, logger))
	}
	router.Use(transport.GetLoggingMiddleware(logger))
	return router
}

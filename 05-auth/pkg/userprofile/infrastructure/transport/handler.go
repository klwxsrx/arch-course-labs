package transport

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/common/app/log"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/common/infrastructure/transport"
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/userprofile/app"
	"net/http"
)

type userProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type route struct {
	Name    string
	Method  string
	Pattern string
	Handler func(*app.UserProfileService, http.ResponseWriter, *http.Request)
}

func getRoutes() []route {
	return []route{
		{
			"getProfile",
			http.MethodGet,
			"/profile/{userID}",
			getProfileHandler,
		},
		{
			"updateProfile",
			http.MethodPut,
			"/profile/{userID}",
			updateProfileHandler,
		},
	}
}

func getProfileHandler(service *app.UserProfileService, w http.ResponseWriter, r *http.Request) {
	subjectID, err := parseSubjectID(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := parseUUID(mux.Vars(r)["userID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profile, err := service.Get(userID, subjectID)
	if errors.Is(err, app.ErrUserProfileNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, app.ErrPermissionDenied) {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(userProfile{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func updateProfileHandler(service *app.UserProfileService, w http.ResponseWriter, r *http.Request) {
	subjectID, err := parseSubjectID(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := parseUUID(mux.Vars(r)["userID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userProfile userProfile
	err = json.NewDecoder(r.Body).Decode(&userProfile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = service.Update(userID, userProfile.FirstName, userProfile.LastName, subjectID)
	if errors.Is(err, app.ErrPermissionDenied) {
		w.WriteHeader(http.StatusForbidden)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func parseSubjectID(r *http.Request) (uuid.UUID, error) {
	id := r.Header.Get("X-Auth-User-ID")
	return parseUUID(id)
}

func parseUUID(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}

func getHandlerFunc(
	userProfileService *app.UserProfileService,
	f func(*app.UserProfileService, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(userProfileService, w, r)
	}
}

func NewHTTPHandler(userProfileService *app.UserProfileService, logger log.Logger) (http.Handler, error) {
	router := mux.NewRouter()

	for _, route := range getRoutes() {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(getHandlerFunc(userProfileService, route.Handler))
	}

	router.Use(transport.NewLoggingMiddleware(logger))
	return router, nil
}
